import { createWebsocketStore } from '../store/websocketstore';
import type { WebsocketStore } from '../store/websocketstore';

// Define the WebSocket URL
const url: string =
	typeof window !== 'undefined'
		? `ws${window.location.protocol === 'https:' ? 's' : ''}://${
				window.location.host
		  }` + "/management/ws"
		: 'ws://localhost:8080/ws';

// Create the WebSocket store using the provided function
export const websocketStore: WebsocketStore = createWebsocketStore(url);
