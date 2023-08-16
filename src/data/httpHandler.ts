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
debugger
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
  "\t\t\"Name\": \"sample application\",\n" +
  "\t\t\"Version\": \"0.1-beta\",\n" +
  "\t\t\"Channel\": \"beta\",\n" +
  "\t\t\"Owner\": \"Tom Barber\",\n" +
  "\t\t\"Description\": \"A demonstration application for the Unity platform\",\n" +
  "\t\t\"Repository\": \"https://github.com/unity-sds/unity-marketplace\",\n" +
  "\t\t\"Tags\": [\n" +
  "\t\t\t\"tag a\",\n" +
  "\t\t\t\"tag b\"\n" +
  "\t\t],\n" +
  "\t\t\"Category\": \"data processing\",\n" +
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
  "\t\t\"Package\": \"http://github.com/path/to/package.zip\",\n" +
  "\t\t\"ManagedDependencies\": [{\n" +
  "\t\t\t\"Eks\": {\n" +
  "\t\t\t\t\"MinimumVersion\": \"1.21\"\n" +
  "\t\t\t}\n" +
  "\t\t}],\n" +
  "\t\t\"Backend\": \"terraform\",\n" +
  "\t\t\"DefaultDeployment\": {\n" +
  "\t\t\t\"Variables\": {\n" +
  "\t\t\t\t\"some_terraform_variable\": \"some_value\"\n" +
  "\t\t\t},\n" +
  "\t\t\t\"EksSpec\": {\n" +
  "\t\t\t\t\"NodeGroups\": [{\n" +
  "\t\t\t\t\t\"NodeGroup1\": {\n" +
  "\t\t\t\t\t\t\"MinNodes\": 1,\n" +
  "\t\t\t\t\t\t\"MaxNodes\": 10,\n" +
  "\t\t\t\t\t\t\"DesiredNodes\": 4,\n" +
  "\t\t\t\t\t\t\"InstanceType\": \"m6.large\"\n" +
  "\t\t\t\t\t}\n" +
  "\t\t\t\t}]\n" +
  "\t\t\t}\n" +
  "\t\t}\n" +
  "\t},\n" +
  "\t{\n" +
  "\t\t\"DisplayName\": \"Unity API Gateway\",\n" +
  "\t\t\"Name\": \"unity-apigateway\",\n" +
  "\t\t\"Version\": \"0.1-beta\",\n" +
  "\t\t\"Channel\": \"beta\",\n" +
  "\t\t\"Owner\": \"U-CS Team\",\n" +
  "\t\t\"Description\": \"A package to install and configure an API gateway for your Unity Venue\",\n" +
  "\t\t\"Repository\": \"https://github.com/unity-sds/unity-cs-infra/\",\n" +
  "\t\t\"Tags\": [\n" +
  "\t\t\t\"api\",\n" +
  "\t\t\t\"http\",\n" +
  "\t\t\t\"rest\"\n" +
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
  "\t\t\"Package\": \"https://github.com/unity-sds/unity-cs-infra/\",\n" +
  "\t\t\"Backend\": \"terraform\",\n" +
  "\t\t\"WorkDirectory\": \"terraform-project-api-gateway_module\",\n" +
  "\t\t\"DefaultDeployment\": {\n" +
  "\t\t\t\"Variables\": {\n" +
  "\t\t\t\t\"some_terraform_variable\": \"some_value\"\n" +
  "\t\t\t}\n" +
  "\t\t}\n" +
  "\t},\n" +
  "\t{\n" +
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
  "\t\t\"DefaultDeployment\": {\n" +
  "\t\t\t\"Variables\": {\n" +
  "\t\t\t\t\"Values\": {\n" +
  "\t\t\t\t\t\"some_terraform_variable\": \"some_value\"\n" +
  "\t\t\t\t},\n" +
  "\t\t\t\t\"NestedValues\": {\n" +
  "\t\t\t\t\t\"NodeGroups\": {\n" +
  "\t\t\t\t\t\t\"Config\": {\n" +
  "\t\t\t\t\t\t\t\"Name\": {\n" +
  "\t\t\t\t\t\t\t\t\"Options\": {\n" +
  "\t\t\t\t\t\t\t\t\t\"type\": \"String\",\n" +
  "\t\t\t\t\t\t\t\t\t\"default\": \"NodeGroup\"\n" +
  "\t\t\t\t\t\t\t\t}\n" +
  "\t\t\t\t\t\t\t},\n" +
  "\t\t\t\t\t\t\t\"MinNodes\": {\n" +
  "\t\t\t\t\t\t\t\t\"Options\": {\n" +
  "\t\t\t\t\t\t\t\t\t\"type\": \"Number\",\n" +
  "\t\t\t\t\t\t\t\t\t\"default\": \"1\"\n" +
  "\t\t\t\t\t\t\t\t}\n" +
  "\t\t\t\t\t\t\t},\n" +
  "\t\t\t\t\t\t\t\"MaxNodes\": {\n" +
  "\t\t\t\t\t\t\t\t\"Options\": {\n" +
  "\t\t\t\t\t\t\t\t\t\"type\": \"Number\",\n" +
  "\t\t\t\t\t\t\t\t\t\"default\": \"3\"\n" +
  "\t\t\t\t\t\t\t\t}\n" +
  "\t\t\t\t\t\t\t},\n" +
  "\t\t\t\t\t\t\t\"DesiredNodes\": {\n" +
  "\t\t\t\t\t\t\t\t\"Options\": {\n" +
  "\t\t\t\t\t\t\t\t\t\"type\": \"Number\",\n" +
  "\t\t\t\t\t\t\t\t\t\"default\": \"1\"\n" +
  "\t\t\t\t\t\t\t\t}\n" +
  "\t\t\t\t\t\t\t},\n" +
  "\t\t\t\t\t\t\t\"InstanceType\": {\n" +
  "\t\t\t\t\t\t\t\t\"Options\": {\n" +
  "\t\t\t\t\t\t\t\t\t\"type\": \"String\",\n" +
  "\t\t\t\t\t\t\t\t\t\"default\": \"m6.xlarge\"\n" +
  "\t\t\t\t\t\t\t\t}\n" +
  "\t\t\t\t\t\t\t}\n" +
  "\t\t\t\t\t\t}\n" +
  "\t\t\t\t\t}\n" +
  "\t\t\t\t}\n" +
  "\t\t\t}\n" +
  "\t\t}\n" +
  "\t}, {\n" +
  "\t\t\"Name\": \"Unity SPS\",\n" +
  "\t\t\"Version\": \"0.1-beta\",\n" +
  "\t\t\"Channel\": \"beta\",\n" +
  "\t\t\"Owner\": \"Tom Barber\",\n" +
  "\t\t\"Description\": \"The Unity SPS Prototype package\",\n" +
  "\t\t\"Repository\": \"https://github.com/unity-sds/unity-sps-prototype\",\n" +
  "\t\t\"Tags\": [\n" +
  "\t\t\t\"sps\",\n" +
  "\t\t\t\"data processing\"\n" +
  "\t\t],\n" +
  "\t\t\"Category\": \"data processing\",\n" +
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
  "\t\t\"Package\": \"https://github.com/unity-sds/unity-sps-prototype/archive/refs/tags/u-cs-deployment.zip\",\n" +
  "\t\t\"ManagedDependencies\": [{\n" +
  "\t\t\t\"Eks\": {\n" +
  "\t\t\t\t\"MinimumVersion\": \"1.21\"\n" +
  "\t\t\t}\n" +
  "\t\t}],\n" +
  "\t\t\"Backend\": \"terraform\",\n" +
  "\t\t\"DefaultDeployment\": {\n" +
  "\t\t\t\"Variables\": {\n" +
  "\t\t\t\t\"some_terraform_variable\": \"some_value\"\n" +
  "\t\t\t},\n" +
  "\t\t\t\"EksSpec\": {\n" +
  "\t\t\t\t\"NodeGroups\": [{\n" +
  "\t\t\t\t\t\"NodeGroup1\": {\n" +
  "\t\t\t\t\t\t\"MinNodes\": 1,\n" +
  "\t\t\t\t\t\t\"MaxNodes\": 10,\n" +
  "\t\t\t\t\t\t\"DesiredNodes\": 4,\n" +
  "\t\t\t\t\t\t\"InstanceType\": \"m6.large\"\n" +
  "\t\t\t\t\t}\n" +
  "\t\t\t\t}]\n" +
  "\t\t\t}\n" +
  "\t\t}\n" +
  "\t}\n" +
  "]"