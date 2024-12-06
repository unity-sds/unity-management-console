export type InstalledMarketplaceApplication = {
	DeploymentName: string;
	PackageName: string;
	Name: string;
	Source: string;
	Version: string;
	Status: string;
};

export async function getInstalledApplications(): Promise<InstalledMarketplaceApplication[]> {
	const res = await fetch('../api/installed_applications');
	if (!res.ok) {
		console.warn('Unable to get application list!');
		return [];
	}
	return await res.json();
}
