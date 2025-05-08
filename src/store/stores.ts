import { Order } from '../data/entities';
import { writable } from 'svelte/store';
import type {
	Install,
	Parameters,
	Deployments
} from '../data/unity-cs-manager/protobuf/extensions';

// Local definition of MarketplaceMetadata instead of using the protobuf-generated one
export interface MarketplaceMetadataStatement {
	Effect: string;
	Action: string[];
	Resource: string[];
}

export interface MarketplaceMetadataIamRoles {
	Statement: MarketplaceMetadataStatement[];
}

export interface MarketplaceMetadataVariables {
	Values: { [key: string]: string };
	AdvancedValues?: { [key: string]: any };
}

export interface MarketplaceMetadataDefaultDeployment {
	Variables?: MarketplaceMetadataVariables;
}

export interface MarketplaceMetadata {
	Name: string;
	DisplayName: string;
	Version: string;
	Channel: string;
	Owner: string;
	Description: string;
	Repository: string;
	Tags: string[];
	Category: string;
	IamRoles?: MarketplaceMetadataIamRoles;
	Package: string;
	ManagedDependencies?: { [key: string]: any };
	Backend: string;
	Entrypoint: string;
	WorkDirectory: string;
	PostInstall: string;
	PreInstall: string;
	DefaultDeployment?: MarketplaceMetadataDefaultDeployment;
	Dependencies: { [key: string]: string };
}

// Helper function to create a MarketplaceMetadata object from JSON data (replacing the protobuf fromJSON method)
export function createMarketplaceMetadataFromJSON(json: any): MarketplaceMetadata {
	return {
		Name: json.Name || '',
		DisplayName: json.DisplayName || '',
		Version: json.Version || '',
		Channel: json.Channel || '',
		Owner: json.Owner || '',
		Description: json.Description || '',
		Repository: json.Repository || '',
		Tags: Array.isArray(json.Tags) ? json.Tags.map((e: any) => String(e)) : [],
		Category: json.Category || '',
		IamRoles: json.IamRoles
			? {
					Statement: Array.isArray(json.IamRoles.Statement)
						? json.IamRoles.Statement.map((s: any) => ({
								Effect: s.Effect || '',
								Action: Array.isArray(s.Action) ? s.Action.map((a: any) => String(a)) : [],
								Resource: Array.isArray(s.Resource) ? s.Resource.map((r: any) => String(r)) : []
						  }))
						: []
			  }
			: undefined,
		Package: json.Package || '',
		ManagedDependencies:
			typeof json.ManagedDependencies === 'object' ? json.ManagedDependencies : undefined,
		Backend: json.Backend || '',
		Entrypoint: json.Entrypoint || '',
		WorkDirectory: json.WorkDirectory || '',
		PostInstall: json.PostInstall || '',
		PreInstall: json.PreInstall || '',
		DefaultDeployment: json.DefaultDeployment
			? {
					Variables: json.DefaultDeployment.Variables
						? {
								Values:
									typeof json.DefaultDeployment.Variables.Values === 'object'
										? Object.entries(json.DefaultDeployment.Variables.Values).reduce(
												(acc: { [key: string]: string }, [key, value]) => {
													acc[key] = String(value);
													return acc;
												},
												{}
										  )
										: {},
								AdvancedValues:
									typeof json.DefaultDeployment.Variables.AdvancedValues === 'object'
										? json.DefaultDeployment.Variables.AdvancedValues
										: undefined
						  }
						: undefined
			  }
			: undefined,
		Dependencies:
			typeof json.Dependencies === 'object'
				? Object.entries(json.Dependencies).reduce(
						(acc: { [key: string]: string }, [key, value]) => {
							acc[key] = String(value);
							return acc;
						},
						{}
				  )
				: {
						shared_services_account: '/unity/shared-services/aws/account',
						shared_services_region: '/unity/shared-services/aws/account/region',
						venue_proxy_baseurl: '/unity/${PROJ}/${VENUE}/management/httpd/loadbalancer-url',
						venue_subnet_list: '/unity/account/network/subnet_list'
				  }
	};
}

// Function to create an empty MarketplaceMetadata object (replacing the protobuf create method)
export function createEmptyMarketplaceMetadata(): MarketplaceMetadata {
	return {
		Name: '',
		DisplayName: '',
		Version: '',
		Channel: '',
		Owner: '',
		Description: '',
		Repository: '',
		Tags: [],
		Category: '',
		IamRoles: undefined,
		Package: '',
		ManagedDependencies: undefined,
		Backend: '',
		Entrypoint: '',
		WorkDirectory: '',
		PostInstall: '',
		PreInstall: '',
		DefaultDeployment: undefined,
		Dependencies: {}
	};
}

export interface NetworkConfig {
	publicsubnets: string[];
	privatesubnets: string[];
}

export interface ApplicationConfig {
	MarketplaceOwner: string;
	MarketplaceUser: string;
	Project: string;
	Venue: string;
	Version: string;
}

export interface Config {
	applicationConfig: ApplicationConfig;
	networkConfig: NetworkConfig;
	lastupdated: string;
	updatedby: string;
	bootstrap: string;
	version: string;
}

// Create a function to fetch the configuration from the API
export async function fetchConfigFromAPI(): Promise<Config | null> {
	try {
		const response = await fetch('../api/config');
		if (response.ok) {
			// API now directly returns data in the expected format
			return await response.json();
		}
	} catch (error) {
		console.error('Error fetching configuration:', error);
	}
	return null;
}

// Create the writable store for configuration
export const config = writable<Config | null>(null);

// Function to refresh the config store from the API
export async function refreshConfig(): Promise<void> {
	const configData = await fetchConfigFromAPI();
	if (configData) {
		config.set(configData);
	}
}

// Initialize the config - this should be called on app startup
if (typeof window !== 'undefined') {
	refreshConfig();
}

export const selectedCategory = writable<string>('');
export const order = writable<Order>(new Order());
export const install = writable<Install>({} as Install);
export const projectStore = writable('');
export const venueStore = writable('');
export const parametersStore = writable<Parameters>({} as Parameters);
export const messageStore = writable<string>('');
export const marketplaceStore = writable<MarketplaceMetadata[]>([]);

export const initialized = writable<boolean>(false);

export const isLoading = writable<boolean>(false);

export const installRunning = writable<boolean>(false);

export const installError = writable<boolean>(false);

export const productInstall = writable<MarketplaceMetadata>();

export const deploymentStore = writable<Deployments>();
