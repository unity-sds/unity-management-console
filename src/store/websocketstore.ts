import { writable, derived } from 'svelte/store';
import type { Writable } from 'svelte/store';
import type { Readable } from 'svelte/store';
import { UnityWebsocketMessage } from '../data/unity-cs-manager/protobuf/extensions';

type WritableStore = Writable<UnityWebsocketMessage[]>;

export type WebsocketStore = {
	subscribe: WritableStore['subscribe'];
	send: (message: Uint8Array) => void;
	filterByType: (type: keyof UnityWebsocketMessage) => Readable<UnityWebsocketMessage[]>;
};

function createWebsocketStore(url: string): WebsocketStore {
	const { subscribe, update } = writable<UnityWebsocketMessage[]>([]);

	let socket: WebSocket;
	let messageQueue: Uint8Array[] = [];

	function connect() {
		if (typeof WebSocket === 'undefined') {
			console.error('WebSocket is not available in this environment');
			return;
		}
		socket = new WebSocket(url);

		socket.onopen = () => {
			console.log('Socket is open');
			sendQueuedMessages();
		};

		socket.onmessage = async (event) => {
			const response = new Response(event.data);
			const arrayBuffer = await response.arrayBuffer();
			const data = new Uint8Array(arrayBuffer);
			const message = UnityWebsocketMessage.decode(data); // Decode the Protobuf message
			update((messages) => [...messages, message]);
		};

		socket.onerror = (error) => {
			console.error('WebSocket error: ', error);
		};

		socket.onclose = (event) => {
			console.log('Socket is closed. Reconnect will be attempted in 1 second.', event.reason);
			setTimeout(function () {
				connect();
			}, 1000);
		};
	}

	connect();

	function sendQueuedMessages() {
		messageQueue.forEach((message) => {
			send(message);
		});
		messageQueue = [];
	}

	function send(message: Uint8Array) {
		if (socket.readyState === WebSocket.OPEN) {
			socket.send(message);
		} else {
			console.log('Socket is not open. Queueing the message.');
			messageQueue.push(message);
		}
	}

	function filterByType(type: keyof UnityWebsocketMessage): Readable<UnityWebsocketMessage[]> {
		return derived({ subscribe }, ($messages, set) => {
			set($messages.filter((message: UnityWebsocketMessage) => message[type] !== undefined));
		});
	}

	return {
		subscribe,
		send,
		filterByType
	};
}

export { createWebsocketStore };
