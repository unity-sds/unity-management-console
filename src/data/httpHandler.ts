import Axios from 'axios';
import { dev } from '$app/environment';
import {
	deploymentStore,
	installError,
	installRunning,
	marketplaceStore,
	messageStore,
	parametersStore
} from '../store/stores';
import { config } from '../store/stores';
import { websocketStore } from '../data/websocketStore';
import { MarketplaceMetadata } from './unity-cs-manager/protobuf/marketplace';
import {
	ConnectionSetup,
	Parameters,
	SimpleMessage,
	UnityWebsocketMessage,
	Install_Applications,
	Install,
	Parameters_Parameter,
	Uninstall
} from './unity-cs-manager/protobuf/extensions';

let headers = {};
let marketplaceowner = 'unity-sds';
let marketplacerepo = 'unity-marketplace';

const unsubscribe = config.subscribe((configValue) => {
	if (configValue && configValue.applicationConfig && configValue.applicationConfig.GithubToken) {
		headers = {
			Authorization: `token ${configValue.applicationConfig.GithubToken}`
		};

		if (
			configValue &&
			configValue.applicationConfig.MarketplaceOwner &&
			configValue.applicationConfig.MarketplaceUser
		) {
			console.log('Setting marketplace owner: ' + configValue.applicationConfig.MarketplaceOwner);
			console.log('Setting marketplace user: ' + configValue.applicationConfig.MarketplaceUser);

			marketplaceowner = configValue.applicationConfig.MarketplaceOwner;
			marketplacerepo = configValue.applicationConfig.MarketplaceUser;

			generateMarketplace();
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

	async installSoftware(install: Install_Applications, deploymentName: string): Promise<string> {
		const installrequest = Install.create({
			applications: install,
			DeploymentName: deploymentName
		});
		const installSoftware = UnityWebsocketMessage.create({ install: installrequest });
		websocketStore.send(UnityWebsocketMessage.encode(installSoftware).finish());
		return '';
	}

	async uninstallAllSoftware(): Promise<string> {
		const uninstallAll = Uninstall.create({
			All: true
		});
		const uninstallSoftware = UnityWebsocketMessage.create({ uninstall: uninstallAll });
		websocketStore.send(UnityWebsocketMessage.encode(uninstallSoftware).finish());
		return '';
	}

	setupws() {
		if (!dev) {
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
						if (message.simplemessage.operation === 'terraform') {
							if (message.simplemessage.payload === 'completed') {
								installRunning.set(false);
							}

							if (message.simplemessage.payload === 'failed') {
								installError.set(true);
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
					} else if (message.deployments) {
						debugger;
						deploymentStore.set(message.deployments);
					}
					lastProcessedIndex = i; // Update the last processed index
				}
			});
		} else {
			generateMarketplace();
		}
	}

	updateParameters(p: { [p: string]: Parameters_Parameter }) {
		const paramMessage = Parameters.fromPartial({ parameterlist: p });
		const message = UnityWebsocketMessage.create({ parameters: paramMessage });
		websocketStore.send(UnityWebsocketMessage.encode(message).finish());
	}
}

interface GithubContent {
	name: string;
	path: string;
	type: string;
}

export async function uninstallApplication(name: string, appPackage: string, deployment: string) {
	const uninstallMessage = SimpleMessage.create({
		operation: 'uninstall application',
		payload:
			'{ "Application": "' +
			name +
			'", "ApplicationPackage": "' +
			appPackage +
			'", "Deployment":"' +
			deployment +
			'"}'
	});

	const m = UnityWebsocketMessage.create({ simplemessage: uninstallMessage });
	websocketStore.send(UnityWebsocketMessage.encode(m).finish());
}

export async function reapplyApplication(name: string, appPackage: string, deployment: string) {
	const reapplyMessage = SimpleMessage.create({
		operation: 'reapply application',
		payload:
			'{ "Application": "' +
			name +
			'", "ApplicationPackage": "' +
			appPackage +
			'", "Deployment":"' +
			deployment +
			'"}'
	});
}

export async function fetchDeployedApplications() {
	if (!dev) {
		const paramrequest = SimpleMessage.create({
			operation: 'request all applications',
			payload: ''
		});
		const wsm2 = UnityWebsocketMessage.create({ simplemessage: paramrequest });

		websocketStore.send(UnityWebsocketMessage.encode(wsm2).finish());
	}
}

async function generateMarketplace() {
	if (!dev) {
		console.log('Checking if manifest.json exists in the repository...');
		const manifestExists = await checkIfFileExists(
			marketplaceowner,
			marketplacerepo,
			'manifest.json'
		);
		if (manifestExists) {
			console.log('manifest.json exists in the repository.');
			const content = await getGitHubFileContents(
				marketplaceowner,
				marketplacerepo,
				'manifest.json'
			);
			const c = JSON.parse(content);
			const products: MarketplaceMetadata[] = [];
			for (const p of c) {
				const prod = MarketplaceMetadata.fromJSON(p);
				products.push(prod);
			}
			marketplaceStore.set(products);
			return;
		}

		console.log('fetching repo contents: ' + marketplaceowner);
		const c = await getRepoContents(marketplaceowner, marketplacerepo);

		const products: MarketplaceMetadata[] = [];
		for (const p of c) {
			const content = await getGitHubFileContents(marketplaceowner, marketplacerepo, p);
			const j = JSON.parse(content);
			const prod = MarketplaceMetadata.fromJSON(j);
			products.push(prod);
		}

		marketplaceStore.set(products);
	} else {
		const j = JSON.parse(mock_marketplace);
		const products: MarketplaceMetadata[] = [];
		for (const p of j) {
			const prod = MarketplaceMetadata.fromJSON(p);
			products.push(prod);
		}
		marketplaceStore.set(products);
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

	console.log('fetching: ' + url);
	const paths: string[] = [];
	try {
		const api = Axios.create({
			baseURL: 'https://api.github.com',
			headers: headers
		});

		const response = await api.get<GithubContent[]>(url);
		const data = response.data;

		for (const item of data) {
			if (item.path.includes('metadata.json')) {
				console.log(item.path);
				paths.push(item.path);
			}
			if (item.type === 'dir') {
				// If the item is a directory, recursively fetch its contents
				const dirPaths = await getRepoContents(user, repo, item.path);
				paths.push(...dirPaths); // Add the results of the recursive call to paths
			}
		}
	} catch (error: unknown) {
		if (Axios.isAxiosError(error)) {
			console.error(`Error fetching directory listing: ${error.message}`);
		}
	}

	console.log('returning: ' + paths);
	return paths;
}

function decodeBase64(data: string): string {
	return atob(data);
}

async function getGitHubFileContents(user: string, repo: string, path: string): Promise<string> {
	const url = `/repos/${user}/${repo}/contents/${path}`;

	console.log('fetching: ' + url);
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

const mock_marketplace =
	'[\n' +
	'  {\n' +
	'    "Name": "sample application",\n' +
	'    "Version": "0.1-beta",\n' +
	'    "Channel": "beta",\n' +
	'    "Owner": "Tom Barber",\n' +
	'    "Description": "A demonstration application for the Unity platform",\n' +
	'    "Repository": "https://github.com/unity-sds/unity-marketplace",\n' +
	'    "Tags": [\n' +
	'      "tag a",\n' +
	'      "tag b"\n' +
	'    ],\n' +
	'    "Category": "data processing",\n' +
	'    "IamRoles": {\n' +
	'      "Statement": [\n' +
	'        {\n' +
	'          "Effect": "Allow",\n' +
	'          "Action": [\n' +
	'            "iam:CreateInstanceProfile"\n' +
	'          ],\n' +
	'          "Resource": [\n' +
	'            "arn:aws:iam::<account_id>:instance-profile/eksctl*"\n' +
	'          ]\n' +
	'        }\n' +
	'      ]\n' +
	'    },\n' +
	'    "Package": "http://github.com/path/to/package.zip",\n' +
	'    "ManagedDependencies": [\n' +
	'      {\n' +
	'        "Eks": {\n' +
	'          "MinimumVersion": "1.21"\n' +
	'        }\n' +
	'      }\n' +
	'    ],\n' +
	'    "Backend": "terraform",\n' +
	'    "DefaultDeployment": {\n' +
	'      "Variables": {\n' +
	'        "some_terraform_variable": "some_value"\n' +
	'      },\n' +
	'      "EksSpec": {\n' +
	'        "NodeGroups": [\n' +
	'          {\n' +
	'            "NodeGroup1": {\n' +
	'              "MinNodes": 1,\n' +
	'              "MaxNodes": 10,\n' +
	'              "DesiredNodes": 4,\n' +
	'              "InstanceType": "m6.large"\n' +
	'            }\n' +
	'          }\n' +
	'        ]\n' +
	'      }\n' +
	'    }\n' +
	'  },\n' +
	'  {\n' +
	'    "DisplayName": "Unity API Gateway",\n' +
	'    "Name": "unity-apigateway",\n' +
	'    "Version": "0.1-beta",\n' +
	'    "Channel": "beta",\n' +
	'    "Owner": "U-CS Team",\n' +
	'    "Description": "A package to install and configure an API gateway for your Unity Venue",\n' +
	'    "Repository": "https://github.com/unity-sds/unity-cs-infra/",\n' +
	'    "Tags": [\n' +
	'      "api",\n' +
	'      "http",\n' +
	'      "rest"\n' +
	'    ],\n' +
	'    "Category": "system",\n' +
	'    "IamRoles": {\n' +
	'      "Statement": [\n' +
	'        {\n' +
	'          "Effect": "Allow",\n' +
	'          "Action": [\n' +
	'            "iam:CreateInstanceProfile"\n' +
	'          ],\n' +
	'          "Resource": [\n' +
	'            "arn:aws:iam::<account_id>:instance-profile/eksctl*"\n' +
	'          ]\n' +
	'        }\n' +
	'      ]\n' +
	'    },\n' +
	'    "Package": "https://github.com/unity-sds/unity-cs-infra/",\n' +
	'    "Backend": "terraform",\n' +
	'    "WorkDirectory": "terraform-project-api-gateway_module",\n' +
	'    "DefaultDeployment": {\n' +
	'      "Variables": {\n' +
	'        "some_terraform_variable": "some_value"\n' +
	'      }\n' +
	'    }\n' +
	'  },\n' +
	'  {\n' +
	'    "DisplayName": "Unity Kubernetes",\n' +
	'    "Name": "unity-eks",\n' +
	'    "Version": "0.1",\n' +
	'    "Channel": "beta",\n' +
	'    "Owner": "Tom Barber",\n' +
	'    "Description": "The Unity Kubernetes package",\n' +
	'    "Repository": "https://github.com/unity-sds/unity-cs-infra",\n' +
	'    "Tags": [\n' +
	'      "eks",\n' +
	'      "kubernetes"\n' +
	'    ],\n' +
	'    "Category": "system",\n' +
	'    "IamRoles": {\n' +
	'      "Statement": [\n' +
	'        {\n' +
	'          "Effect": "Allow",\n' +
	'          "Action": [\n' +
	'            "iam:CreateInstanceProfile"\n' +
	'          ],\n' +
	'          "Resource": [\n' +
	'            "arn:aws:iam::<account_id>:instance-profile/eksctl*"\n' +
	'          ]\n' +
	'        }\n' +
	'      ]\n' +
	'    },\n' +
	'    "Package": "https://github.com/unity-sds/unity-cs-infra",\n' +
	'    "WorkDirectory": "terraform-unity-eks_module",\n' +
	'    "Backend": "terraform",\n' +
	'    "ManagedDependencies": [],\n' +
	'    "PostInstall": "scripts/postinstall.sh",\n' +
	'    "DefaultDeployment": {\n' +
	'      "Variables": {\n' +
	'        "Values": {\n' +
	'          "cluster_version": "1.27"\n' +
	'        },\n' +
	'        "AdvancedValues": {\n' +
	'          "nodegroups": {\n' +
	'            "UnityNodeGroup": {\n' +
	'              "min_size": 1,\n' +
	'              "max_size": 10,\n' +
	'              "desired_size": 1,\n' +
	'              "instance_types": [\n' +
	'                "t3.large"\n' +
	'              ],\n' +
	'              "capacity_type": "SPOT"\n' +
	'            }\n' +
	'          }\n' +
	'        }\n' +
	'      }\n' +
	'    }\n' +
	'  },\n' +
	'  {\n' +
	'    "Name": "Unity SPS",\n' +
	'    "Version": "0.1-beta",\n' +
	'    "Channel": "beta",\n' +
	'    "Owner": "Tom Barber",\n' +
	'    "Description": "The Unity SPS Prototype package",\n' +
	'    "Repository": "https://github.com/unity-sds/unity-sps-prototype",\n' +
	'    "Tags": [\n' +
	'      "sps",\n' +
	'      "data processing"\n' +
	'    ],\n' +
	'    "Category": "data processing",\n' +
	'    "IamRoles": {\n' +
	'      "Statement": [\n' +
	'        {\n' +
	'          "Effect": "Allow",\n' +
	'          "Action": [\n' +
	'            "iam:CreateInstanceProfile"\n' +
	'          ],\n' +
	'          "Resource": [\n' +
	'            "arn:aws:iam::<account_id>:instance-profile/eksctl*"\n' +
	'          ]\n' +
	'        }\n' +
	'      ]\n' +
	'    },\n' +
	'    "Package": "https://github.com/buggtb/unity-sps-prototype@bebb4f5cc092b88fb583c4f36cecefdb7f037244",\n' +
	'    "PreInstall": "scripts/preinstall.sh",\n' +
	'    "ManagedDependencies": [\n' +
	'      {\n' +
	'        "unity-eks": {\n' +
	'          "MinimumVersion": "1.21"\n' +
	'        }\n' +
	'      }\n' +
	'    ],\n' +
	'    "Backend": "terraform",\n' +
	'    "DefaultDeployment": {\n' +
	'      "Variables": {\n' +
	'        "Values": {\n' +
	'          "release": "23.1",\n' +
	'          "deployment_name": "",\n' +
	'          "eks_cluster_name": "",\n' +
	'          "kubeconfig_filepath": "",\n' +
	'          "venue": "",\n' +
	'          "uads_development_efs_fsmt_id": "",\n' +
	'          "default_group_node_group_name": "LoadBalancer",\n' +
	'          "default_group_node_group_launch_template_name": "",\n' +
	'          "elb_subnets": ""\n' +
	'        },\n' +
	'        "NestedValues": {}\n' +
	'      }\n' +
	'    }\n' +
	'  }\n' +
	']';
