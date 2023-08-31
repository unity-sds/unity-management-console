import { createWebsocketStore } from '../store/websocketstore';
import type { WebsocketStore } from '../store/websocketstore';

// ... (rest of your existing code)

// Define the WebSocket URL
//const url: string = 'ws://' + window.location.host + '/ws';
const url: string =
	typeof window !== 'undefined' ? `ws://${window.location.host}/ws` : 'ws://localhost/ws';

// Create the WebSocket store using the provided function
export const websocketStore: WebsocketStore = createWebsocketStore(url);
