package types

type AdvancedValue map[string]interface{}

type ApplicationInstallParams struct {
	Name           string
	Version        string
	DisplayName    string
	DeploymentName string
	Variables      map[string]string
	AdvancedValues AdvancedValue
	ModuleReferences []ModuleReference `json:"moduleReferences,omitempty"`
}

// ModuleReference represents a reference to a module in the registry
type ModuleReference struct {
	Name      string                 `json:"name"`
	Version   string                 `json:"version"`
	Alias     string                 `json:"alias,omitempty"`
	Config    map[string]interface{} `json:"config"`
	DependsOn []string               `json:"depends_on,omitempty"`
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
	ModuleReferences []ModuleReference `json:"moduleReferences,omitempty"`
}