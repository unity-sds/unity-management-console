package models

type CoreConfig struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Key   string `gorm:"index;unique" json:"key"`
	Value string `json:"value"`
}
