package types

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