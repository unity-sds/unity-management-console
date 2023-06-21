package processes

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"github.com/unity-sds/unity-control-plane/backend/internal/application/config"
	"github.com/unity-sds/unity-control-plane/backend/internal/database"
	"github.com/unity-sds/unity-control-plane/backend/internal/database/models"
	"github.com/unity-sds/unity-control-plane/backend/internal/marketplace"
	ws "github.com/unity-sds/unity-control-plane/backend/internal/websocket"
	"os"
	"path/filepath"
	"testing"
)

type mockDB struct{}

type MockStore struct{}

var conf config.AppConfig

func (m *MockStore) FetchCoreParams() ([]models.CoreConfig, error) {
	// return test data
	return []models.CoreConfig{}, nil
}
func (db *mockDB) FetchConfig() ([]models.CoreConfig, error) {
	return []models.CoreConfig{
		{Key: "project", Value: "testProject"},
		{Key: "venue", Value: "testVenue"},
		{Key: "privateSubnets", Value: "testPrivateSubnets"},
		{Key: "publicSubnets", Value: "testPublicSubnets"},
	}, nil
}

type MockActRunner struct {
	//mockRunAct func(path string, inputs, env, secrets map[string]string, conn *websocket.Conn) error
	RunActFn                func(path string, inputs, env, secrets map[string]string, conn *websocket.Conn) error
	UpdateCoreConfigFn      func(conn *websocket.Conn, store database.Datastore) error
	InstallMarketplaceAppFn func(conn *websocket.Conn, store database.Datastore, meta string) error
}

func (m *MockActRunner) RunAct(path string, inputs, env, secrets map[string]string, conn *websocket.Conn) error {
	if m.RunActFn != nil {
		return m.RunActFn(path, inputs, env, secrets, conn)
	}
	return nil
}

func (m *MockActRunner) UpdateCoreConfig(conn *websocket.Conn, store database.Datastore) error {
	if m.UpdateCoreConfigFn != nil {
		return m.UpdateCoreConfigFn(conn, store)
	}

	return nil
}

func (m *MockActRunner) InstallMarketplaceApplication(conn *websocket.Conn, store database.Datastore, meta string) error {
	if m.InstallMarketplaceAppFn != nil {
		return m.InstallMarketplaceAppFn(conn, store, meta)
	}

	return nil
}

func TestUpdateCoreConfig(t *testing.T) {
	Convey("Given a mock sql connection and a websocket connection", t, func() {

		Convey("And a mock implementation of fetchCoreParams", func() {

			Convey("When UpdateCoreConfig is called", func() {
				// temporarily point the global DB variable to our mock
				mockStore := &MockStore{}
				mockRunner := &MockActRunner{
					RunActFn: func(path string, inputs, env, secrets map[string]string, conn *websocket.Conn) error {
						// Do something here to mock the act.RunAct behavior
						return nil
					},
				}

				err := mockRunner.UpdateCoreConfig(nil, mockStore)

				Convey("Then the expectations should be met", func() {
					So(err, ShouldBeNil)
				})
			})
		})
	})
}

func fetchConfig() {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Errorf("Error fetching home directory: %v", err)
		return
	}

	configdir := filepath.Join(dir, ".unity")

	if _, err := os.Stat(configdir); os.IsNotExist(err) {
		errDir := os.MkdirAll(configdir, 0755)
		if errDir != nil {
			log.Errorf("Error creating directory: %v", errDir)
			return
		}
	}

	// Search config in home directory with name ".cobra" (without extension).
	viper.AddConfigPath(configdir)
	viper.SetConfigType("yaml")
	viper.SetConfigName("unity")
	viper.SetDefault("GithubToken", "unset")
	viper.SetDefault("MarketplaceURL", "unset")
	viper.SetDefault("WorkflowBasePath", "unset")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			file, createErr := os.Create(filepath.Join(configdir, "unity.yaml"))
			if createErr != nil {
				log.Fatalf("Failed to create config file: %v", createErr)
			}
			defer file.Close()
			log.Infof("Created config file: %s", viper.ConfigFileUsed())
		} else {
			// Config file was found but another error was produced
			log.Errorf("Failed to read config file: %v", err)
		}
	}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Errorf("Unable to decode into struct, %v", err)
	}
}
func TestRunSPSDemo(t *testing.T) {
	r := ActRunnerImpl{}
	mockStore := &MockStore{}

	fetchConfig()
	meta := "{\n\t\"metadata\": {\n\t\t\"metadataversion\": \"unity-cs-0.1\",\n\t\t\"exectarget\": \"act\",\n\t\t\"deploymentname\": \"managementdashboard\",\n\t\t\"services\": [\n\t\t\t{\"name\":\"ryantestdeploy\",\"source\":\"unity-sds/unity-sps-prototype\",\"version\":\"xxx\",\"branch\":\"main\"}\n\t\t],\n\t\t\"extensions\":{\n\t\t\t\"kubernetes\":{\n\t\t\t\t\"clustername\":\"unity-sps-managementdashboard\",\n\t\t\t\t\"owner\":\"ryan\",\n\t\t\t\t\"projectname\":\"testproject\",\n\t\t\t\t\"nodegroups\":{\n\t\t\t\t\t\"group1\": {\n\t\t\t\t\t\t\"instancetype\": \"m5.xlarge\",\n\t\t\t\t\t\t\"nodecount\":\"1\"\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t}\n\t}\n}"
	err := r.InstallMarketplaceApplication(nil, mockStore, meta, conf, "")
	log.Errorf("Error: %v", err)
}

func TestRun(t *testing.T) {
	r := ActRunnerImpl{}
	mockStore := &MockStore{}

	ng1 := marketplace.Extensions_Nodegroups{
		Name:         "my ng",
		Instancetype: "m5.xlarge",
		Nodecount:    "5",
	}

	narr := []*marketplace.Extensions_Nodegroups{&ng1}

	eks := marketplace.Extensions_Eks{
		Clustername: "test cluster",
		Owner:       "tom",
		Projectname: "testing",
		Nodegroups:  narr,
	}

	msg := marketplace.Extensions{
		Eks: &eks,
	}
	m := models.AppInstall{
		Name:      "unity-eks",
		Version:   "0.1",
		Variables: nil,
	}
	c := models.InstallConfig{
		Install:    []models.AppInstall{m},
		Extensions: msg,
	}
	payload := []models.InstallConfig{c}
	rec := ws.InstallMessage{
		Action:  "install software",
		Payload: payload,
	}
	r.TriggerInstall(nil, mockStore, rec, conf)
}
