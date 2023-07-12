package action

import (
	"github.com/unity-sds/unity-control-plane/backend/internal/act"
	"github.com/unity-sds/unity-control-plane/backend/internal/application/config"
	"github.com/unity-sds/unity-control-plane/backend/internal/websocket"
	"github.com/unity-sds/unity-cs-manager/marketplace"
)

type ActRunner interface {
	RunAct(path string, inputs, env, secrets map[string]string, conn *websocket.WebSocketManager, appConfig config.AppConfig) error
}

type ActRunnerImpl struct {
}

// NewActRunner creates a new ActRunnerImpl instance.
func NewActRunner() *ActRunnerImpl {
	return &ActRunnerImpl{}
}
func (r *ActRunnerImpl) RunAct(path string, inputs, env, secrets map[string]string, conn *websocket.WebSocketManager, appConfig config.AppConfig) error {
	return act.RunAct(path, inputs, env, secrets, conn, appConfig)
}
func RunInstall(wsmanager *websocket.WebSocketManager, userid string, install *marketplace.Install, appConfig config.AppConfig) error {

	/*if install.Extensions != nil {
		err := spinUpExtensions(conn, appConfig, install.Extensions, r)
		if err != nil {
			return err
		}
	}*/

	if install.Applications != nil {
		spinUpProjects(install.Applications)
	}
	return nil
}

func spinUpProjects(applications *marketplace.Install_Applications) {

}
