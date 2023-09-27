package processes

import (
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/aws"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"github.com/unity-sds/unity-management-console/backend/internal/websocket"
)

func ProcessSimpleMessage(message *marketplace.SimpleMessage, conf *config.AppConfig, store database.Datastore, wsmgr *websocket.WebSocketManager, userid string) ([]byte, error) {
	if message.Operation == "request config" {
		log.Info("Request Config received")
		return fetchConfig(conf, store)
	} else if message.Operation == "request parameters" {
		log.Info("Request Parameters received")
		return fetchParameters(conf)
	} else if message.Operation == "update config" {
		log.Info("Update config received")
		//return updateParameters()
	} else if message.Operation == "request all applications" {
		log.Info("Request all applications received")
		err := fetchAllApplications(store)
		return nil, err
	} else if message.Operation == "uninstall application" {
		log.Info("Request to uninstall application")
		err := uninstallApplication(message.Payload, conf, store)
		return nil, err
	} else if message.Operation == "uninstall deployment" {
		log.Info("Request to uninstall deployment")
		return uninstallDeployment(message.Payload)
	} else if message.Operation == "reapply application" {
		log.Info("Request to reapply application")
		err := reapplyApplication(message.Payload, conf, store, wsmgr, userid)
		return nil, err
	}
	return nil, nil
}

func UpdateParameters(params *marketplace.Parameters, store database.Datastore, appconf *config.AppConfig, wsmgr *websocket.WebSocketManager, userid string) {

	log.Info("Storing parameters")
	var parr []config.SSMParameter
	for _, p := range params.Parameterlist {
		if p.Name != "" && p.Value != "" {
			np := config.SSMParameter{
				Name:  p.Name,
				Type:  p.Type,
				Value: p.Value,
			}
			parr = append(parr, np)

		}
	}
	log.Infof("Saving %v parameters", len(parr))
	err := store.StoreSSMParams(parr, "test")
	if err != nil {
		log.WithError(err).Error("Error storing parameters")
		return
	}

	store.AddToAudit(application.Config_Updated, "test")

	err = UpdateCoreConfig(appconf, store, wsmgr, userid)
	if err != nil {
		log.WithError(err).Error("Error updating config")
		return
	}

}

func fetchParameters(conf *config.AppConfig) ([]byte, error) {
	db, err := database.NewGormDatastore()

	params, err := db.FetchSSMParams()

	ssm, err := aws.ReadSSMParameters(params)

	paramwrap := marketplace.UnityWebsocketMessage_Parameters{Parameters: ssm}
	msg := &marketplace.UnityWebsocketMessage{Content: &paramwrap}
	data, err := proto.Marshal(msg)
	if err != nil {
		log.WithError(err).Error("Failed to marshal config")
		return nil, err
	}

	return data, nil
}
func fetchConfig(conf *config.AppConfig, store database.Datastore) ([]byte, error) {

	//coreconf, err := store.FetchCoreParams()
	//if err != nil {
	//	log.WithError(err).Error("Error fetching core config")
	//	return nil, err
	//}
	pub, priv, err := aws.FetchSubnets()
	if err != nil {
		log.WithError(err).Error("Error fetching subnets")
		return nil, err
	}

	netconfig := marketplace.Config_NetworkConfig{
		Publicsubnets:  pub,
		Privatesubnets: priv,
	}

	appConfig := marketplace.Config_ApplicationConfig{
		GithubToken:      conf.GithubToken,
		MarketplaceOwner: conf.MarketplaceOwner,
		MarketplaceUser:  conf.MarketplaceRepo,
	}
	auditline := application.Config_Updated
	audit, err := store.FindLastAuditLineByOperation(auditline)
	genconfig := &marketplace.Config{

		ApplicationConfig: &appConfig,
		NetworkConfig:     &netconfig,
		Lastupdated:       audit.CreatedAt.Format("2006-01-02T15:04:05.000"),
		Updatedby:         audit.Owner,
	}

	mpcfg := marketplace.UnityWebsocketMessage_Config{Config: genconfig}

	msg := &marketplace.UnityWebsocketMessage{Content: &mpcfg}

	log.WithFields(log.Fields{
		"Config": genconfig,
	}).Info("Config Generated")

	data, err := proto.Marshal(msg)
	if err != nil {
		log.WithError(err).Error("Failed to marshal config")
		return nil, err
	}

	return data, nil
}

func uninstallApplication(name string, conf *config.AppConfig, store database.Datastore) error {

	log.Infof("Uninstalling application %s", name)

	return UninstallApplication(name, conf, store)
}

func reapplyApplication(name string, conf *config.AppConfig, store database.Datastore, wsmgr *websocket.WebSocketManager, userid string) error {
	log.Infof("Repplying application %s", name)

	return ReapplyApplication(name, conf, store, wsmgr, userid)
}

func uninstallDeployment(name string) ([]byte, error) {
	log.Warn("Uninstall Deployment not yet implemented")

	return nil, nil
}
