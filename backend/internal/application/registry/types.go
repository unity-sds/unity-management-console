package registry

import "time"

// ModuleRegistry represents the entire module registry catalog
type ModuleRegistry struct {
	Version  string                     `json:"version"`
	Metadata ModuleRegistryMetadata     `json:"metadata"`
	Modules  map[string]*ModuleDefinition `json:"modules"`
}

// ModuleRegistryMetadata contains metadata about the registry
type ModuleRegistryMetadata struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	LastUpdated time.Time `json:"last_updated"`
}

// ModuleDefinition defines a Terraform module available in the registry
type ModuleDefinition struct {
	Description   string                        `json:"description"`
	Source        string                        `json:"source"`
	Documentation string                        `json:"documentation,omitempty"`
	Versions      map[string]*ModuleVersionInfo `json:"versions"`
	Inputs        map[string]*ModuleInput       `json:"inputs,omitempty"`
	Outputs       map[string]string             `json:"outputs,omitempty"`
}

// ModuleVersionInfo contains version-specific information
type ModuleVersionInfo struct {
	Ref               string            `json:"ref"`
	TerraformVersion  string            `json:"terraform_version,omitempty"`
	Providers         map[string]string `json:"providers,omitempty"`
	MinManagementConsoleVersion string            `json:"min_mc_version,omitempty"`
}

// ModuleInput describes an input variable for a module
type ModuleInput struct {
	Description string      `json:"description"`
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
	Sensitive   bool        `json:"sensitive,omitempty"`
}

// ModuleReference represents a user's reference to a module in their configuration
type ModuleReference struct {
	Name      string                 `json:"name"`
	Version   string                 `json:"version"`
	Alias     string                 `json:"alias,omitempty"`
	Config    map[string]interface{} `json:"config"`
	DependsOn []string               `json:"depends_on,omitempty"`
}

// ResolvedModule contains the resolved module information ready for Terraform
type ResolvedModule struct {
	ModuleReference
	Definition   *ModuleDefinition
	VersionInfo  *ModuleVersionInfo
	ResolvedSource string
	CachePath    string
}