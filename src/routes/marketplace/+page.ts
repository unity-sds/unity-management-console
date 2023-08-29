import type { PageLoad } from './$types';
import { HttpHandler } from '../../data/httpHandler';
import { selectedCategory, order } from '../../store/stores';
import { Order } from '../../data/entities';

export const ssr = false;
export const load = (async () => {
	order.set(new Order());
	selectedCategory.set('All');
}) satisfies PageLoad;
