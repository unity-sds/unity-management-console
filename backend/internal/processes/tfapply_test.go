package processes

import (
	"testing"
	"os"
	"context"
	tfexec "github.com/hashicorp/terraform-exec/tfexec"
)

var TEST_REPO_URL = "https://github.com/unity-sds/unity-sps-prototype"
var TEST_REPO_VARS_URL = "file:///Users/ryhunter/projects/unity/unity-sps-prototype/terraform-unity/tfvars/MCP-DEV_unity-dev-sps-hysds-eks-ryan.tfvars" 
var TEST_DEPLOYABLE_NAME = "unity-dev-sps-hysds-eks-ryan"

func TestCreateDeployment(t *testing.T) {
	_, deploymentStagingPath, err := CreateDeployment("./deployments")
	defer os.RemoveAll(deploymentStagingPath)

	if err != nil {
		t.Fatalf("CreateDeployment failed with %s", err)		
	}

	deploymentStagingDirInfo, err := os.Stat(deploymentStagingPath)
	if err != nil {
		t.Fatalf("Couldn't stat %s, got error %s", deploymentStagingPath, err)		
	}

	if !deploymentStagingDirInfo.IsDir() {
		t.Fatalf("%s is not a directory", deploymentStagingPath)		
	}

	return
}

func TestGetDeployableVariableFile(t *testing.T) {

	varFile, err := os.Create("test.tfvar")
	if err != nil {
		t.Fatalf("failed to create test file: %s", err)
	}
	defer varFile.Close()
	defer os.Remove("test.tfvar")

	err = getDeployableVariableFile(varFile, TEST_REPO_VARS_URL)
	if err != nil {
		t.Fatalf("failed to get deployable: %s", err)
	}

	return
}

func TestStageDeployable(t *testing.T) {
	_, deploymentStagingPath, err := CreateDeployment("./deployments")

	defer os.RemoveAll(deploymentStagingPath)

	stagePath, err := StageDeployable(TEST_REPO_URL, TEST_REPO_VARS_URL, deploymentStagingPath)

	stageDirInfo, err := os.Stat(stagePath)
	if err != nil {
		t.Fatalf("Couldn't stat %s, got error %s", stagePath, err)		
	}

	if !stageDirInfo.IsDir() {
		t.Fatalf("%s is not a directory", stagePath)		
	}

	return
}

func TestValidateDeployable(t *testing.T) {
	_, deploymentStagingPath, err := CreateDeployment("./deployments")

	defer os.RemoveAll(deploymentStagingPath)

	stagePath, err := StageDeployable(TEST_REPO_URL, TEST_REPO_VARS_URL, deploymentStagingPath)
	tf, err := tfexec.NewTerraform(stagePath, "terraform")

	tf.Init(context.Background(), tfexec.Upgrade(true))

	isValid, tfErrorString, err := validateDeployable(tf)

	if err != nil {
		t.Fatalf("Validation error: %s", err)		
	}

	if ! isValid {
		t.Fatalf("Validation failed for staged deployable: %s", tfErrorString)
	}
	return
}

func TestInstallMarketplaceApplication(t *testing.T) {
	t.SkipNow()
	InstallMarketplaceApplication(TEST_REPO_URL, TEST_REPO_VARS_URL, TEST_DEPLOYABLE_NAME)
}

func TestInstallMarketplaceApplicationRightInterface(t *testing.T) {
	meta := &MarketplaceMetadata{
		Package: TEST_REPO_URL,
	}

	install := &Install{
		VariablesReference: TEST_REPO_VARS_URL,
		DeploymentName: TEST_DEPLOYABLE_NAME,
	}

	appConfig := &AppConfig{}

	ds := Datastore{}

	ws := WebSocketManager{}

	InstallMarketplaceApplicationRightInterface(ws, "ryan", appConfig, meta, "./deployments", install, ds)
}