package processes

import "github.com/unity-sds/unity-control-plane/backend/internal/marketplace"

func fetchInstalledApplications() {

}

func checkExistingInstallation(metadata marketplace.MarketplaceMetadata) (bool, error) {

	return true, nil
}
func checkDependencies(metadata marketplace.MarketplaceMetadata) (bool, error) {

	return false, nil
}
