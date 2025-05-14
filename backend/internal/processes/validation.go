package processes

import (
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/types"
)

func CheckIAMPolicies() error {
	// Check IAM policies

	// Get default polcies from marketplace

	// Run IAM Simulator
	return nil
}

func ValidateMarketplaceInstallation(name string, version string, appConfig *config.AppConfig) (bool, types.MarketplaceMetadata, error) {
	// Validate installation

	// Check Marketplace Installation Exists
	meta, err := FetchMarketplaceMetadata(name, version, appConfig)
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

func validateAndPrepareInstallation(app *types.ApplicationInstallParams, appConfig *config.AppConfig) (*types.MarketplaceMetadata, error) {
	log.Info("Validating installation")
	validMarket, meta, err := ValidateMarketplaceInstallation(app.Name, app.Version, appConfig)
	if err != nil || !validMarket {
		log.Error("Invalid marketplace installation:", err)
		return nil, err
	}

	log.Info("Checking IAM Policies")
	if err := CheckIAMPolicies(); err != nil {
		log.Error("Error checking IAM policies:", err)
		return nil, err
	}

	return &meta, nil
}
