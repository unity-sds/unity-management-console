package processes

import (
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"github.com/unity-sds/unity-management-console/backend/internal/websocket"
)

func fetchAllApplications(store database.Datastore) error {

	dep, err := store.FetchAllApplicationStatus()
	if err != nil {
		return err
	}

	var deployments []*marketplace.Deployment
	for _, d := range dep {
		var apps []*marketplace.Application
		for _, a := range d.Applications {
			app := marketplace.Application{
				ApplicationName: a.Name,
				DisplayName:     a.DisplayName,
				PackageName:     a.PackageName,
				Version:         a.Version,
				Source:          a.Source,
				Status:          a.Status,
			}
			apps = append(apps, &app)
		}
		deployment := marketplace.Deployment{
			Name:         d.Name,
			Creator:      d.Creator,
			Creationdate: d.CreationDate.Format("2006-01-02T15:04:05"),
			Application:  apps,
		}

		deployments = append(deployments, &deployment)
	}

	dea := &marketplace.Deployments{Deployment: deployments}

	de := &marketplace.UnityWebsocketMessage_Deployments{Deployments: dea}
	msg := &marketplace.UnityWebsocketMessage{Content: de}

	data, err := proto.Marshal(msg)
	if err != nil {
		log.WithError(err).Error("Failed to marshal config")
		return err
	}

	websocket.WsManager.SendMessageToAllClients(data)
	return nil
}
