package types

type ApplicationInstallParams struct {
	Name           string
	Version        string
	DisplayName string
	DeploymentName string
	Variables      map[string]string
}
