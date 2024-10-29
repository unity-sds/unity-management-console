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
