package database

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/database/models"
	"gorm.io/gorm/clause"
	"gorm.io/gorm"
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

func (g GormDatastore) GetInstalledApplicationByName(name string) (*models.InstalledMarketplaceApplication, error) {
	var application models.InstalledMarketplaceApplication
	result := g.db.Where("name = ?", name).Where("status != 'UNINSTALLED'").First(&application)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound){
			return nil, nil
		}

		log.WithError(result.Error).Error("Error finding application")
		return nil, result.Error
	}
	return &application, nil
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

func (g GormDatastore) FetchAllInstalledMarketplaceApplications() ([]models.InstalledMarketplaceApplication, error) {
	var applications []models.InstalledMarketplaceApplication
	result := g.db.Find(&applications)
	if result.Error != nil {
		return nil, result.Error
	}
	return applications, nil
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

func (g GormDatastore) StoreInstalledMarketplaceApplication(model models.InstalledMarketplaceApplication) (error) {
	if err := g.db.Save(&model).Error; err != nil {
		// Handle error for Save
		log.WithError(err).Error("Problem saving record to database")
		return err
	}
	return nil
}

func (g GormDatastore) GetInstalledMarketplaceApplicationStatusByName(appName string, deploymentName string) (*models.InstalledMarketplaceApplication, error) {
	var application models.InstalledMarketplaceApplication
	err := g.db.Where("Name = ? AND deployment_name = ?", appName, deploymentName).First(&application).Error
	if err != nil {
		log.WithError(err).Error("Problem getting application status")
		return nil, err
	}
	return &application, nil
}

func (g GormDatastore) UpdateInstalledMarketplaceApplicationStatusByName(appName string, deploymentName string, status string) (error) {
	var app models.InstalledMarketplaceApplication
	
	g.db.Where("name = ? AND deployment_name = ?", appName, deploymentName).First(&app)
	app.Status = status

	if err := g.db.Save(&app).Error; err != nil {
		// Handle error for Save
		log.WithError(err).Error("Problem saving record to database")
		return err
	}
	return nil
}

func (g GormDatastore) RemoveInstalledMarketplaceApplicationByName(appName string) (error) {
	if err := g.db.Where("name != ?", appName).Delete(&models.InstalledMarketplaceApplication{}).Error; err != nil {
		return err
	}
	return nil
}