package processes

import (
	"github.com/unity-sds/unity-management-console/backend/types"
)

func fetchInstalledApplications() {

}

func checkExistingInstallation(metadata types.MarketplaceMetadata) (bool, error) {

	// Fetch installed applications
	fetchInstalledApplications()

	// Check that selected application isn't installed

	// If it is installed ensure that the name varies to prevent clashes

	return false, nil
}
func checkDependencies(metadata types.MarketplaceMetadata) (bool, error) {

	// Check metadata for dependencies

	// If dependencies exist ensure that they are fulfilled

	return true, nil
}
