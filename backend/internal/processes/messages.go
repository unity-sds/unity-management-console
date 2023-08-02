package processes

import (
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/aws"
	"github.com/unity-sds/unity-management-console/backend/internal/database"
	"github.com/unity-sds/unity-management-console/backend/internal/websocket"
)

func ProcessSimpleMessage(message *marketplace.SimpleMessage, conf *config.AppConfig) ([]byte, error) {
	if message.Operation == "request config" {
		log.Info("Request Config received")
		return fetchConfig(conf)
	} else if message.Operation == "request parameters" {
		log.Info("Request Parameters received")
		return fetchParameters(conf)
	}
	return nil, nil
}

func UpdateParameters(params *marketplace.Parameters, store database.Datastore, appconf *config.AppConfig, wsmgr *websocket.WebSocketManager, userid string) {

	log.Info("Storing parameters")
	var parr []config.SSMParameter
	for _, p := range params.Parameterlist {
		np := config.SSMParameter{
			Name:  p.Name,
			Type:  p.Type,
			Value: p.Value,
		}
		parr = append(parr, np)
	}
	log.Infof("Saving %v parameters", len(parr))
	err := store.StoreSSMParams(parr, "test")
	if err != nil {
		log.WithError(err).Error("Error storing parameters")
		return
	}

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
func fetchConfig(conf *config.AppConfig) ([]byte, error) {

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
	genconfig := &marketplace.Config{

		ApplicationConfig: &appConfig,
		NetworkConfig:     &netconfig,
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
