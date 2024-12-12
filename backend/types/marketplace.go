package types


type SSMProviderApp struct {
	MinimumVersion string `json:MinimumVersion`
}

type SSMParam map[string]string

type Dependencies struct {
	SSMParams []SSMParam `json:SSMParams`
	Apps 	  map[string]SSMProviderApp `json:Apps`
}

type MarketplaceMetadata struct {
	Dependencies Dependencies `json:Dependencies`
}

type AdvancedValue map[string]interface{}

type ApplicationInstallParams struct {
	Name           string
	Version        string
	DisplayName    string
	DeploymentName string
	Variables      map[string]string
	AdvancedValues AdvancedValue
}


type InstalledMarketplaceApplication struct {
	Name         string
	DeploymentName  string
	Version      string
	Source       string
	Status       string
	PackageName  string	
	TerraformModuleName string
	Variables    map[string]string
	AdvancedValues    AdvancedValue
}
