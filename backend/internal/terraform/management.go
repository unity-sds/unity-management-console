package terraform

import (
	"fmt"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/google/uuid"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/application/registry"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"github.com/unity-sds/unity-management-console/backend/types"
	"github.com/zclconf/go-cty/cty"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Varstruct struct {
	Name      string
	Value     interface{}
	Tfobjtype cty.Type
}

func convertValue(v *structpb.Value) interface{} {
	switch k := v.Kind.(type) {
	case *structpb.Value_NullValue:
		return nil
	case *structpb.Value_NumberValue:
		return k.NumberValue
	case *structpb.Value_StringValue:
		return k.StringValue
	case *structpb.Value_BoolValue:
		return k.BoolValue
	case *structpb.Value_StructValue:
		return convertStruct(k.StructValue)
	case *structpb.Value_ListValue:
		arr := make([]interface{}, len(k.ListValue.Values))
		for i, lv := range k.ListValue.Values {
			arr[i] = convertValue(lv)
		}
		return arr
	}
	return nil
}

func convertStruct(s *structpb.Struct) map[string]interface{} {
	m := make(map[string]interface{})
	for key, value := range s.Fields {
		m[key] = convertValue(value)
	}

	return m
}

func convertToCty(data interface{}) cty.Value {
	switch v := data.(type) {
	case map[string]interface{}:
		ctyMap := make(map[string]cty.Value)
		for k, val := range v {
			ctyMap[k] = convertToCty(val)
		}
		return cty.ObjectVal(ctyMap)
	case []interface{}:
		ctyList := make([]cty.Value, len(v))
		for i, val := range v {
			ctyList[i] = convertToCty(val)
		}
		return cty.ListVal(ctyList)
	case string:
		return cty.StringVal(v)
	case bool:
		return cty.BoolVal(v)
	case float64: // JSON numbers are float64
		return cty.NumberFloatVal(v)
	}

	// Return a zero value if none of the above cases match
	return cty.NilVal
}

func parseAdvancedVariables(advancedVars *types.AdvancedValue, cloudenv *map[string]cty.Value) {
	if advancedVars == nil {
		return
	}

	for key, value := range *advancedVars {
		switch v := value.(type) {
		case types.AdvancedValue:
			// Create a new map for nested structure
			nestedMap := make(map[string]cty.Value)
			// Recursively process nested AdvancedValue
			parseAdvancedVariables(&v, &nestedMap)
			// Add the nested map to cloudenv
			(*cloudenv)[key] = cty.ObjectVal(nestedMap)
		default:
			// Convert the value to cty.Value and add to cloudenv
			(*cloudenv)[key] = convertToCty(v)
		}
	}
}

func generateMetadataHeader(cloudenv *hclwrite.Body, id string, application string, applicationName string, version string, creator string) {
	currentTime := time.Now()
	dateString := currentTime.Format("2006-01-02")
	comment := hclwrite.Tokens{
		&hclwrite.Token{
			Type:         hclsyntax.TokenComment,
			Bytes:        []byte(fmt.Sprintf("# id: %v\n# application: %v\n# applicationName: %v\n# version: %v\n# creator: %v\n# creationDate: %v\n", id, application, applicationName, version, creator, dateString)),
			SpacesBefore: 0,
		},
	}
	cloudenv.AppendUnstructuredTokens(comment)
}

func GenerateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// createFile creates a file with the given filename in the specified directory.
// If the directory does not exist, it will be created with the given permissions.
func createFile(directory string, filename string, perm os.FileMode) (*os.File, error) {
	// Ensure the directory exists, creating it if necessary
	err := os.MkdirAll(directory, perm)
	if err != nil {
		return nil, fmt.Errorf("could not create directory %s: %w", directory, err)
	}

	// Create the file within the directory
	filePath := filepath.Join(directory, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not create file %s: %w", filePath, err)
	}

	return file, nil
}

// AppendBlockToBody appends a new block to an HCL body with given type, labels, source, and attributes.
func appendBlockToBody(body *hclwrite.Body, blockType string, labels []string, source string, attributes map[string]cty.Value) {
	block := body.AppendNewBlock(blockType, labels)
	blockBody := block.Body()
	blockBody.SetAttributeValue("source", cty.StringVal(source))

	// Iterate through the attributes map and set each attribute in the block
	for key, value := range attributes {
		blockBody.SetAttributeValue(key, value)
	}
}

// AddApplicationToStack adds the given application configuration to the stack.
// It takes care of creating the necessary workspace directory, generating the
// HCL file, and writing the required attributes.
func AddApplicationToStack(appConfig *config.AppConfig, location string, meta *marketplace.MarketplaceMetadata, application *types.InstalledMarketplaceApplication, db database.Datastore) error {
	log.Infof("Adding application to stack. Location: %v, meta %v, install: %v, module ID: %s", location, meta, application.TerraformModuleName)
	rand.Seed(time.Now().UnixNano())

	// s := GenerateRandomString(8)
	hclFile := hclwrite.NewEmptyFile()

	directory := filepath.Join(appConfig.Workdir, "workspace")
	log.Errorf("Application name: %s", application.Name)
	filename := fmt.Sprintf("%v-%v.tf", application.Name, application.DeploymentName)

	// Handle module references if present
	if len(application.ModuleReferences) > 0 {
		log.Info("Application has module references, initializing module resolver")
		resolver, err := registry.NewModuleResolver(appConfig)
		if err != nil {
			log.WithError(err).Warn("Failed to initialize module resolver, continuing without modules")
		} else if resolver.HasRegistry() {
			// Add modules to the Terraform configuration
			if err := AddModulesToTerraform(appConfig, application.ModuleReferences, resolver, hclFile, directory); err != nil {
				log.WithError(err).Error("Failed to add modules to Terraform configuration")
				return err
			}
			log.Infof("Successfully added %d module references", len(application.ModuleReferences))
		} else {
			log.Warn("Module references found but no registry available, skipping modules")
		}
	}

	log.Errorf("Creating file with the name: %s", filename)
	tfFile, err := createFile(directory, filename, 0755)
	if err != nil {
		log.WithError(err).Error("Problem creating tf file")
		return err
	}

	path := filepath.Join(location, meta.WorkDirectory)
	// initialize the body of the new file object
	rootBody := hclFile.Body()

	u, err := uuid.NewRandom()
	if err != nil {
		log.WithError(err).Error("Failed to generate UUID")
		return err
	}

	log.Info("Generating header")
	generateMetadataHeader(rootBody, u.String(), meta.Name, application.DeploymentName, application.Version, "admin")

	log.Info("adding attributes")
	attributes := map[string]cty.Value{
		"deployment_name": cty.StringVal(application.DeploymentName),
		"tags":            cty.MapValEmpty(cty.String), // Example of setting an empty map
		"project":         cty.StringVal(appConfig.Project),
		"venue":           cty.StringVal(appConfig.Venue),
		"installprefix":   cty.StringVal(appConfig.InstallPrefix),
	}

	log.Info("Organising variable replacement")
	if application.Variables != nil {
		for key, element := range application.Variables {
			log.Infof("Adding variable: %s, %s", key, element)
			attributes[key] = cty.StringVal(element)
		}
	}
	log.Info("Parsing advanced vars")
	parseAdvancedVariables(&application.AdvancedValues, &attributes)
	rand.Seed(time.Now().UnixNano())
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	randomChars := make([]byte, 5)
	for i, v := range rand.Perm(52)[:5] {
		randomChars[i] = chars[v]
	}
	log.Info("Appending block to body")
	appendBlockToBody(rootBody, "module", []string{application.TerraformModuleName}, path, attributes)

	log.Info("Writing hcl file.")
	_, err = tfFile.Write(hclFile.Bytes())
	if err != nil {
		log.WithError(err).Error("error writing hcl file")
		return err
	}

	return nil
}

func lookUpVariablePointer(element string, inst *marketplace.Install) (string, error) {
	val, err := lookUpFromDependencies(element, inst.Applications)
	if err != nil {
		return "", err
	}
	if val != "" {
		return val, err
	}

	return "", nil
}


func lookUpFromDependencies(element string, inst *marketplace.Install_Applications) (string, error) {
	deps := inst.Dependencies
	for k, v := range deps {
		//Stripping * from string
		newStr := strings.Replace(element, "*", "", 1)
		log.Infof("Checking dependency: %s, %s for value %s", k, v, element)
		if k == newStr {
			return v, nil
		}
	}

	return "", nil

}

// AddModulesToTerraform adds module references to the Terraform configuration
func AddModulesToTerraform(appConfig *config.AppConfig, moduleRefs []types.ModuleReference, resolver *registry.ModuleResolver, tfFile *hclwrite.File, terraformDir string) error {
	if len(moduleRefs) == 0 {
		return nil
	}

	rootBody := tfFile.Body()

	for _, moduleRef := range moduleRefs {
		// Resolve the module
		resolvedModule, err := resolver.ResolveModule(&registry.ModuleReference{
			Name:      moduleRef.Name,
			Version:   moduleRef.Version,
			Alias:     moduleRef.Alias,
			Config:    moduleRef.Config,
			DependsOn: moduleRef.DependsOn,
		})
		if err != nil {
			return fmt.Errorf("failed to resolve module %s: %w", moduleRef.Name, err)
		}

		// Validate module configuration
		if err := resolver.ValidateModuleConfig(resolvedModule); err != nil {
			return fmt.Errorf("invalid configuration for module %s: %w", moduleRef.Name, err)
		}

		// Cache the module locally
		if err := resolver.CacheModule(resolvedModule); err != nil {
			return fmt.Errorf("failed to cache module %s: %w", moduleRef.Name, err)
		}

		// Determine module name for Terraform
		moduleName := moduleRef.Alias
		if moduleName == "" {
			moduleName = fmt.Sprintf("%s_%s", moduleRef.Name, GenerateRandomString(5))
		}

		// Convert config to cty.Value attributes
		attributes := map[string]cty.Value{
			"source": cty.StringVal(resolvedModule.CachePath),
		}

		// Add module configuration
		for key, value := range moduleRef.Config {
			attributes[key] = convertToCty(value)
		}

		// Add default attributes from app config
		attributes["project"] = cty.StringVal(appConfig.Project)
		attributes["venue"] = cty.StringVal(appConfig.Venue)
		attributes["installprefix"] = cty.StringVal(appConfig.InstallPrefix)

		// Append module block
		appendBlockToBody(rootBody, "module", []string{moduleName}, resolvedModule.CachePath, attributes)

		// Handle depends_on if specified
		if len(moduleRef.DependsOn) > 0 {
			// Note: This would need more sophisticated handling for proper depends_on syntax
			log.Infof("Module %s depends on: %v", moduleName, moduleRef.DependsOn)
		}
	}

	return nil
}
