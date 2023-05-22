import type { Product } from '../data/entities';
import { Order, Install } from '../data/entities';
import { writable } from 'svelte/store';

export const products = writable<Product[]>([]);
export const selectedCategory = writable<string>('');
export const order = writable<Order>(new Order());
export const install = writable<Install>(new Install())
export const projectStore = writable('');
export const venueStore = writable('');
export const privateSubnetsStore = writable<string[]>([]);
export const publicSubnetsStore = writable<string[]>([]);
export const listStore = writable<{ key: string; value: string }[]>([]);

export const messageStore = writable<string>('')