package terraform

import (
	"fmt"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/hashicorp/hcl/v2/hclwrite"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
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

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func WriteTFVars(vars []Varstruct, appconf *config.AppConfig) {
	// Open a file for writing
	file, err := os.Create(filepath.Join(appconf.Workdir, "workspace", ".auto.tfvars"))
	if err != nil {
		log.WithError(err).Errorf("Error creating file")
		return
	}
	defer file.Close()

	sliceToString := func(slice []string) string {
		for i, v := range slice {
			slice[i] = `"` + v + `"`
		}
		return "[" + strings.Join(slice, ", ") + "]"
	}
	// Iterate through the map and write key-value pairs to the file
	for _, variable := range vars {
		line := ""
		switch v := variable.Value.(type) {
		case map[string][]string:
			line += fmt.Sprintf("%s = { ", variable.Name)
			for k, val := range v {
				line += fmt.Sprintf("%s: %s, ", k, sliceToString(val))
			}
			line = line[:len(line)-2] + " }\n" // Remove trailing comma and add closing bracket
		case []string:
			line = fmt.Sprintf("%s = %s\n", variable.Name, sliceToString(v))
		default:
			line = fmt.Sprintf("%s = \"%v\"\n", variable.Name, variable.Value)
		}
		_, err := file.WriteString(line)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	log.Info("File written successfully")

	hclFile := hclwrite.NewEmptyFile()

	err = os.MkdirAll(filepath.Join(appconf.Workdir, "workspace"), 0755)
	if err != nil {
		log.WithError(err).Error("Could not create workspace directory")
	}
	// create new file on system
	tfFile, err := os.Create(filepath.Join(appconf.Workdir, "workspace", "variables.tf"))
	if err != nil {
		log.WithError(err).Error("Problem creating tf file")
		return
	}
	// initialize the body of the new file object
	rootBody := hclFile.Body()
	for _, variable := range vars {
		cloudenv := rootBody.AppendNewBlock("variable", []string{variable.Name})
		cloudenvBody := cloudenv.Body()
		typeTokens := hclwrite.Tokens{
			{
				Type:  9,
				Bytes: []byte(strings.ToLower(strings.ReplaceAll(variable.Tfobjtype.GoString(), "cty.", ""))),
			},
		}
		cloudenvBody.SetAttributeRaw("type", typeTokens)
	}
	_, err = tfFile.Write(hclFile.Bytes())
	if err != nil {
		log.WithError(err).Error("error writing hcl file")
		return
	}
	return
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

func parseAdvancedVariables(install *marketplace.Install, cloudenv *hclwrite.Block) {
	for key, str := range convertStruct(install.Applications.Variables.AdvancedValues) {
		ctyValue := convertToCty(str)
		cloudenv.Body().SetAttributeValue(key, ctyValue)
	}
}

func AddApplicationToStack(appConfig *config.AppConfig, location string, meta *marketplace.MarketplaceMetadata, install *marketplace.Install) error {
	rand.Seed(time.Now().UnixNano())

	s := String(8)
	hclFile := hclwrite.NewEmptyFile()

	err := os.MkdirAll(filepath.Join(appConfig.Workdir, "workspace"), 0755)
	if err != nil {
		log.WithError(err).Error("Could not create workspace directory")
	}
	// create new file on system
	tfFile, err := os.Create(filepath.Join(appConfig.Workdir, "workspace", fmt.Sprintf("%v%v%v", meta.Name, s, ".tf")))
	if err != nil {
		log.WithError(err).Error("Problem creating tf file")
		return err
	}
	path := filepath.Join(location, meta.WorkDirectory)
	// initialize the body of the new file object
	rootBody := hclFile.Body()

	cloudenv := rootBody.AppendNewBlock("module", []string{meta.Name})
	cloudenvBody := cloudenv.Body()
	cloudenvBody.SetAttributeValue("source", cty.StringVal(path))
	cloudenvBody.SetAttributeValue("name", cty.StringVal(install.DeploymentName))
	typeTokens := hclwrite.Tokens{
		{
			Type:  9,
			Bytes: []byte("{}"),
		},
	}
	cloudenvBody.SetAttributeRaw("tags", typeTokens)
	for key, element := range install.Applications.Variables.Values {
		cloudenv.Body().SetAttributeValue(key, cty.StringVal(element))
	}

	//for key, element := range install.Applications.Variables.NestedValues {
	//	tmap := map[string]cty.Value{}
	//	//t := element.GetType()
	//	con := element.GetConfig()
	//	for k, e := range con {
	//		tmap[k] = cty.StringVal(e.GetOptions().GetDefault())
	//	}
	//	cloudenv.Body().SetAttributeValue(key, cty.MapVal(tmap))
	//}

	parseAdvancedVariables(install, cloudenv)

	_, err = tfFile.Write(hclFile.Bytes())
	if err != nil {
		log.WithError(err).Error("error writing hcl file")
		return err
	}
	return nil
}
