package processes

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-control-plane/backend/internal/marketplace"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"net/http"
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
