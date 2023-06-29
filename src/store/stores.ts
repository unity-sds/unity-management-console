import type {Product } from "../data/entities";
import type {Config} from "../data/unity-cs-manager/protobuf/config"
import { Order } from '../data/entities';
import { writable } from 'svelte/store';
import type { Install } from "../data/unity-cs-manager/protobuf/extensions";
import type { Parameters} from "../data/unity-cs-manager/protobuf/config";

export const config = writable<Config | null>(null);
export const products = writable<Product[]>([]);
export const selectedCategory = writable<string>('');
export const order = writable<Order>(new Order());
export const install = writable<Install>({} as Install)
export const projectStore = writable('');
export const venueStore = writable('');

export const parametersStore = writable<Parameters>({} as Parameters)
export const privateSubnetsStore = writable<string[]>([]);
export const publicSubnetsStore = writable<string[]>([]);
export const listStore = writable<{ key: string; value: string }[]>([]);

export const messageStore = writable<string>('')

export const installComplete = writable<boolean>(false)