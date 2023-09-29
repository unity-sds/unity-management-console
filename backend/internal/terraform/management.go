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
	"github.com/unity-sds/unity-management-console/backend/internal/database"
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

func parseAdvancedVariables(install *marketplace.Install, cloudenv *map[string]cty.Value) {
	if install.Applications.Variables.AdvancedValues != nil && len(install.Applications.Variables.AdvancedValues.Fields) > 0 {
		for key, str := range convertStruct(install.Applications.Variables.AdvancedValues) {
			ctyValue := convertToCty(str)
			(*cloudenv)[key] = ctyValue
		}
	}
}

func generateMetadataHeader(cloudenv *hclwrite.Body, id string, application string, applicationName string, version string, creator string, deploymentID uint) {
	currentTime := time.Now()
	dateString := currentTime.Format("2006-01-02")
	comment := hclwrite.Tokens{
		&hclwrite.Token{
			Type:         hclsyntax.TokenComment,
			Bytes:        []byte(fmt.Sprintf("# id: %v\n# application: %v\n# applicationName: %v\n# version: %v\n# creator: %v\n# creationDate: %v\n# deploymentID: %v\n", id, application, applicationName, version, creator, dateString, deploymentID)),
			SpacesBefore: 0,
		},
	}
	cloudenv.AppendUnstructuredTokens(comment)
}

func generateRandomString(length int) string {
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
func AddApplicationToStack(appConfig *config.AppConfig, location string, meta *marketplace.MarketplaceMetadata, install *marketplace.Install, db database.Datastore, deploymentID uint) error {
	rand.Seed(time.Now().UnixNano())

	s := generateRandomString(8)
	hclFile := hclwrite.NewEmptyFile()

	directory := filepath.Join(appConfig.Workdir, "workspace")
	filename := fmt.Sprintf("%v%v%v", install.Applications.Name, s, ".tf")
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

	generateMetadataHeader(rootBody, u.String(), meta.Name, install.Applications.Displayname, install.Applications.Version, "admin", deploymentID)

	attributes := map[string]cty.Value{
		"deployment_name": cty.StringVal(install.Applications.Displayname),
		"tags":            cty.MapValEmpty(cty.String), // Example of setting an empty map
		// Add other attributes as needed
	}
	for key, element := range install.Applications.Variables.Values {
		if strings.HasPrefix("*", element) {
			log.Infof("Element %s has prefix: %s", key, element)
			element, err = lookUpVariablePointer(element, install)
			if err != nil {
				return err
			}
		} else if strings.HasPrefix("\\*", element) {
			element = strings.Replace(element, "\\", "", 1)
		}
		attributes[key] = cty.StringVal(element)
	}
	parseAdvancedVariables(install, &attributes)
	rand.Seed(time.Now().UnixNano())
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	randomChars := make([]byte, 5)
	for i, v := range rand.Perm(52)[:5] {
		randomChars[i] = chars[v]
	}
	appendBlockToBody(rootBody, "module", []string{fmt.Sprintf("%s-%s", install.Applications.Displayname, string(randomChars))}, path, attributes)

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
		log.Infof("Checking dependency: %s, %s for value %s", k, v, element)
		if k == element {
			return v, nil
		}
	}

	return "", nil

}
