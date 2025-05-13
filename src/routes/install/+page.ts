import { browser } from '$app/environment';
import type { PageLoad } from './$types';
import { get } from 'svelte/store';
import { marketplaceStore, refreshConfig } from '../../store/stores';

export const load: PageLoad = async ({ url, fetch }) => {
  // Only run in the browser
  if (!browser) {
    return { success: true };
  }

  const name = url.searchParams.get('name');
  const version = url.searchParams.get('version');
  
  let product = null;
  
  // If we have marketplace data already, use it
  if (get(marketplaceStore).length > 0) {
    product = get(marketplaceStore).find(
      (p) => p.Name === name && p.Version === version
    );
  } 
  // Otherwise fetch just this specific product
  else if (name && version) {
    try {
      const response = await fetch(`/api/marketplace/item/${name}/${version}`);
      if (response.ok) {
        product = await response.json();
      }
    } catch (error) {
      console.error("Failed to fetch product:", error);
    }
  }

  return {
    name,
    version,
    product,
    hasMarketplaceData: get(marketplaceStore).length > 0,
    success: true
  };
};
