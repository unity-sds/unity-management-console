package processes

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"github.com/unity-sds/unity-management-console/backend/internal/terraform"
	"github.com/zclconf/go-cty/cty"
	"os"
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
func UpdateCoreConfig(appConfig *config.AppConfig, db database.Datastore) error {

	//
	//varstr := ""
	//for _, v := range cParams {
	//	if v.Key == "/unity/core/project" {
	//		env["TF_VAR_project"] = v.Value
	//	} else if v.Key == "/unity/core/venue" {
	//		env["TF_VAR_venue"] = v.Value
	//
	//	} else {
	//		varstr = varstr + fmt.Sprintf("{name  = \"%v\", type  = \"%v\", value = \"%v\"},", v.Key, v.Type, v.Value)
	//	}
	//}
	//
	//if varstr != "" {
	//	varstr = strings.TrimRight(varstr, ",")
	//	varstr = fmt.Sprintf("[%v]", varstr)
	//
	//	env["TF_VAR_ssm_parameters"] = varstr
	//
	//}
	//secrets := map[string]string{}
	//return r.RunAct(config.Workdir+"/environment-provisioner.yml", inputs, env, secrets, conn, config)

	updateAutotfvars(appConfig)
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// create new file on system
	tfFile, err := os.Create(appConfig.Workdir + "/bservelist.tf")
	if err != nil {
		log.WithError(err).Error("Problem creating tf file")
		return err
	}
	// initialize the body of the new file object
	rootBody := hclFile.Body()

	cloudenv := rootBody.AppendNewBlock("module", []string{"example_module"})
	cloudenvBody := cloudenv.Body()
	cloudenvBody.SetAttributeValue("source", cty.StringVal(appConfig.Workdir+"/../terraform/modules/unity-cloud-env"))
	cloudenvBody.SetAttributeValue("venue", cty.StringVal("ami-123"))
	cloudenvBody.SetAttributeValue("project", cty.StringVal("t2.micro"))
	cloudenvBody.SetAttributeValue("publicsubnets", cty.StringVal("123"))
	cloudenvBody.SetAttributeValue("privatesubnets", cty.StringVal("123"))
	ssmParameters, err := generateSSMParameters(db)
	if err != nil {
		log.WithError(err).Error("Problem fetching params")
		return err
	}
	cloudenvBody.SetAttributeValue("ssm_parameters", cty.ListVal(ssmParameters))

	_, err = tfFile.Write(hclFile.Bytes())

	terraform.RunTerraform(appConfig)
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
