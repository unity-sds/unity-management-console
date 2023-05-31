package database

import (
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-control-plane/backend/internal/database/models"
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
//    updatedConfigs, err := StoreConfig(configs)
//    if err != nil {
//        // Handle the error
//    }
//
// Parameters:
//   config : A slice of CoreConfig models to be stored in the database
//
// Returns:
//   ([]models.CoreConfig, error) : A slice of updated CoreConfig models and error (if any)
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
//    configs, err := FetchConfig()
//    if err != nil {
//        // Handle the error
//    }
//
// Returns:
//   ([]models.CoreConfig, error) : A slice of CoreConfig models and an error (if any)
func FetchConfig() ([]models.CoreConfig, error) {
	var config []models.CoreConfig

	// Fetch all records from the database
	if err := DB.Find(&config).Error; err != nil {
		return nil, err
	}

	return config, nil
}