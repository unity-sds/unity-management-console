import Axios from 'axios';
import {dev} from '$app/environment';
import { installError, installRunning, marketplaceStore, messageStore, parametersStore } from "../store/stores";
import {config} from "../store/stores"
import { websocketStore } from '../data/websocketStore';
import {
  MarketplaceMetadata, MarketplaceMetadata_InnerMap, MarketplaceMetadata_SubMap,
  MarketplaceMetadata_TypeMap,
  MarketplaceMetadata_Variables
} from "./unity-cs-manager/protobuf/marketplace";
import {
    ConnectionSetup,
    Parameters,
    SimpleMessage,
    UnityWebsocketMessage,
    Install_Applications, Install, Parameters_Parameter
} from "./unity-cs-manager/protobuf/extensions";
let headers = {};
let marketplaceowner = "unity-sds"
let marketplacerepo = "unity-marketplace"

const unsubscribe = config.subscribe(configValue => {
    if (configValue && configValue.applicationConfig && configValue.applicationConfig.GithubToken) {
        headers = {
            'Authorization': `token ${configValue.applicationConfig.GithubToken}`
        };

    if(configValue && configValue.applicationConfig.MarketplaceOwner && configValue.applicationConfig.MarketplaceUser) {
        console.log("Setting marketplace owner: " + configValue.applicationConfig.MarketplaceOwner)
        console.log("Setting marketplace user: " + configValue.applicationConfig.MarketplaceUser)

        marketplaceowner = configValue.applicationConfig.MarketplaceOwner
        marketplacerepo = configValue.applicationConfig.MarketplaceUser

        generateMarketplace()
    }
    } else {
        // default or error headers if GithubToken is not available
        headers = {};
    }
});
const urls = {
    products: '/api/products',
    orders: '/api/orders',
    install: '/api/application/install'
};

export class HttpHandler {
	public message = '';

	async installSoftware(
		install: Install_Applications,
		deploymentName: string
	): Promise<string> {
		const installrequest = Install.create({
			applications: install,
			DeploymentName: deploymentName
		});
		const installSoftware = UnityWebsocketMessage.create({ install: installrequest });
		websocketStore.send(UnityWebsocketMessage.encode(installSoftware).finish());
		return '';
	}

	async setupws() {
    if(!dev) {
      const set = ConnectionSetup.create({ type: 'register', userID: 'test' });
      console.log(ConnectionSetup.toJSON(set))
      websocketStore.send(ConnectionSetup.encode(set).finish());
      const configrequest = SimpleMessage.create({ operation: 'request config', payload: '' });
      const wsm = UnityWebsocketMessage.create({ simplemessage: configrequest });
      const paramrequest = SimpleMessage.create({ operation: 'request parameters', payload: '' });
      const wsm2 = UnityWebsocketMessage.create({ simplemessage: paramrequest });

      websocketStore.send(UnityWebsocketMessage.encode(wsm).finish());
      websocketStore.send(UnityWebsocketMessage.encode(wsm2).finish());
      let lastProcessedIndex = -1;
      const unsubscribe = websocketStore.subscribe((receivedMessages) => {
        // loop through the received messages
        for (let i = lastProcessedIndex + 1; i < receivedMessages.length; i++) {
          const message = receivedMessages[i];
          if (message.simplemessage) {

            if (message.simplemessage.operation === "terraform") {
              if (message.simplemessage.payload === "completed") {
                installRunning.set(false)
              }

              if (message.simplemessage.payload === "failed") {
                installError.set(true)
              }

            }
          } else if (message.parameters) {
            parametersStore.set(message.parameters);
          } else if (message.config) {
            config.set(message.config);
          } else if (message.logs) {
            if (message.logs.line != undefined) {
              console.log(message.logs?.line);
              messageStore.update(
                (messages) => `${messages}[${message.logs?.level}] ${message.logs?.line}\n`
              );
            }
          }
          lastProcessedIndex = i; // Update the last processed index
        }
      });
    } else {
      generateMarketplace()
    }
	}

	updateParameters(p: { [p: string]: Parameters_Parameter }) {
		const paramMessage = Parameters.fromPartial({parameterlist: p})
		const message = UnityWebsocketMessage.create({ parameters: paramMessage });
		websocketStore.send(UnityWebsocketMessage.encode(message).finish());
	}
}

interface GithubContent {
    name: string;
    path: string;
    type: string;
}


async function generateMarketplace(){
  if(!dev) {
    console.log("Checking if manifest.json exists in the repository...");
    const manifestExists = await checkIfFileExists(marketplaceowner, marketplacerepo, 'manifest.json');
    if (manifestExists) {
      console.log("manifest.json exists in the repository.");
      const content = await getGitHubFileContents(marketplaceowner, marketplacerepo, "manifest.json")
      const c = JSON.parse(content)
      const products: MarketplaceMetadata[] = []
      for (const p of c) {
        const prod= MarketplaceMetadata.fromJSON(p)
        products.push(prod)
      }
      marketplaceStore.set(products)
      return;
    }

    console.log("fetching repo contents: " + marketplaceowner)
    const c = await getRepoContents(marketplaceowner, marketplacerepo);

    const products: MarketplaceMetadata[] = []
    for (const p of c) {
      const content = await getGitHubFileContents(marketplaceowner, marketplacerepo, p)
      const j = JSON.parse(content)
      const prod= MarketplaceMetadata.fromJSON(j)
      products.push(prod)
    }

    marketplaceStore.set(products)
  } else {
    const j = JSON.parse(mock_marketplace)
    const products: MarketplaceMetadata[] = []
    for (const p of j) {
      const prod= MarketplaceMetadata.fromJSON(p)
      products.push(prod)
    }
    marketplaceStore.set(products)
    return;
  }
}

async function checkIfFileExists(user: string, repo: string, filePath: string): Promise<boolean> {
    const url = `/repos/${user}/${repo}/contents/${filePath}`;
    try {
        const api = Axios.create({
            baseURL: 'https://api.github.com',
            headers: headers
        });

        const response = await api.get<GithubContent[]>(url);
        return response.status === 200; // Success status means the file exists.
    } catch (error: unknown) {
        if (Axios.isAxiosError(error)) {
            // If the error status is 404, the file does not exist.
            if (error.response?.status === 404) {
                return false;
            }
            console.error(`Error checking file existence: ${error.message}`);
        }
        return false;
    }
}
async function getRepoContents(user: string, repo: string, path = ''): Promise<string[]> {
    const url = `/repos/${user}/${repo}/contents/${path}`;

    console.log("fetching: "+url)
    const paths: string[] = [];
    try {
        const api = Axios.create({
            baseURL: 'https://api.github.com',
            headers: headers
        });

        const response = await api.get<GithubContent[]>(url);
        const data = response.data;

        for (const item of data) {
            if (item.path.includes("metadata.json")) {
                console.log(item.path);
                paths.push(item.path);
            }
            if (item.type === 'dir') {
                // If the item is a directory, recursively fetch its contents
                const dirPaths = await getRepoContents(user, repo, item.path);
                paths.push(...dirPaths);  // Add the results of the recursive call to paths
            }
        }
    } catch (error: unknown) {
        if (Axios.isAxiosError(error)) {
            console.error(`Error fetching directory listing: ${error.message}`);
        }
    }

    console.log("returning: " + paths);
    return paths;
}

function decodeBase64(data: string): string {
    return atob(data)
}

async function getGitHubFileContents(user: string, repo: string, path: string): Promise<string> {
    const url = `/repos/${user}/${repo}/contents/${path}`;

    console.log("fetching: "+url)
    try {
        const api = Axios.create({
            baseURL: 'https://api.github.com',
            headers: headers
        });

        const response = await api.get(url);
        const fileContent = decodeBase64(response.data.content);
        return fileContent;
    } catch (error: unknown) {
        if (Axios.isAxiosError(error)) {
            console.error(`Error fetching file contents: ${error.message}`);
        }
        throw error;
    }
}
const mock_marketplace = "[{\n" +
  "\t\t\"DisplayName\": \"Unity Kubernetes\",\n" +
  "\t\t\"Name\": \"unity-eks\",\n" +
  "\t\t\"Version\": \"0.1\",\n" +
  "\t\t\"Channel\": \"beta\",\n" +
  "\t\t\"Owner\": \"Tom Barber\",\n" +
  "\t\t\"Description\": \"The Unity Kubernetes package\",\n" +
  "\t\t\"Repository\": \"https://github.com/unity-sds/unity-cs-infra\",\n" +
  "\t\t\"Tags\": [\n" +
  "\t\t\t\"eks\",\n" +
  "\t\t\t\"kubernetes\"\n" +
  "\t\t],\n" +
  "\t\t\"Category\": \"system\",\n" +
  "\t\t\"IamRoles\": {\n" +
  "\t\t\t\"Statement\": [{\n" +
  "\t\t\t\t\"Effect\": \"Allow\",\n" +
  "\t\t\t\t\"Action\": [\n" +
  "\t\t\t\t\t\"iam:CreateInstanceProfile\"\n" +
  "\t\t\t\t],\n" +
  "\t\t\t\t\"Resource\": [\n" +
  "\t\t\t\t\t\"arn:aws:iam::<account_id>:instance-profile/eksctl*\"\n" +
  "\t\t\t\t]\n" +
  "\t\t\t}]\n" +
  "\t\t},\n" +
  "\t\t\"Package\": \"https://github.com/unity-sds/unity-cs-infra\",\n" +
  "\t\t\"WorkDirectory\": \"terraform-unity-eks_module\",\n" +
  "\t\t\"Backend\": \"terraform\",\n" +
  "\t\t\"ManagedDependencies\": [{\n" +
  "\t\t\t\"Eks\": {\n" +
  "\t\t\t\t\"MinimumVersion\": \"1.21\"\n" +
  "\t\t\t}\n" +
  "\t\t}],\n" +
  "\t\t\"PostInstall\": \"scripts/postinstall.sh\",\n" +
  "\t\t\"DefaultDeployment\": {\n" +
  "\t\t\t\"Variables\": {\n" +
  "\t\t\t\t\"Values\": {\n" +
  "\t\t\t\t\t\"cluster_version\": \"1.27\"\n" +
  "\t\t\t\t},\n" +
  "\t\t\t\t\"AdvancedValues\": {\n" +
  "\t\t\t\t\t\"nodegroups\": {\n" +
  "\t\t\t\t\t\t\"blue\": {\n" +
  "\t\t\t\t\t\t\t\"create_iam_role\":            false,\n" +
  "\t\t\t\t\t\t\t\"iam_role_arn\":               \"data.aws_ssm_parameter.eks_iam_node_role.value\",\n" +
  "\t\t\t\t\t\t\t\"min_size\":                   1,\n" +
  "\t\t\t\t\t\t\t\"max_size\":                   10,\n" +
  "\t\t\t\t\t\t\t\"desired_size\":               1,\n" +
  "\t\t\t\t\t\t\t\"ami_id\":                     \"ami-0c0e3c5bfa15ba56b\",\n" +
  "\t\t\t\t\t\t\t\"instance_types\":             [\"t3.large\"],\n" +
  "\t\t\t\t\t\t\t\"capacity_type\":              \"SPOT\",\n" +
  "\t\t\t\t\t\t\t\"enable_bootstrap_user_data\": true\n" +
  "\t\t\t\t\t\t},\n" +
  "\t\t\t\t\t\t\"green\" :{\n" +
  "\t\t\t\t\t\t}\n" +
  "\t\t\t\t\t}\n" +
  "\n" +
  "\t\t\t\t}\n" +
  "\t\t\t}\n" +
  "\t\t}\n" +
  "\t}]"