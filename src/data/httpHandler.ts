import type {Product, Order} from './entities';
import Axios from 'axios';
import {dev} from '$app/environment';
import {fetchline} from "../routes/progress/text";
import { messageStore, parametersStore } from "../store/stores";
import {config} from "../store/stores"
import { Config, Parameters } from "./unity-cs-manager/protobuf/config";

let text = '';
let lines = 0;
const maxLines = 100;
let headers = {};
let marketplaceowner = "unity-sds"
let marketplacerepo = "unity-marketplace"
const unsubscribe = config.subscribe(configValue => {
    if (configValue && configValue.applicationConfig && configValue.applicationConfig.GithubToken) {
        headers = {
            'Authorization': `token ${configValue.applicationConfig.GithubToken}`
        };

    if(configValue && configValue.applicationConfig.MarketplaceOwner){
        marketplaceowner = configValue.applicationConfig.MarketplaceOwner
    }
    if(configValue && configValue.applicationConfig.MarketplaceUser){
        marketplacerepo = configValue.applicationConfig.MarketplaceUser
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
    private websocket: WebSocket | null = null;
    public message = '';
    async loadProducts(): Promise<Product[]> {
        if (!dev) {
            //const response = await Axios.get<Product[]>(urls.products);
            //return response.data;
            const c = generateMarketplace()
            console.log(c);
            return c;
        } else {
            return mock_products;
        }
    }

    installSoftwareSocket(meta: string) {
        if (!dev){
            this.websocket = new WebSocket('ws://localhost:8080/ws');

            this.websocket.onmessage = (event) => {
                // Append new messages to the existing text
                this.message += event.data + '\n';
                messageStore.update(message => message + event.data + '\n');
            };

            this.websocket.onerror = (error) => {
                this.message += 'WebSocket error: ' + error + '\n';
            };

            this.websocket.onclose = () => {
                this.message += 'WebSocket connection closed\n';
            };
            console.log("Sending message")
            let message: { payload: { value: string; key: string }[]; action: string }
            if (meta != null) {
                message = {
                    action: "install software",
                    payload: [{"key": "sps", "value": meta}]
                }
            } else {
                message = {
                    action: "config upgrade",
                    payload: [{ "key": "abc", "value": "def" }]
                };
            }


            this.websocket.onopen = () => {
                if (this.websocket!=null){
                    this.websocket.send(JSON.stringify(message));
                    console.log("Message sent")
                }
            }


        } else {
            const interval2 = setInterval(() => {
                lines++;
                text += fetchline(lines)
                messageStore.update(message => message + text)
                // Scroll to the bottom
                const textarea = document.getElementById('console');
                if(textarea != null){
                    textarea.scrollTop = textarea.scrollHeight;
                }
                if (lines >= maxLines) {
                    clearInterval(interval2);
                }
            }, 100);
            return () => {
                clearInterval(interval2);
            }
        }
    }

    closeSocket(): void{
        if (this.websocket) {
        this.websocket.close();
    }
    }

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

    async fetchParams(): Promise<string> {
        if (!dev) {
            const relativePath = '/ws'; // Relative path to the WebSocket endpoint
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const host = window.location.host;
            const websocketUrl = `${protocol}//${host}${relativePath}`;
            this.websocket = new WebSocket(websocketUrl);

            this.websocket.onmessage = async (event) => {
                console.log("Parameter message received: " + event.data)
                const arrayBuffer = await new Response(event.data).arrayBuffer();
                const encodedConfig = new Uint8Array(arrayBuffer);
                const decodeConfig = Parameters.decode(encodedConfig)
                console.log(decodeConfig)
                parametersStore.set(decodeConfig)
            };

            this.websocket.onerror = (error) => {
                console.log('WebSocket error: ' + error);
            };

            this.websocket.onclose = () => {
                console.log('WebSocket connection closed');
            };
            console.log("Sending message")
            const message = {
                action: "request parameters",
                payload: null
            }

            this.websocket.onopen = () => {
                if (this.websocket!=null){
                    this.websocket.send(JSON.stringify(message));
                    console.log("Config requested")
                }
            }
            return "";
        } else {
            return "abc";
        }
    }
    async fetchConfig(): Promise<string> {
        if (!dev) {
            const relativePath = '/ws'; // Relative path to the WebSocket endpoint
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const host = window.location.host;
            const websocketUrl = `${protocol}//${host}${relativePath}`;
            //const response = await Axios.post<{ id: string }>(urls.install, installData);
            this.websocket = new WebSocket(websocketUrl);

            this.websocket.onmessage = async (event) => {
                console.log("Config message received: " + event.data)
                const arrayBuffer = await new Response(event.data).arrayBuffer();
                const encodedConfig = new Uint8Array(arrayBuffer);
                const decodeConfig = Config.decode(encodedConfig)
                console.log(decodeConfig)
                config.set(decodeConfig)
            };

            this.websocket.onerror = (error) => {
                console.log('WebSocket error: ' + error);
            };

            this.websocket.onclose = () => {
                console.log('WebSocket connection closed');
            };
            console.log("Sending message")
            const message = {
                action: "request config",
                payload: null
            }

            this.websocket.onopen = () => {
                if (this.websocket!=null){
                    this.websocket.send(JSON.stringify(message));
                    console.log("Config requested")
                }
            }
            return "";
        } else {
            return "abc";
        }
    }

    async installSoftware(install: any): Promise<string> {
        console.log("installing: "+JSON.stringify(install))
        if (!dev) {
            const relativePath = '/ws'; // Relative path to the WebSocket endpoint
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const host = window.location.host;
            const websocketUrl = `${protocol}//${host}${relativePath}`;
            //const response = await Axios.post<{ id: string }>(urls.install, installData);
            this.websocket = new WebSocket(websocketUrl);

            this.websocket.onmessage = (event) => {
                // Append new messages to the existing text
                this.message += event.data + '\n';
                messageStore.update(message => message + event.data + '\n');
            };

            this.websocket.onerror = (error) => {
                this.message += 'WebSocket error: ' + error + '\n';
            };

            this.websocket.onclose = () => {
                this.message += 'WebSocket connection closed\n';
            };
            console.log("Sending message")
            let message: { payload: any; action: string }
            let extrapayload = new Uint8Array
            if (install != null) {
                message = {
                    action: "install software",
                    payload: null
                }
                extrapayload = install
            } else {
                message = {
                    action: "config upgrade",
                    payload: install
                };
            }


            this.websocket.onopen = () => {
                if (this.websocket!=null){
                    this.websocket.send(JSON.stringify(message));
                    console.log("Message sent")
                    if (extrapayload != null && extrapayload.length>0){
                        console.log("Other Message sent")

                        this.websocket.send(extrapayload)
                    }
                }
            }
            return "";
        } else {
            return "abc";
        }
    }
}

interface GithubContent {
    name: string;
    path: string;
    type: string;
}


async function generateMarketplace(): Promise<Product[]> {


    const c = await getRepoContents(marketplaceowner, marketplacerepo);

    const products: Product[] = []
    for (const p of c) {
        const content = await getGitHubFileContents(marketplaceowner, marketplacerepo, p)
        const prod: Product = JSON.parse(content)
        products.push(prod)
    }

    return products
}

async function getRepoContents(user: string, repo: string, path = ''): Promise<string[]> {
    const url = `/repos/${user}/${repo}/contents/${path}`;

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
