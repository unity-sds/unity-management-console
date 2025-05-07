import { browser } from '$app/environment';
import type { PageLoad } from './$types';
import { HttpHandler } from '../../data/httpHandler';
import { get } from 'svelte/store';
import { marketplaceStore } from '../../store/stores';

export const load: PageLoad = async ({ url }) => {
  // Only run in the browser
  if (!browser) {
    return { success: true };
  }

  const name = url.searchParams.get('name');
  const version = url.searchParams.get('version');

  // If we have the parameters and no marketplace data, fetch it
  if (name && version && get(marketplaceStore).length === 0) {
    const httpHandler = new HttpHandler();
    await httpHandler.fetchConfig();
  }

  return {
    name,
    version,
    success: true
  };
};