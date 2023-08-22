package database

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/database/models"
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

func (g GormDatastore) UpdateApplicationStatus(deploymentID uint, targetAppName string, newStatus string) error {
	var deployment models.Deployment
	result := g.db.Preload("Applications").First(&deployment, deploymentID) // Replace deploymentID with the actual ID you're looking for
	if result.Error != nil {
		// Handle the error
		log.WithError(result.Error).Error("Error finding deployment")
		return result.Error
	}

	// Then, find the specific application by name within the deployment
	var appToUpdate *models.Application
	for _, app := range deployment.Applications {
		if app.Name == targetAppName {
			appToUpdate = &app
			break
		}
	}

	if appToUpdate == nil {
		err := errors.New("Application not found")
		log.WithError(err).Error("Problem finding application")
		return err
	}

	// Finally, update the status field of the application
	appToUpdate.Status = newStatus
	if err := g.db.Save(appToUpdate).Error; err != nil {
		log.WithError(err).Error("Problem updating application status")
		return err
	}

	return nil
}
