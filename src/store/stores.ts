import { Order } from '../data/entities';
import { writable } from 'svelte/store';
import type {
	Config,
	Install,
	Parameters,
	Deployments
} from '../data/unity-cs-manager/protobuf/extensions';
import type { MarketplaceMetadata } from '../data/unity-cs-manager/protobuf/marketplace';

export const config = writable<Config | null>(null);
export const selectedCategory = writable<string>('');
export const order = writable<Order>(new Order());
export const install = writable<Install>({} as Install);
export const projectStore = writable('');
export const venueStore = writable('');
export const parametersStore = writable<Parameters>({} as Parameters);
export const messageStore = writable<string>('');
export const marketplaceStore = writable<MarketplaceMetadata[]>([]);

export const initialized = writable<boolean>(false);

export const isLoading = writable<boolean>(false);

export const installRunning = writable<boolean>(false);

export const installError = writable<boolean>(false);

export const productInstall = writable<MarketplaceMetadata>();

export const deploymentStore = writable<Deployments>();
