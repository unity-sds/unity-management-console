package models

type CoreConfig struct {
	ID     string   `json:"id" gorm:"primary_key"`
	Value  string `json:"value"`
}

