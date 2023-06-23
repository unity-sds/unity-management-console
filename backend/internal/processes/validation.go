package processes

import (
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-control-plane/backend/internal/marketplace"
)

func CheckIAMPolicies() error {
	// Check IAM policies

	// Get default polcies from marketplace

	// Run IAM Simulator
	return nil
}

func ValidateMarketplaceInstallation(name string, version string) (bool, marketplace.MarketplaceMetadata, error) {
	// Validate installation

	// Check Marketplace Installation Exists
	meta, err := fetchMarketplaceMetadata(name, version)
	if err != nil {
		return false, meta, err
	}

	// Is already installed?
	exists, err := checkExistingInstallation(meta)
	if exists == true || err != nil {
		return false, meta, err
	}

	// Do dependencies match?
	depvalid, err := checkDependencies(meta)
	if depvalid == false || err != nil {
		return false, meta, err
	}

	return true, meta, nil
}

func validateAndPrepareInstallation(app *marketplace.Install_Applications) (*marketplace.MarketplaceMetadata, error) {
	log.Info("Validating installation")
	validMarket, meta, err := ValidateMarketplaceInstallation(app.Name, app.Version)
	if err != nil || !validMarket {
		log.Error("Invalid marketplace installation:", err)
		return &marketplace.MarketplaceMetadata{}, err
	}

	log.Info("Checking IAM Policies")
	if err := CheckIAMPolicies(); err != nil {
		log.Error("Error checking IAM policies:", err)
		return &marketplace.MarketplaceMetadata{}, err
	}

	return &meta, nil
}
