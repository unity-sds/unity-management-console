import type {Product } from "../data/entities";
import type {Config} from "../data/protobuf/config"
import { Order } from '../data/entities';
import { writable } from 'svelte/store';
import type { Install } from "../data/protobuf/extensions";

export const config = writable<Config | null>(null);
export const products = writable<Product[]>([]);
export const selectedCategory = writable<string>('');
export const order = writable<Order>(new Order());
export const install = writable<Install>({} as Install)
export const projectStore = writable('');
export const venueStore = writable('');
export const privateSubnetsStore = writable<string[]>([]);
export const publicSubnetsStore = writable<string[]>([]);
export const listStore = writable<{ key: string; value: string }[]>([]);

export const messageStore = writable<string>('')