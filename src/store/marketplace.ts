import { readable } from 'svelte/store';
import type { MarketplaceMetadata } from './stores';

const marketplaceowner = 'unity-sds';
const marketplacerepo = 'unity-marketplace';

export const marketplaceData = readable<MarketplaceMetadata[]>([], (set) => {
	const url = `https://api.github.com/repos/${marketplaceowner}/${marketplacerepo}/contents/manifest.json`;
	fetch(url)
		.then((res) => res.json())
		.then((json) => {
			// GitHub API returns content as base64 encoded
			if (json.content) {
				// Decode the base64 content
				const decodedContent = atob(json.content.replace(/\n/g, ''));
				// Parse the JSON content
				try {
					const parsedContent = JSON.parse(decodedContent);
					set(parsedContent);
				} catch (e) {
					console.error('Error parsing marketplace data:', e);
					set([]);
				}
			} else {
				console.error('No content found in GitHub response');
				set([]);
			}
		})
		.catch((error) => {
			console.error('Error fetching marketplace data:', error);
			set([]);
		});
});

type InstalledMarketplaceApplication = {
	DeploymentName: string;
	PackageName: string;
	Name: string;
	Source: string;
	Version: string;
	Status: string;
};

export const installedApplications = readable<InstalledMarketplaceApplication[]>([], (set) => {
	fetch('../api/installed_applications')
		.then((res) => res.json())
		.then((json) => set(json))
		.catch((e) => {
			console.warn('Unable to get application list!');
		});
});
