package processes

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"net/http"
	"os"
	"strings"
)

func fetchMarketplaceMetadata(name string, version string) (marketplace.MarketplaceMetadata, error) {

	log.Infof("Fetching marketplace metadata for, %s, %s", name, version)
	resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/unity-sds/unity-marketplace/main/applications/%s/%s/metadata.json", name, version))
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
	log.Infof("Error unmarshalling file: %v", err)
	return *req, err
}

func FetchPackage(meta *marketplace.MarketplaceMetadata) (string, error) {
	// Get package

	locationdir := ""
	if strings.HasSuffix(meta.Package, ".zip") {
		// Fetch from zip

	} else {
		// Checkout git repo
		locationdir, err := gitclone(meta.Package)
		return locationdir, err
	}
	return locationdir, nil
}

func gitclone(url string) (string, error) {
	tempDir, err := os.MkdirTemp("", "git-")
	if err != nil {
		return tempDir, err
	}
	_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	return tempDir, err
}
