package models

import (
	"gorm.io/gorm"
)

type CoreConfig struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Key   string `gorm:"index;unique" json:"key"`
	Value string `json:"value"`
}

type SSMParameters struct {
	gorm.Model
	Key   string `gorm:"index;unique"`
	Value string
	Type  string
	Owner string
}
