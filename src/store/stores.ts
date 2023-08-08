import type {Product } from "../data/entities";
import { Order } from '../data/entities';
import { writable } from 'svelte/store';
import type { Config, Install, Parameters } from "../data/unity-cs-manager/protobuf/extensions";

export const config = writable<Config | null>(null);
export const selectedCategory = writable<string>('');
export const order = writable<Order>(new Order());
export const install = writable<Install>({} as Install)
export const projectStore = writable('');
export const venueStore = writable('');
export const parametersStore = writable<Parameters>({} as Parameters)
export const messageStore = writable<string>('')
export const marketplaceStore = writable<Product[]>([])

export const initialized = writable<boolean>(false)

export const installRunning = writable<boolean>(false)

export const installError = writable<boolean>(false)