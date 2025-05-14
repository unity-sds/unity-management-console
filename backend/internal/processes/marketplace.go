package processes

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/types"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func FetchMarketplaceMetadata(name string, version string, appConfig *config.AppConfig) (types.MarketplaceMetadata, error) {

	log.Infof("Fetching marketplace metadata for, %s, %s", name, version)
	url := fmt.Sprintf("%sunity-sds/unity-marketplace/main/applications/%s/%s/metadata.json", appConfig.MarketplaceBaseUrl, name, version)

	log.Infof("Fetching marketplace metadata at: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Error fetching from github: %v", err)
		errMsg := fmt.Sprintf("Error fetching metadata from url: %s", url)
		return types.MarketplaceMetadata{}, errors.New(errMsg)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading file: %v", err)
		errMsg := fmt.Sprintf("Error fetching metadata from url: %s", url)
		return types.MarketplaceMetadata{}, errors.New(errMsg)
	}

	content := string(body)
	var req types.MarketplaceMetadata
	err = json.Unmarshal([]byte(content), &req)
	if err != nil {
		errMsg := fmt.Sprintf("Error fetching metadata from url: %s", url)
		return types.MarketplaceMetadata{}, errors.New(errMsg)
	}
	return req, err
}

func FetchPackage(meta *types.MarketplaceMetadata, appConfig *config.AppConfig) (string, error) {
	// Get package
	basedir := "/tmp"
	if meta.TerraformModuleName != "" {
		basedir = filepath.Join(appConfig.Workdir, "terraform", "modules", meta.Name, meta.Version)
	}
	if strings.HasSuffix(meta.Package, ".zip") {
		// Fetch from zip
		return "", nil
	} else {
		// Checkout git repo
		locationdir, err := gitClone(meta.Package, basedir)
		return locationdir, err
	}

}

func gitClone(url string, basedir string) (string, error) {
	log.Infof("Cloning: %s into %s", url, basedir)
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

	log.Infof("Basedir: %s, url: %s, sha: %s", basedir, url, sha)

	repo, err := git.PlainClone(basedir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		if err.Error() != "repository already exists" {
			return "", err
		}

		// If the repository already exists, open it and pull the latest changes
		repo, err = git.PlainOpen(basedir)
		if err != nil {
			return "", err
		}
		
		// Pull the latest changes
		log.Infof("Repository already exists, pulling latest changes in %s", basedir)
		w, err := repo.Worktree()
		if err != nil {
			log.Errorf("Couldn't get worktree: %v", err)
			return "", err
		}
		
		err = w.Pull(&git.PullOptions{
			Progress: os.Stdout,
		})
		
		if err != nil && err != git.NoErrAlreadyUpToDate {
			log.Errorf("Couldn't pull updates: %v", err)
			return "", err
		} else if err == git.NoErrAlreadyUpToDate {
			log.Infof("Repository already up to date")
		}
	}

	// Checkout the specific SHA if provided
	if sha != "" {
		w, err := repo.Worktree()
		if err != nil {
			log.Infof("Couldn't open repo: %v", err)
			return "", err
		}

		err = w.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(sha),
		})
		if err != nil {
			log.Infof("Couldn't open sha: %v", err)
			return "", err
		}
	}

	return basedir, err
}
