package processes

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func fetchMarketplaceMetadata(name string, version string, appConfig *config.AppConfig) (marketplace.MarketplaceMetadata, error) {

	log.Infof("Fetching marketplace metadata for, %s, %s", name, version)
	url := fmt.Sprintf("%sunity-sds/unity-marketplace/main/applications/%s/%s/metadata.json", appConfig.MarketplaceBaseUrl, name, version)
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Error fetching from github: %v", err)
		return marketplace.MarketplaceMetadata{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading file: %v", err)
		return marketplace.MarketplaceMetadata{}, err
	}

	content := string(body)
	req := &marketplace.MarketplaceMetadata{}
	err = protojson.Unmarshal([]byte(content), req)
	if err != nil {
		log.Infof("Error unmarshalling file: %v", err)
	}
	return *req, err
}

func FetchPackage(meta *marketplace.MarketplaceMetadata, appConfig *config.AppConfig) (string, error) {
	// Get package
	basedir := "/tmp"
	if meta.Backend == "terraform" {
		basedir = filepath.Join(appConfig.Workdir, "..", "terraform", "modules")
	}
	if strings.HasSuffix(meta.Package, ".zip") {
		// Fetch from zip
		return "", nil
	} else {
		// Checkout git repo
		locationdir, err := gitclone(meta.Package, basedir)
		return locationdir, err
	}

}

func gitclone(url string, basedir string) (string, error) {
	tempDir, err := os.MkdirTemp(basedir, "git-")
	if err != nil {
		return tempDir, err
	}
	_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	return tempDir, err
}
