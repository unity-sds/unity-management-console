import { browser } from '$app/environment';
import type { PageLoad } from './$types';
import { get } from 'svelte/store';
import { marketplaceStore, refreshConfig } from '../../store/stores';

export const load: PageLoad = async ({ url }) => {
  // Only run in the browser
  if (!browser) {
    return { success: true };
  }

  const name = url.searchParams.get('name');
  const version = url.searchParams.get('version');
  
  // We'll fetch the data if needed in the component's onMount
  // This is to prevent potential loops from both the load function
  // and the component initializing data

  return {
    name,
    version,
    hasMarketplaceData: get(marketplaceStore).length > 0,
    success: true
  };
};