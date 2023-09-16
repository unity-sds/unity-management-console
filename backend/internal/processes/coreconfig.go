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

func UpdateCoreConfig(appConfig *config.AppConfig, db database.Datastore, websocketmgr *websocket.WebSocketManager, userid string) error {

	err := os.MkdirAll(filepath.Join(appConfig.Workdir, "workspace"), 0755)
	if err != nil {
		log.WithError(err).Error("Couldn't create workdir for core config")
		return err
	}
	// create new file on system
	tfFile, err := os.Create(filepath.Join(appConfig.Workdir, "workspace", "params.tf"))
	if err != nil {
		log.WithError(err).Error("Problem creating tf file")
		return err
	}
	defer tfFile.Close()

	venue, err := getSSMParameterValueFromDatabase("venue", db)
	if err != nil {
		log.WithError(err).Error("Problem fetching venue")
		return err
	}
	project, err := getSSMParameterValueFromDatabase("project", db)
	if err != nil {
		log.WithError(err).Error("Problem fetching project")
		return err
	}
	publicsubnets, err := getSSMParameterValueFromDatabase("publicsubnets", db)
	if err != nil {
		log.WithError(err).Error("Problem fetching public subnets")
		return err
	}
	privatesubnets, err := getSSMParameterValueFromDatabase("privatesubnets", db)
	if err != nil {
		log.WithError(err).Error("Problem fetching private subnets")
		return err
	}

	ssmParameters, err := generateSSMParameters(db)
	if err != nil {
		log.WithError(err).Error("Problem fetching params")
		return err
	}

	hclFile := generateFileStructure(appConfig, venue, project, publicsubnets, privatesubnets, ssmParameters)

	_, err = tfFile.Write(hclFile.Bytes())
	executor := &terraform.RealTerraformExecutor{}

	err = terraform.RunTerraform(appConfig, websocketmgr, userid, executor, "")
	return err
}

func generateFileStructure(appConfig *config.AppConfig, venue, project, publicsubnets, privatesubnets string, ssmParameters []cty.Value) *hclwrite.File {
	hclFile := hclwrite.NewEmptyFile()
	// initialize the body of the new file object
	rootBody := hclFile.Body()

	cloudenv := rootBody.AppendNewBlock("module", []string{"unity-cloud-env"})
	cloudenvBody := cloudenv.Body()
	cloudenvBody.SetAttributeValue("source", cty.StringVal(filepath.Join(appConfig.Workdir, "..", "terraform", "modules", "unity-cloud-env")))
	cloudenvBody.SetAttributeValue("venue", cty.StringVal(venue))
	cloudenvBody.SetAttributeValue("project", cty.StringVal(project))
	cloudenvBody.SetAttributeValue("publicsubnets", cty.StringVal(publicsubnets))
	cloudenvBody.SetAttributeValue("privatesubnets", cty.StringVal(privatesubnets))

	if ssmParameters != nil {
		cloudenvBody.SetAttributeValue("ssm_parameters", cty.ListVal(ssmParameters))
	}
	return hclFile
}

func generateSSMParameters(db database.Datastore) ([]cty.Value, error) {
	params, err := db.FetchSSMParams()
	var ssmParameters []cty.Value

	if err != nil {
		log.WithError(err).Error("Error fetching ssm parameters")
		return ssmParameters, err
	}

	for _, p := range params {
		if p.Key != "" {
			a := cty.ObjectVal(map[string]cty.Value{
				"name":  cty.StringVal(p.Key),
				"type":  cty.StringVal(p.Type),
				"value": cty.StringVal(p.Value),
			})
			ssmParameters = append(ssmParameters, a)
		}
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
