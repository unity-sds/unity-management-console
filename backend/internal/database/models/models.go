package models

type CoreConfig struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Key   string `gorm:"index;unique" json:"key"`
	Value string `json:"value"`
}

type InstallConfig struct {
	Install    []AppInstall `json:"install,omitempty"`
	Extensions Extensions   `json:"extenstions,omitempty"`
}

type AppInstall struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Variables map[string]string `json:"variables"`
}

type Extensions map[string]bool
