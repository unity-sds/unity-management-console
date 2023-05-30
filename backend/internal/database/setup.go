package database

import (
	"github.com/unity-sds/unity-control-plane/backend/internal/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type GormDatastore struct {
	db *gorm.DB
}

func (g GormDatastore) FetchCoreParams() ([]models.CoreConfig, error) {
	var config []models.CoreConfig
	result := g.db.Find(&config)

	if result.Error != nil {
		return nil, result.Error
	}

	return config, nil
}

func NewGormDatastore() (*GormDatastore, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &GormDatastore{
		db: db,
	}, nil
}

type Datastore interface {
	FetchCoreParams() ([]models.CoreConfig, error)
}
