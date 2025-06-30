package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/gitutil"
)

// ModuleResolver handles module resolution and caching
type ModuleResolver struct {
	registry   *ModuleRegistry
	cacheDir   string
	appConfig  *config.AppConfig
	mu         sync.RWMutex
	httpClient *gitutil.Downloader
}

// NewModuleResolver creates a new module resolver instance
func NewModuleResolver(appConfig *config.AppConfig) (*ModuleResolver, error) {
	cacheDir := filepath.Join(appConfig.Workdir, "module_cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create module cache directory: %w", err)
	}

	resolver := &ModuleResolver{
		cacheDir:   cacheDir,
		appConfig:  appConfig,
		httpClient: gitutil.NewDownloader(),
	}

	// Load registry if it exists
	if err := resolver.LoadRegistry(); err != nil {
		// Registry not existing is not an error - feature is optional
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load module registry: %w", err)
		}
	}

	return resolver, nil
}

// LoadRegistry loads the module registry from the configured location
func (r *ModuleResolver) LoadRegistry() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// First, check for local registry file
	localRegistryPath := filepath.Join(r.appConfig.Workdir, "module-registry.json")
	if _, err := os.Stat(localRegistryPath); err == nil {
		return r.loadRegistryFromFile(localRegistryPath)
	}

	// If not found locally, try to download from marketplace
	registryURL := r.buildRegistryURL()
	if registryURL != "" {
		return r.downloadAndLoadRegistry(registryURL)
	}

	return os.ErrNotExist
}

// buildRegistryURL constructs the URL to download the registry from
func (r *ModuleResolver) buildRegistryURL() string {
	if r.appConfig.MarketplaceBaseUrl == "" || r.appConfig.MarketplaceOwner == "" || r.appConfig.MarketplaceRepo == "" {
		return ""
	}

	// Construct URL similar to marketplace pattern
	return fmt.Sprintf("%s%s/%s/main/module-registry.json",
		r.appConfig.MarketplaceBaseUrl,
		r.appConfig.MarketplaceOwner,
		r.appConfig.MarketplaceRepo,
	)
}

// downloadAndLoadRegistry downloads and loads the registry from a URL
func (r *ModuleResolver) downloadAndLoadRegistry(url string) error {
	tempFile := filepath.Join(r.cacheDir, "module-registry-download.json")
	
	// Download registry file
	if err := r.httpClient.DownloadFile(url, tempFile); err != nil {
		return fmt.Errorf("failed to download registry: %w", err)
	}
	defer os.Remove(tempFile)

	// Load from downloaded file
	if err := r.loadRegistryFromFile(tempFile); err != nil {
		return err
	}

	// Cache the registry locally
	localPath := filepath.Join(r.appConfig.Workdir, "module-registry.json")
	data, err := os.ReadFile(tempFile)
	if err != nil {
		return err
	}
	
	return os.WriteFile(localPath, data, 0644)
}

// loadRegistryFromFile loads the registry from a file
func (r *ModuleResolver) loadRegistryFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read registry file: %w", err)
	}

	var registry ModuleRegistry
	if err := json.Unmarshal(data, &registry); err != nil {
		return fmt.Errorf("failed to parse registry file: %w", err)
	}

	r.registry = &registry
	return nil
}

// HasRegistry returns true if a module registry is loaded
func (r *ModuleResolver) HasRegistry() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.registry != nil
}

// ResolveModule resolves a module reference to its full definition
func (r *ModuleResolver) ResolveModule(ref *ModuleReference) (*ResolvedModule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.registry == nil {
		return nil, fmt.Errorf("no module registry loaded")
	}

	module, exists := r.registry.Modules[ref.Name]
	if !exists {
		return nil, fmt.Errorf("module '%s' not found in registry", ref.Name)
	}

	// Resolve version
	version := ref.Version
	if version == "" {
		version = "latest"
	}

	versionInfo, exists := module.Versions[version]
	if !exists {
		return nil, fmt.Errorf("version '%s' not found for module '%s'", version, ref.Name)
	}

	// Build resolved module
	resolved := &ResolvedModule{
		ModuleReference: *ref,
		Definition:      module,
		VersionInfo:     versionInfo,
	}

	// Generate cache path
	cacheKey := fmt.Sprintf("%s-%s", ref.Name, version)
	resolved.CachePath = filepath.Join(r.cacheDir, cacheKey)

	// Build source URL
	if strings.Contains(module.Source, "//") {
		// Already has subdir notation
		resolved.ResolvedSource = fmt.Sprintf("git::https://%s?ref=%s", module.Source, versionInfo.Ref)
	} else {
		resolved.ResolvedSource = fmt.Sprintf("git::https://%s?ref=%s", module.Source, versionInfo.Ref)
	}

	return resolved, nil
}

// CacheModule downloads and caches a module locally
func (r *ModuleResolver) CacheModule(resolved *ResolvedModule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if already cached
	if _, err := os.Stat(resolved.CachePath); err == nil {
		return nil // Already cached
	}

	// Create temp directory for download
	tempDir := filepath.Join(r.cacheDir, fmt.Sprintf("tmp-%d", time.Now().UnixNano()))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Clone the module
	if err := gitutil.CloneRepo(resolved.ResolvedSource, tempDir); err != nil {
		return fmt.Errorf("failed to clone module: %w", err)
	}

	// Move to cache location
	if err := os.Rename(tempDir, resolved.CachePath); err != nil {
		return fmt.Errorf("failed to cache module: %w", err)
	}

	return nil
}

// GetCachedModulePath returns the local path to a cached module
func (r *ModuleResolver) GetCachedModulePath(moduleName, version string) string {
	cacheKey := fmt.Sprintf("%s-%s", moduleName, version)
	return filepath.Join(r.cacheDir, cacheKey)
}

// ValidateModuleConfig validates that required inputs are provided
func (r *ModuleResolver) ValidateModuleConfig(resolved *ResolvedModule) error {
	if resolved.Definition.Inputs == nil {
		return nil
	}

	for inputName, input := range resolved.Definition.Inputs {
		if input.Required {
			if _, exists := resolved.Config[inputName]; !exists {
				if input.Default == nil {
					return fmt.Errorf("required input '%s' not provided for module '%s'", inputName, resolved.Name)
				}
			}
		}
	}

	return nil
}

// GetRegistry returns the loaded registry (for display purposes)
func (r *ModuleResolver) GetRegistry() *ModuleRegistry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.registry
}

// RefreshRegistry forces a refresh of the module registry
func (r *ModuleResolver) RefreshRegistry() error {
	r.mu.Lock()
	r.registry = nil
	r.mu.Unlock()
	
	return r.LoadRegistry()
}