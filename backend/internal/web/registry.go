package web

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/application/registry"
)

// TerraformRegistryDiscovery represents the service discovery response
type TerraformRegistryDiscovery struct {
	ModulesV1 string `json:"modules.v1"`
}

// ModuleVersionsResponse represents the module versions API response
type ModuleVersionsResponse struct {
	Modules []ModuleVersionInfo `json:"modules"`
}

type ModuleVersionInfo struct {
	Versions []ModuleVersion `json:"versions"`
}

type ModuleVersion struct {
	Version    string      `json:"version"`
	Submodules []string    `json:"submodules"`
	Root       ModuleRoot  `json:"root"`
}

type ModuleRoot struct {
	Providers    []ModuleProvider `json:"providers"`
	Dependencies []string         `json:"dependencies"`
	Variables    []ModuleVariable `json:"variables"`
	Outputs      []ModuleOutput   `json:"outputs"`
}

type ModuleProvider struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Source    string `json:"source"`
	Version   string `json:"version"`
}

type ModuleVariable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type ModuleOutput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// handleRegistryDiscovery handles the Terraform registry service discovery
func handleRegistryDiscovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := TerraformRegistryDiscovery{
			ModulesV1: "/v1/modules/",
		}
		c.JSON(http.StatusOK, response)
	}
}

// handleModuleVersions handles the module versions API endpoint
func handleModuleVersions(appConfig config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		provider := c.Param("provider")

		log.Infof("Fetching module versions for %s/%s/%s", namespace, name, provider)

		// Initialize module resolver
		resolver, err := registry.NewModuleResolver(&appConfig)
		if err != nil {
			log.WithError(err).Error("Failed to initialize module resolver")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize module resolver"})
			return
		}

		if !resolver.HasRegistry() {
			c.JSON(http.StatusNotFound, gin.H{"error": "No module registry configured"})
			return
		}

		// Get the registry
		moduleRegistry := resolver.GetRegistry()
		if moduleRegistry == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Module registry not loaded"})
			return
		}

		// Find the module
		fullModuleName := fmt.Sprintf("%s-%s", namespace, name)
		module, exists := moduleRegistry.Modules[fullModuleName]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Module %s not found", fullModuleName)})
			return
		}

		// Build version information
		var versions []ModuleVersion
		for versionStr, _ := range module.Versions {
			moduleVersion := ModuleVersion{
				Version:    versionStr,
				Submodules: []string{},
				Root: ModuleRoot{
					Providers:    []ModuleProvider{},
					Dependencies: []string{},
					Variables:    []ModuleVariable{},
					Outputs:      []ModuleOutput{},
				},
			}

			// Provider information comes from module's versions.tf, not registry metadata

			// Add inputs as variables
			if module.Inputs != nil {
				for inputName, input := range module.Inputs {
					moduleVersion.Root.Variables = append(moduleVersion.Root.Variables, ModuleVariable{
						Name:        inputName,
						Description: input.Description,
						Type:        input.Type,
					})
				}
			}

			// Add outputs
			if module.Outputs != nil {
				for outputName, outputDescription := range module.Outputs {
					moduleVersion.Root.Outputs = append(moduleVersion.Root.Outputs, ModuleOutput{
						Name:        outputName,
						Description: outputDescription,
					})
				}
			}

			versions = append(versions, moduleVersion)
		}

		response := ModuleVersionsResponse{
			Modules: []ModuleVersionInfo{
				{
					Versions: versions,
				},
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

// handleModuleDownload handles the module download endpoint
func handleModuleDownload(appConfig config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		provider := c.Param("provider")
		version := c.Param("version")

		log.Infof("Handling module download for %s/%s/%s@%s", namespace, name, provider, version)

		// Set the X-Terraform-Get header pointing to the archive endpoint
		archiveURL := fmt.Sprintf("/v1/modules/%s/%s/%s/%s/archive", namespace, name, provider, version)
		c.Header("X-Terraform-Get", archiveURL)
		c.Status(http.StatusNoContent)
	}
}

// handleModuleArchive handles serving the module archive
func handleModuleArchive(appConfig config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")
		provider := c.Param("provider")
		version := c.Param("version")

		log.Infof("Serving module archive for %s/%s/%s@%s", namespace, name, provider, version)

		// Initialize module resolver
		resolver, err := registry.NewModuleResolver(&appConfig)
		if err != nil {
			log.WithError(err).Error("Failed to initialize module resolver")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize module resolver"})
			return
		}

		// Create module reference
		fullModuleName := fmt.Sprintf("%s-%s", namespace, name)
		moduleRef := &registry.ModuleReference{
			Name:    fullModuleName,
			Version: version,
			Config:  make(map[string]interface{}),
		}

		// Resolve the module
		resolved, err := resolver.ResolveModule(moduleRef)
		if err != nil {
			log.WithError(err).Errorf("Failed to resolve module %s", fullModuleName)
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Module %s not found", fullModuleName)})
			return
		}

		// Cache the module
		err = resolver.CacheModule(resolved)
		if err != nil {
			log.WithError(err).Error("Failed to cache module")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cache module"})
			return
		}

		// Create tar.gz archive from cached module
		archivePath, err := createModuleArchive(resolved.CachePath)
		if err != nil {
			log.WithError(err).Error("Failed to create module archive")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create module archive"})
			return
		}
		defer os.Remove(archivePath)

		// Serve the archive
		c.Header("Content-Type", "application/gzip")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s-%s-%s.tar.gz", namespace, name, version))
		c.File(archivePath)
	}
}

// createModuleArchive creates a tar.gz archive from the module directory
func createModuleArchive(modulePath string) (string, error) {
	// Create temporary file for archive
	tempFile, err := os.CreateTemp("", "module-*.tar.gz")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	// Create gzip writer
	gzipWriter := gzip.NewWriter(tempFile)
	defer gzipWriter.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Walk the module directory and add files to archive
	err = filepath.Walk(modulePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and hidden files
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// Create relative path from module root
		relPath, err := filepath.Rel(modulePath, path)
		if err != nil {
			return err
		}

		// Create tar header
		header := &tar.Header{
			Name: relPath,
			Size: info.Size(),
			Mode: int64(info.Mode()),
		}

		// Write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// Write file content
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tarWriter, file)
		return err
	})

	if err != nil {
		return "", fmt.Errorf("failed to create archive: %w", err)
	}

	return tempFile.Name(), nil
}