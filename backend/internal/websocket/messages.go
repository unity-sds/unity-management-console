package websocket

import "github.com/unity-sds/unity-control-plane/backend/internal/database/models"

type BareMessage struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}
type ConfigMessage struct {
	Action  string              `json:"action"`
	Payload []models.CoreConfig `json:"payload"`
}

type InstallMessage struct {
	Action  string                 `json:"action"`
	Payload []models.InstallConfig `json:"payload"`
}
