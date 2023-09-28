package database

import (
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-management-console/backend/internal/application"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/database/models"
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

func (g GormDatastore) AddToAudit(operation application.AuditLine, owner string) error {
	a := models.Audit{
		Model:     gorm.Model{},
		Operation: operation.String(),
		Owner:     owner,
	}
	result := g.db.Create(&a)

	if result.Error != nil {
		log.WithError(result.Error).Error("could not insert new Audit")
		return result.Error
	}

	return nil
}

func (g GormDatastore) FindLastAuditLineByOperation(operation application.AuditLine) (models.Audit, error) {
	var audit models.Audit

	// Query the latest entry where Operation equals "config updated"
	result := g.db.Where("operation = ?", operation.String()).Order("created_at desc").First(&audit)

	if result.Error != nil {
		// Handle error
		log.WithError(result.Error).Error("Error looking up audit line")
		return audit, result.Error
	}

	return audit, nil
}

func NewGormDatastore() (Datastore, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Audit{})
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

	err = db.AutoMigrate(&models.Deployment{}, &models.Application{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Application{}, &models.Deployment{})
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
	StoreDeployment(p models.Deployment) (uint, error)
	UpdateApplicationStatus(deploymentid uint, application string, displayName string, status string) error
	FetchDeploymentIDByName(deploymentname string) (uint, error)
	FetchAllApplicationStatus() ([]models.Deployment, error)
	FetchAllApplicationStatusByDeployment(deploymentid uint) ([]models.Application, error)
	AddToAudit(operation application.AuditLine, owner string) error
	FindLastAuditLineByOperation(operation application.AuditLine) (models.Audit, error)

	FetchDeploymentNames() ([]string, error)
	RemoveDeploymentByName(name string) error
	RemoveApplicationByName(deploymentName string, applicationName string) error
}
