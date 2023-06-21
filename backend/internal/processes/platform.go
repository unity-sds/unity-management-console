package processes

import "github.com/unity-sds/unity-control-plane/backend/internal/marketplace"

func fetchInstalledApplications() {

}

func checkExistingInstallation(metadata marketplace.MarketplaceMetadata) (bool, error) {

	// Fetch installed applications
	fetchInstalledApplications()

	// Check that selected application isn't installed

	// If it is installed ensure that the name varies to prevent clashes

	return false, nil
}
func checkDependencies(metadata marketplace.MarketplaceMetadata) (bool, error) {

	// Check metadata for dependencies

	// If dependencies exist ensure that they are fulfilled

	return true, nil
}
