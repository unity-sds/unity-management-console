package models

import (
	"gorm.io/gorm"
	"time"
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
