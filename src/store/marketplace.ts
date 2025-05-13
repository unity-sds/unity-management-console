import { readable } from 'svelte/store';

const marketplaceowner = 'unity-sds';
const marketplacerepo = 'unity-marketplace';

export const marketplaceData = readable({}, (set) => {
	const url = `https://api.github.com/repos/${marketplaceowner}/${marketplacerepo}/contents/manifest.json`;
	fetch(url)
		.then((res) => res.json())
		.then((json) => set(json))
		.catch((e) => set({}));
});

// async function generateMarketplace() {
// 	if (!dev) {
// 		console.log('Checking if manifest.json exists in the repository...');
// 		const manifestExists = await checkIfFileExists(
// 			marketplaceowner,
// 			marketplacerepo,
// 			'manifest.json'
// 		);
// 		if (manifestExists) {
// 			console.log('manifest.json exists in the repository.');
// 			const content = await getGitHubFileContents(
// 				marketplaceowner,
// 				marketplacerepo,
// 				'manifest.json'
// 			);
// 			const c = JSON.parse(content);
// 			const products: MarketplaceMetadata[] = [];
// 			for (const p of c) {
// 				const prod = MarketplaceMetadata.fromJSON(p);
// 				products.push(prod);
// 			}
// 			marketplaceStore.set(products);
// 			return;
// 		}

// 		console.log('fetching repo contents: ' + marketplaceowner);
// 		const c = await getRepoContents(marketplaceowner, marketplacerepo);

// 		const products: MarketplaceMetadata[] = [];
// 		for (const p of c) {
// 			const content = await getGitHubFileContents(marketplaceowner, marketplacerepo, p);
// 			const j = JSON.parse(content);
// 			const prod = MarketplaceMetadata.fromJSON(j);
// 			products.push(prod);
// 		}

// 		marketplaceStore.set(products);
// 	} else {
// 		const j = JSON.parse(mock_marketplace);
// 		const products: MarketplaceMetadata[] = [];
// 		for (const p of j) {
// 			const prod = MarketplaceMetadata.fromJSON(p);
// 			products.push(prod);
// 		}
// 		marketplaceStore.set(products);
// 		return;
// 	}
// }
