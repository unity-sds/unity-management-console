package processes

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"github.com/unity-sds/unity-management-console/backend/internal/terraform"
	"github.com/unity-sds/unity-management-console/backend/internal/websocket"
	"github.com/zclconf/go-cty/cty"
	"os"
	"path/filepath"
)

func updateAutotfvars(appConf *config.AppConfig) {
	file, err := os.Create(appConf.Workdir + "/.auto.tfvars")
	if err != nil {
		log.WithError(err).Error("Failed to create tfvars file")
	}
	defer file.Close()

	varconfig := `project = "unity-nightly"

privatesubnets = "subnet-059bc4f467275b59d,subnet-0ebdd997cc3ebe58d"
publicsubnets = "subnet-087b54673c7549e2d,subnet-009c32904a8bf3b92"

venue = "dev"
ssm_parameters = [{name  = "parameter1dev", type  = "String", value = "value1"}]`
	_, err = file.WriteString(varconfig)
	if err != nil {
		log.WithError(err).Error("Failed to write string to tfvars file")
	}

	// Save changes to the file
	err = file.Sync()
	if err != nil {
		log.WithError(err).Error("Failed to sync tfvars file")
	}

	log.Println("File '.auto.tf' has been written")
}
func UpdateCoreConfig(appConfig *config.AppConfig, db database.Datastore, websocketmgr *websocket.WebSocketManager, userid string) error {

	//updateAutotfvars(appConfig)
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// create new file on system
	tfFile, err := os.Create(appConfig.Workdir + "/params.tf")
	if err != nil {
		log.WithError(err).Error("Problem creating tf file")
		return err
	}
	// initialize the body of the new file object
	rootBody := hclFile.Body()

	venue, err := getSSMParameterValueFromDatabase("venue", db)
	project, err := getSSMParameterValueFromDatabase("project", db)
	publicsubnets, err := getSSMParameterValueFromDatabase("publicsubnets", db)
	privatesubnets, err := getSSMParameterValueFromDatabase("privatesubnets", db)

	cloudenv := rootBody.AppendNewBlock("module", []string{"unity-cloud-env"})
	cloudenvBody := cloudenv.Body()
	cloudenvBody.SetAttributeValue("source", cty.StringVal(filepath.Join(appConfig.Workdir, "..", "terraform", "modules", "unity-cloud-env")))
	cloudenvBody.SetAttributeValue("venue", cty.StringVal(venue))
	cloudenvBody.SetAttributeValue("project", cty.StringVal(project))
	cloudenvBody.SetAttributeValue("publicsubnets", cty.StringVal(publicsubnets))
	cloudenvBody.SetAttributeValue("privatesubnets", cty.StringVal(privatesubnets))

	ssmParameters, err := generateSSMParameters(db)
	if err != nil {
		log.WithError(err).Error("Problem fetching params")
		return err
	}
	if ssmParameters != nil {
		cloudenvBody.SetAttributeValue("ssm_parameters", cty.ListVal(ssmParameters))
	}

	_, err = tfFile.Write(hclFile.Bytes())

	terraform.RunTerraform(appConfig, websocketmgr, userid)
	return err
}

func generateSSMParameters(db database.Datastore) ([]cty.Value, error) {
	params, err := db.FetchSSMParams()
	var ssmParameters []cty.Value

	if err != nil {
		log.WithError(err).Error("Error fetching ssm parameters")
		return ssmParameters, err
	}

	for _, p := range params {
		a := cty.ObjectVal(map[string]cty.Value{
			"name":  cty.StringVal(p.Key),
			"type":  cty.StringVal(p.Type),
			"value": cty.StringVal(p.Value),
		})
		ssmParameters = append(ssmParameters, a)
	}

	return ssmParameters, err
}

func getSSMParameterValueFromDatabase(paramName string, db database.Datastore) (string, error) {

	params, err := db.FetchSSMParams()
	if err != nil {
		log.WithError(err).Error("Error fetching ssm parameters")
		return "", err
	}
	for _, p := range params {
		if p.Key == paramName {
			return p.Value, nil
		}
	}

	return "", nil

}
