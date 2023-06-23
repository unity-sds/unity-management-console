package database

import (
	"github.com/unity-sds/unity-control-plane/backend/internal/application/config"
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

func (g GormDatastore) FetchSSMParams() ([]models.SSMParameters, error) {

	var ssmparams []models.SSMParameters

	result := g.db.Find(&ssmparams)

	if result.Error != nil {
		return nil, result.Error
	}

	return ssmparams, nil
}

func NewGormDatastore() (Datastore, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.SSMParameters{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.CoreConfig{})
	if err != nil {
		return nil, err
	}

	return &GormDatastore{
		db: db,
	}, nil
}

type Datastore interface {
	FetchCoreParams() ([]models.CoreConfig, error)
	FetchSSMParams() ([]models.SSMParameters, error)
	StoreSSMParams(p []config.SSMParameter, owner string) error
}
