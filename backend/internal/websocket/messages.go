package websocket

import (
	"encoding/json"
	"github.com/unity-sds/unity-control-plane/backend/internal/database/models"
)

type BareMessage struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}
type ConfigMessage struct {
	Action  string              `json:"action"`
	Payload []models.CoreConfig `json:"payload"`
}
