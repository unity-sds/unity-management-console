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

	log.Info("Received Project: " + conf.Project)
	log.Info("Received Venue: " + conf.Venue)

	appConfig := marketplace.Config_ApplicationConfig{
		GithubToken:      conf.GithubToken,
		MarketplaceOwner: conf.MarketplaceOwner,
		MarketplaceUser:  conf.MarketplaceRepo,
		Project:          conf.Project,
		Venue:            conf.Venue,
	}

	log.Info("Received App Config: " + appConfig)
	auditline := application.Config_Updated
	audit, err := store.FindLastAuditLineByOperation(auditline)

	auditline = application.Bootstrap_Unsuccessful
	bootstrapfailed, err := store.FindLastAuditLineByOperation(auditline)

	auditline = application.Bootstrap_Successful
	bootstrapsuccess, err := store.FindLastAuditLineByOperation(auditline)

	bsoutput := ""
	if bootstrapsuccess.Owner != "" {
		bsoutput = "complete"
	} else if bootstrapfailed.Owner != "" {
		bsoutput = "failed"
	}
	genconfig := &marketplace.Config{

		ApplicationConfig: &appConfig,
		NetworkConfig:     &netconfig,
		Lastupdated:       audit.CreatedAt.Format("2006-01-02T15:04:05.000"),
		Updatedby:         audit.Owner,
		Bootstrap:         bsoutput,
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

func reapplyApplication(name string, conf *config.AppConfig, store database.Datastore, wsmgr *websocket.WebSocketManager, userid string) error {
	log.Infof("Repplying application %s", name)

	return ReapplyApplication(name, conf, store, wsmgr, userid)
}
