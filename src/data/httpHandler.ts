import type {Product, Order, AppInstall, Extensions, Install} from './entities';
import Axios, {AxiosError} from 'axios';
import {dev} from '$app/environment';
import {fetchline} from "../routes/progress/text";
import {writable} from "svelte/store";
import {messageStore} from "../store/stores";
let progress = 0;
let installationComplete = false;
let text = '';
let lines = 0;
const maxLines = 100;

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

    async installSoftware(install: Install): Promise<string> {
        if (!dev) {
            const installData = {}

            const response = await Axios.post<{ id: string }>(urls.install, installData);
            return response.data.id;
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

const token = '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';
const api = Axios.create({
    baseURL: 'https://api.github.com',
    headers: {
        'Authorization': `token ${token}`
    }
});

async function generateMarketplace(): Promise<Product[]> {

    const c = await getRepoContents("unity-sds", "unity-marketplace");

    const products: Product[] = []
    for (const p of c) {
        const content = await getGitHubFileContents("unity-sds", "unity-marketplace", p)
        const prod: Product = JSON.parse(content)
        products.push(prod)
    }

    return products
}

async function getRepoContents(user: string, repo: string, path = ''): Promise<string[]> {
    const url = `/repos/${user}/${repo}/contents/${path}`;

    const paths: string[] = [];
    try {
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
    }
];
