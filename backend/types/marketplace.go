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

type MarketplaceMetadata struct {
	Name                string                 `json:"name"`
	DisplayName         string                 `json:"displayName"`
	Description         string                 `json:"description"`
	Version             string                 `json:"version"`
	Category            string                 `json:"category"`
	ShortDescription    string                 `json:"shortDescription"`
	TerraformModuleName string                 `json:"terraformModuleName"`
	Variables           []MarketplaceVariable  `json:"variables"`
	AdvancedVariables   []MarketplaceVariable  `json:"advancedVariables"`
	Outputs             []MarketplaceOutput    `json:"outputs"`
	Screenshots         []MarketplaceScreenshot `json:"screenshots"`
	OutputSsmParameters []string               `json:"OutputSsmParameters"`
	PreInstall          string                 `json:"preInstall"`
	PostInstall         string                 `json:"postInstall"`
	Package             string                 `json:"package"`
	Dependencies		map[string]string	   `json:"Dependencies"`
}

type MarketplaceVariable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Default     string `json:"default"`
	Required    bool   `json:"required"`
}

type MarketplaceOutput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MarketplaceScreenshot struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}
