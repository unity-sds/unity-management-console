package processes

import (
	"fmt"
	"github.com/unity-sds/unity-control-plane/backend/internal/marketplace"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"net/http"
)

func fetchMarketplaceMetadata(name string, version string) (marketplace.MarketplaceMetadata, error) {

	resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/unity-sds/unity-marketplace/main/applications/%s/%s/metadata.json", name, version))
	if err != nil {
		fmt.Println("Error:", err)
		return marketplace.MarketplaceMetadata{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return marketplace.MarketplaceMetadata{}, err
	}

	content := string(body)
	req := &marketplace.MarketplaceMetadata{}
	err = protojson.Unmarshal([]byte(content), req)
	return *req, err
}
