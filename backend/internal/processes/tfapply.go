package processes

import (
	"os"
	"io"
	tfexec "github.com/hashicorp/terraform-exec/tfexec"
	"context"
	"path/filepath"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"strings"
	"fmt"
	"net/http"
	"net/url"
	tfjson "github.com/hashicorp/terraform-json"
)

var DEPLOYABLE_TERRAFORM_DIRECTORY = "terraform-unity"
var DEPLOYMENT_VARIABLES_FILE_NAME = "provided.auto.tfvars"

func cloneDeployableRepo(url string, basedir string) (string, error) {
	sha := ""
	err := os.MkdirAll(basedir, 0755)
	if err != nil {
		return "", err
	}

	// Splitting the URL and SHA if they are in the combined format
	if strings.Contains(url, "@") {
		parts := strings.SplitN(url, "@", 2)
		url = parts[0]
		sha = parts[1]
	}

	repo, err := git.PlainClone(basedir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		if err.Error() != "repository already exists" {
			return "", err
		}

		// If the repository already exists, open it
		repo, err = git.PlainOpen(basedir)
		if err != nil {
			return "", err
		}
	}

	// Checkout the specific SHA if provided
	if sha != "" {
		w, err := repo.Worktree()
		if err != nil {
			return "", err
		}

		err = w.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(sha),
		})
		if err != nil {
			return "", err
		}
	}

	return basedir, err
}

func getDeployableVariableFile(dst io.Writer, variablesUrl string) (error){
	/*
	Fetches file at variablesUrl with http.Get() or os.Open() depending on url scheme and use io.Copy to copy file from memory to dst
	*/
	parsedUrl, err := url.Parse(variablesUrl)
	if err != nil {
		log.Printf("Failed to parse url %s: %s", variablesUrl, err)
		return err
	}

	var variablesReadCloser io.ReadCloser
	switch parsedUrl.Scheme {
	case "http", "https":
		resp, err := http.Get(variablesUrl)
		if err != nil {
			log.Printf("Unable to fetch variables file from %s: %s", variablesUrl, err)
			return err
		}

		variablesReadCloser = resp.Body
	case "file":
		filePath := strings.TrimPrefix(variablesUrl, "file://")
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Unable to open variables files at %s: %s", filePath, err)
			return err
		}
		variablesReadCloser = io.ReadCloser(file)

	}
	_, err = io.Copy(dst, variablesReadCloser)
	if err != nil {
		log.Printf("Failed to copy GET %s response body to destination", variablesUrl)
		return err
	}

	return err
}

func formatTerraformValidationErrors(validationJson *tfjson.ValidateOutput) (string) {
	var errorString = ""
	for _, diagnostic := range validationJson.Diagnostics {
		errorString += diagnostic.Summary
		errorString += "\n"
	}
	return errorString
}

func validateDeployable(tf *tfexec.Terraform) (isValid bool, errorString string, err error) {
	ctx := context.Background()

	validationJson, err := tf.Validate(ctx)

	if err != nil {
		log.Printf("Error while validating terraform")
		return false, "", err
	}

	if !(validationJson.Valid) {
		errorString = formatTerraformValidationErrors(validationJson)
		log.Printf("Terraform validation failed with the following errors:\n%s", errorString)
		return false, errorString, nil
	}

	return true, "", nil
}

func StageDeployable(repoUrl string, variablesUrl string, deploymentStagingPath string) (string, error) {

	// deploymentStaginPath/deployable/DEPLOYABLE_TERRAFORM_DIRECTORY
	repoPath := filepath.Join(deploymentStagingPath, "deployable")
	terraformPath := filepath.Join(repoPath, DEPLOYABLE_TERRAFORM_DIRECTORY)

	cloneDeployableRepo(repoUrl, repoPath)

	terraformDirInfo, err := os.Stat(terraformPath)
	if err != nil || !terraformDirInfo.IsDir() {
		log.Printf("Deployable repo at %s does not contain required directory %s", repoUrl, DEPLOYABLE_TERRAFORM_DIRECTORY)		
		return terraformPath, err
	}

	if variablesUrl != "" {
		varFileName := filepath.Join(terraformPath, DEPLOYMENT_VARIABLES_FILE_NAME)
		varFile, err := os.Create(varFileName)
		if err != nil {
			log.Printf("Failed to created local file %s for variables: %s", varFileName, err)
			return terraformPath, err
		}
		defer varFile.Close()
		getDeployableVariableFile(varFile, variablesUrl)
	}

	return terraformPath, err
}

func ValidateAndDeployDeployable(stagePath string, deployableName string) (bool, error) {
	tf, err := tfexec.NewTerraform(stagePath, "terraform")

	if err != nil {
		log.Printf("Error creating terraform instance: %s", err)
		return false, err
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.BackendConfig(fmt.Sprintf("key=%s", deployableName)))

	if err != nil {
		log.Printf("Failed to init terraform environment: %s", err)
		return false, err
	}

	isValid, tfErrorString, err := validateDeployable(tf)
	if err != nil {
		log.Printf("Error validating deployable: %s", err)
		return false, err
	}

	if !(isValid) {
		log.Printf("Deployable failed validation, with the following errors:")
		log.Printf(tfErrorString)
		return isValid, nil
	}
 
	err = tf.Apply(context.Background())
	if err != nil {
		log.Printf("Faile to apply terraform: %s", err)
		return isValid, err
	}

	return isValid, err

}

func addDeploymentToDatabase(deploymentId uuid.UUID) error {
	return nil
}

func CreateDeployment(deploymentsParentDir string) (uuid.UUID, string, error) {
	/*
	Configures all necesarry elements of a deployment.

	Creates a Deployment ID, an entry in the Deployment database, and a staging directory for the deployment.
	*/
	deploymentId := uuid.New()
	deploymentStagingPath := filepath.Join(deploymentsParentDir, deploymentId.String())
	err := os.MkdirAll(deploymentStagingPath, 0755)

	return deploymentId, deploymentStagingPath, err
}

func InstallMarketplaceApplication(repoUrl string, variablesUrl string, deployableName string) (string, error){
	deploymentId, deploymentStagingPath, err := CreateDeployment("./deployments")
	if err != nil {
		log.Printf("Failed to create deployment")
		return "", err
	}

	stagePath, err := StageDeployable(repoUrl, variablesUrl, deploymentStagingPath)
	if err != nil {
		log.Printf("Failed to stage deployable from %s with variables from %s: %s", repoUrl, variablesUrl, err)
		return deploymentStagingPath, err
	}

	err = addDeploymentToDatabase(deploymentId)
	if err != nil {
		log.Printf("Failed to add %s to database", deploymentId.String())
		return deploymentStagingPath, err
	}

	isValid, err := ValidateAndDeployDeployable(stagePath, deployableName)
	if !(isValid) {
		log.Printf("Invalid deployable staged at %s, from %s", stagePath, repoUrl)
	}
	if err != nil {
		return deploymentStagingPath, err
	}

	return deploymentStagingPath, err	
}

func InstallMarketplaceApplicationRightInterface(conn WebSocketManager,
								userid string, appConfig *AppConfig,
								meta *MarketplaceMetadata,
								location string,
								install *Install,
								db Datastore) error {

	deploymentId, deploymentStagingPath, err := CreateDeployment(location)

	if err != nil {
		log.Printf("Failed to create deployment")
		return err
	}

	stagePath, err := StageDeployable(meta.Package, install.VariablesReference, deploymentStagingPath)
	if err != nil {
		log.Printf("Failed to stage deployable from %s with variables from %s: %s", meta.Package, install.VariablesReference, err)
		return err
	}

	err = addDeploymentToDatabase(deploymentId)
	if err != nil {
		log.Printf("Failed to add %s to database", deploymentId.String())
		return err
	}

	isValid, err := ValidateAndDeployDeployable(stagePath, install.DeploymentName)

	if !(isValid) {
		log.Printf("Invalid deployable staged at %s, from %s", stagePath, meta.Package)
	}
	if err != nil {
		return err
	}

	return err
}