package database

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/database/models"
	"github.com/unity-sds/unity-management-console/backend/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// StoreConfig stores the given configuration in the database. It uses a
// transaction to ensure all operations are atomic. If any of the operations
// fail, it will rollback the transaction and return an error. In case of a
// conflict on the key field, it updates the existing record's value field
// with the new value.
//
// The function accepts a slice of CoreConfig models and returns the updated
// configurations and an error. The error is non-nil if any issues occur during
// the database operation.
//
// Example usage:
//
//	updatedConfigs, err := StoreConfig(configs)
//	if err != nil {
//	    // Handle the error
//	}
//
// Parameters:
//
//	config : A slice of CoreConfig models to be stored in the database
//
// Returns:
//
//	([]models.CoreConfig, error) : A slice of updated CoreConfig models and error (if any)
func StoreConfig(config []models.CoreConfig) ([]models.CoreConfig, error) {
	// Begin a new transaction
	tx := DB.Begin()

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&config).Error; err != nil {
		log.Errorf("Error inserting into database. %v", err)
		tx.Rollback() // rollback the transaction in case of error
		return nil, err
	}

	// Commit the transaction
	tx.Commit()

	return config, nil
}

// FetchConfig retrieves all configuration records from the database.
// It returns a slice of CoreConfig models and an error. The error
// is non-nil if any issues occur during the database operation.
//
// Example usage:
//
//	configs, err := FetchConfig()
//	if err != nil {
//	    // Handle the error
//	}
//
// Returns:
//
//	([]models.CoreConfig, error) : A slice of CoreConfig models and an error (if any)
func FetchConfig() ([]models.CoreConfig, error) {
	var config []models.CoreConfig

	// Fetch all records from the database
	if err := DB.Find(&config).Error; err != nil {
		return nil, err
	}

	return config, nil
}

func (g GormDatastore) StoreSSMParams(p []config.SSMParameter, owner string) error {
	for _, param := range p {
		model := models.SSMParameters{}

		// Use FirstOrInit to get existing record or initialize a new one
		if err := g.db.Where(models.SSMParameters{Key: param.Name}).FirstOrInit(&model).Error; err != nil {
			// Handle error for FirstOrInit
			log.WithError(err).Error("Problem finding or initializing record in database")
			return err
		}

		// Update the record's fields
		model.Type = param.Type
		model.Value = param.Value
		model.Owner = owner

		log.WithFields(log.Fields{
			"name":  param.Name,
			"type":  param.Type,
			"owner": owner,
		}).Info("Storing SSM parameter")

		// Use Save to insert or update the record
		if err := g.db.Save(&model).Error; err != nil {
			// Handle error for Save
			log.WithError(err).Error("Problem saving record to database")
			return err
		}
	}
	return nil
}

func FetchSSMParams() ([]models.SSMParameters, error) {
	var params []models.SSMParameters

	if err := DB.Find(&params).Error; err != nil {
		return nil, err
	}

	return params, nil
}

func (g GormDatastore) StoreDeployment(model models.Deployment) (uint, error) {
	if err := g.db.Save(&model).Error; err != nil {
		// Handle error for Save
		log.WithError(err).Error("Problem saving record to database")
		return 0, err
	}
	return model.ID, nil
}

func (g GormDatastore) FetchDeploymentIDByName(deploymentID string) (uint, error) {
	var deployment models.Deployment

	result := g.db.Preload("Applications").First(&deployment, deploymentID)
	if result.Error != nil {
		log.WithError(result.Error).Error("Error finding deployment")
		return 0, result.Error
	}
	return deployment.ID, nil
}


func (g GormDatastore) FetchDeploymentIDByApplicationName(deploymentName string) (uint, error) {
	var application models.Application
	result := g.db.Where("display_name = ?", deploymentName).First(&application)
	if result.Error != nil {
		return 0, fmt.Errorf("error finding application: %v", result.Error)
	}
	return application.DeploymentID, nil
}

func (g GormDatastore) UpdateApplicationStatus(deploymentID uint, targetAppName string, displayName string, newStatus string) error {
	var deployment models.Deployment
	result := g.db.Preload("Applications").First(&deployment, deploymentID)
	if result.Error != nil {
		log.WithError(result.Error).Error("Error finding deployment")
		return result.Error
	}

	// Directly find and update the application by name within the deployment
	for index, app := range deployment.Applications {
		if app.Name == targetAppName && app.DisplayName == displayName {
			deployment.Applications[index].Status = newStatus
			if err := g.db.Save(&deployment.Applications[index]).Error; err != nil {
				log.WithError(err).Error("Problem updating application status")
				return err
			}
			return nil
		}
	}

	err := errors.New("application not found")
	log.WithError(err).Error("Problem finding application")
	return err
}

func (g GormDatastore) FetchAllApplicationStatus() ([]models.Deployment, error) {
	var deployments []models.Deployment
	result := g.db.Preload("Applications").Find(&deployments)
	if result.Error != nil {
		return nil, result.Error
	}
	return deployments, nil
}

func parseMarketplaceApplication(dbApp *models.InstalledMarketplaceApplicationDB) (*types.InstalledMarketplaceApplication, error) {
	typeApp := &types.InstalledMarketplaceApplication{
		Name:                dbApp.Name,
		DeploymentName:     dbApp.DeploymentName,
		Version:            dbApp.Version,
		Source:             dbApp.Source,
		Status:             dbApp.Status,
		PackageName:        dbApp.PackageName,
		TerraformModuleName: dbApp.TerraformModuleName,
		Variables:          make(map[string]string),
		AdvancedValues:     make(types.AdvancedValue),
	}

	// Convert Variables JSON to map
	if dbApp.Variables != nil {
		if err := json.Unmarshal([]byte(dbApp.Variables), &typeApp.Variables); err != nil {
			log.WithError(err).Error("Failed to unmarshal Variables")
			return nil, err
		}
	}

	// Convert AdvancedValues JSON to map
	if dbApp.AdvancedValues != nil {
		if err := json.Unmarshal([]byte(dbApp.AdvancedValues), &typeApp.AdvancedValues); err != nil {
			log.WithError(err).Error("Failed to unmarshal AdvancedValues")
			return nil, err
		}
	}

	return typeApp, nil
}

func (g GormDatastore) FetchAllInstalledMarketplaceApplications() ([]*types.InstalledMarketplaceApplication, error) {
	var dbApps []models.InstalledMarketplaceApplicationDB
	result := g.db.Find(&dbApps)
	if result.Error != nil {
		return nil, result.Error
	}

	typeApps := make([]*types.InstalledMarketplaceApplication, 0, len(dbApps))
	for _, dbApp := range dbApps {
		typeApp, err := parseMarketplaceApplication(&dbApp)
		if err != nil {
			return nil, err
		}
		typeApps = append(typeApps, typeApp)
	}

	return typeApps, nil
}

func (g GormDatastore) FetchAllApplicationStatusByDeployment(deploymentid uint) ([]models.Application, error) {
	var deployments models.Deployment
	result := g.db.Preload("Applications").First(&deployments, deploymentid)
	if result.Error != nil {
		return nil, result.Error
	}
	return deployments.Applications, nil
}

func (g GormDatastore) FetchDeploymentNames() ([]string, error) {
	var deployments []models.Deployment

	// Fetch all deployments
	if err := g.db.Find(&deployments).Error; err != nil {
		return nil, err
	}

	// Extract names into a slice of strings
	var names []string
	for _, deployment := range deployments {
		names = append(names, deployment.Name)
	}

	return names, nil
}

func (g GormDatastore) RemoveDeploymentByName(name string) error {
	if err := g.db.Where("name != ?", name).Delete(&models.Deployment{}).Error; err != nil {
		return err
	}
	return nil
}

func (g GormDatastore) RemoveApplicationByName(deploymentName string, applicationName string) error {
	var deployment models.Deployment

	// Retrieve the deployment by name
	err := g.db.Where("name = ?", deploymentName).First(&deployment).Error
	if err != nil {
		return fmt.Errorf("error retrieving deployment: %v", err)
	}

	var application models.Application

	// Retrieve the application by name and DeploymentID
	err = g.db.Where("name = ? AND deployment_id = ?", applicationName, deployment.ID).First(&application).Error
	if err != nil {
		return fmt.Errorf("error retrieving application: %v", err)
	}

	// Delete the application
	err = g.db.Delete(&application).Error
	if err != nil {
		return fmt.Errorf("error deleting application: %v", err)
	}

	return nil
}

func (g GormDatastore) StoreInstalledMarketplaceApplication(application *types.InstalledMarketplaceApplication) error {
	// Convert types model to database model
	dbApp := &models.InstalledMarketplaceApplicationDB{
		Name:                application.Name,
		DeploymentName:     application.DeploymentName,
		Version:            application.Version,
		Source:             application.Source,
		Status:             application.Status,
		PackageName:        application.PackageName,
		TerraformModuleName: application.TerraformModuleName,
	}

	// Convert Variables map to JSON
	if application.Variables != nil {
		varsJSON, err := json.Marshal(application.Variables)
		if err != nil {
			log.WithError(err).Error("Failed to marshal Variables")
			return err
		}
		dbApp.Variables = models.JSON(varsJSON)
	}

	// Convert AdvancedValues map to JSON
	if application.AdvancedValues != nil {
		advJSON, err := json.Marshal(application.AdvancedValues)
		if err != nil {
			log.WithError(err).Error("Failed to marshal AdvancedValues")
			return err
		}
		dbApp.AdvancedValues = models.JSON(advJSON)
	}

	if err := g.db.Save(dbApp).Error; err != nil {
		log.WithError(err).Error("Problem saving record to database")
		return err
	}
	return nil
}

func (g GormDatastore) GetInstalledMarketplaceApplication(appName string, deploymentName string) (*types.InstalledMarketplaceApplication, error) {
	var dbApp models.InstalledMarketplaceApplicationDB
	err := g.db.Where("Name = ? AND deployment_name = ?", appName, deploymentName).First(&dbApp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		log.WithError(err).Error("Problem getting application status")
		return nil, err
	}

	return parseMarketplaceApplication(&dbApp)
}


func (g GormDatastore) UpdateInstalledMarketplaceApplication(application *types.InstalledMarketplaceApplication) error {
	// Find existing record
	var existingApp models.InstalledMarketplaceApplicationDB
	if err := g.db.Where("name = ? AND deployment_name = ?", application.Name, application.DeploymentName).First(&existingApp).Error; err != nil {
		log.WithError(err).Error("Problem finding existing application")
		return err
	}

	// Update fields
	existingApp.Version = application.Version
	existingApp.Source = application.Source
	existingApp.Status = application.Status
	existingApp.PackageName = application.PackageName
	existingApp.TerraformModuleName = application.TerraformModuleName

	// Convert Variables map to JSON
	if application.Variables != nil {
		varsJSON, err := json.Marshal(application.Variables)
		if err != nil {
			log.WithError(err).Error("Failed to marshal Variables")
			return err
		}
		existingApp.Variables = models.JSON(varsJSON)
	}

	// Convert AdvancedValues map to JSON
	if application.AdvancedValues != nil {
		advJSON, err := json.Marshal(application.AdvancedValues)
		if err != nil {
			log.WithError(err).Error("Failed to marshal AdvancedValues")
			return err
		}
		existingApp.AdvancedValues = models.JSON(advJSON)
	}

	if err := g.db.Save(&existingApp).Error; err != nil {
		log.WithError(err).Error("Problem saving record to database")
		return err
	}
	return nil
}


func (g GormDatastore) RemoveInstalledMarketplaceApplication(appName string, deploymentName string) error {
	if err := g.db.Where("name = ? AND deployment_name = ?", appName, deploymentName).Delete(&models.InstalledMarketplaceApplicationDB{}).Error; err != nil {
		return err
	}
	return nil
}
