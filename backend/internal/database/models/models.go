package models

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"database/sql/driver"
	"encoding/json"
)

type CoreConfig struct {
	gorm.Model
	ID    uint   `gorm:"primarykey" json:"id"`
	Key   string `gorm:"index;unique" json:"key"`
	Value string `json:"value"`
	Owner string `json:"owner"`
}

type SSMParameters struct {
	gorm.Model
	Key   string `gorm:"index;unique"`
	Value string
	Type  string
	Owner string
}

type Audit struct {
	gorm.Model
	Operation string `gorm:"index"`
	Owner     string
}
type Application struct {
	gorm.Model
	Name         string
	DisplayName  string
	Version      string
	Source       string
	Status       string
	DeploymentID uint
	PackageName  string
	Deployment   Deployment `gorm:"foreignKey:DeploymentID"`
}
type Deployment struct {
	gorm.Model
	Name         string
	Applications []Application `gorm:"foreignKey:DeploymentID;references:ID"`
	Creator      string
	CreationDate time.Time
}

type InstalledMarketplaceApplicationDB struct {
	gorm.Model
	Name         string
	DeploymentName  string
	Version      string
	Source       string
	Status       string
	PackageName  string	
	TerraformModuleName string
	Variables    JSON `json:"variables"`
	AdvancedValues    JSON `json:"advanced_values"`
}

type JSON json.RawMessage

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return errors.New("Invalid Scan Source")
	}
	*j = append((*j)[0:0], s...)
	return nil
}
