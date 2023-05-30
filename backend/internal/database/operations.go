package database

import (
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-control-plane/backend/internal/database/models"
	"gorm.io/gorm/clause"
)

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

func FetchConfig() ([]models.CoreConfig, error) {
	var config []models.CoreConfig

	// Fetch all records from the database
	if err := DB.Find(&config).Error; err != nil {
		return nil, err
	}

	return config, nil
}
