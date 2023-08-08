import type { Product, Order, Application, InstallationApplication } from "./entities";
import Axios from 'axios';
import {dev} from '$app/environment';
import { marketplaceStore, messageStore, parametersStore } from "../store/stores";
import {config} from "../store/stores"
import { websocketStore } from '../data/websocketStore';
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

	async storeOrder(order: Order): Promise<number> {
		if (!dev) {
			const orderData = {
				lines: [...order.orderLines.values()].map((ol) => ({
					productId: ol.product.Id,
					productName: ol.product.Name,
					quantity: ol.quantity
				}))
			};
			const response = await Axios.post<{ id: number }>(urls.orders, orderData);
			return response.data.id;
		} else {
			return 1;
		}
	}

	async installSoftware(
		install: InstallationApplication[],
		deploymentName: string
	): Promise<string> {
		const appRequest = Install_Applications.create({
			name: install[0].name,
			version: install[0].version,
			variables: { test: '' }
		});
		const installrequest = Install.create({
			applications: appRequest,
			DeploymentName: deploymentName
		});
		const installSoftware = UnityWebsocketMessage.create({ install: installrequest });

		websocketStore.send(UnityWebsocketMessage.encode(installSoftware).finish());
		return '';
	}

	async setupws() {
		const messages = console.log('Setting up connection');
		const set = ConnectionSetup.create({ type: 'register', userID: 'test' });
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
				if (message.parameters) {
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
    console.log("Checking if manifest.json exists in the repository...");
    const manifestExists = await checkIfFileExists(marketplaceowner, marketplacerepo, 'manifest.json');
    if (manifestExists) {
        console.log("manifest.json exists in the repository.");
        const content = await getGitHubFileContents(marketplaceowner, marketplacerepo, "manifest.json")
        const c = JSON.parse(content)
        const products: Product[] = []
        for (const p of c) {
            const prod: Product = p
            products.push(prod)
        }
        marketplaceStore.set(products)
        return; // You can decide what to do if the file exists. Here I'm just exiting the function.
    }

		console.log("fetching repo contents: "+marketplaceowner)
    const c = await getRepoContents(marketplaceowner, marketplacerepo);

    const products: Product[] = []
    for (const p of c) {
        const content = await getGitHubFileContents(marketplaceowner, marketplacerepo, p)
        const prod: Product = JSON.parse(content)
        products.push(prod)
    }

    marketplaceStore.set(products)
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

const mock_products = [
    {
        Id: 1,
        Name: 'sample application',
        Version: '1.2.3',
        Branch: '',
        Category: 'data processing',
        Description: 'A demonstration application for the Unity platform',
        Tags: ["tag a", "tag b"],
        IamRoles: [
            {
                "Version": "2012-10-17",
                "Statement": [
                    {
                        "Effect": "Allow",
                        "Action": "service-prefix:action-name",
                        "Resource": "*",
                        "Condition": {
                            "DateGreaterThan": {"aws:CurrentTime": "2020-04-01T00:00:00Z"},
                            "DateLessThan": {"aws:CurrentTime": "2020-06-30T23:59:59Z"}
                        }
                    }
                ]
            }
        ],
        Package: "http://github.com/path/to/package.zip",
        ManagedDependencies: [
            {
                Eks: {
                    MinimumVersion: "1.21"
                }
            }
        ],
        Backend: "terraform",
        DefaultDeployment: {
            Variables: {
                "some_terraform_variable": "some value"
            },
            EksSpec: {
                NodeGroups: [
                    {
                        NodeGroup1: {
                            MinNodes: 1,
                            MaxNodes: 10,
                            DesiredNodes: 4,
                            InstanceType: "m6.large"
                        }
                    }
                ]
            }
        }
    },
    {
        Id: 1,
        Name: "Unity Kubernetes",
        Version: "0.1-beta",
        Branch: '',
        Channel: "beta",
        Owner: "Tom Barber",
        Description: "The Unity Kubernetes package",
        Repository: "https://github.com/unity-sds/unity-cs-infra",
        Tags: [
            "eks",
            "kubernetes"
        ],
        Category: "system",
        IamRoles: [{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Action": "service-prefix:action-name",
                    "Resource": "*",
                    "Condition": {
                        "DateGreaterThan": {"aws:CurrentTime": "2020-04-01T00:00:00Z"},
                        "DateLessThan": {"aws:CurrentTime": "2020-06-30T23:59:59Z"}
                    }
                }
            ]
        }],
        Package: "https://github.com/unity-sds/unity-cs-infra",
        Backend: "./github/workflows/deploy_eks.yml",
        ManagedDependencies: [
            {
                "Eks": {
                    "MinimumVersion": "1.21"
                }
            }
        ],
        DefaultDeployment: {
            Variables: {
                "some_terraform_variable": "some_value"
            },
            EksSpec: {
                NodeGroups: [
                    {
                        NodeGroup1: {
                            MinNodes: 1,
                            MaxNodes: 10,
                            DesiredNodes: 4,
                            InstanceType: "m6.large"
                        }
                    }
                ]
            }
        }
    }
];
