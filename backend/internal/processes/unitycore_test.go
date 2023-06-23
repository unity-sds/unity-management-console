package processes

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/unity-sds/unity-control-plane/backend/internal/application/config"
	"github.com/unity-sds/unity-control-plane/backend/internal/database/models"
	"github.com/unity-sds/unity-control-plane/backend/internal/marketplace"
	"os"
	"path/filepath"
	"testing"
)

type mockDB struct{}

type MockStore struct {
	mock.Mock
}

func (m *MockStore) StoreSSMParams(p []config.SSMParameter, owner string) error {
	
	return nil
}

func (m *MockStore) FetchSSMParams() ([]models.SSMParameters, error) { // Replace YourProject.YourPackage.SSMParam with the appropriate type
	args := m.Called()
	return args.Get(0).([]models.SSMParameters), args.Error(1)
}

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
	ActRunner
}

func (m *MockActRunner) RunAct(path string, inputs, env, secrets map[string]string, conn *websocket.Conn, appConfig config.AppConfig) error {
	return nil
}

func TestUpdateCoreConfig(t *testing.T) {
	Convey("Given a mock sql connection and a websocket connection", t, func() {

		Convey("And a mock implementation of fetchCoreParams", func() {

			Convey("When UpdateCoreConfig is called", func() {
				// temporarily point the global DB variable to our mock
				mockStore := new(MockStore)

				mockData := []models.SSMParameters{{
					Key:   "/unity/core/project",
					Value: "project-value",
					Type:  "String",
				}, {
					Key:   "/unity/core/venue",
					Value: "venue-value",
					Type:  "String",
				}}

				mockStore.On("FetchSSMParams").Return(mockData, nil)

				mockRunner := &MockActRunner{}

				conn := new(websocket.Conn) // You might want to mock this as well if your RunAct function interacts with it.
				fetchConfig()
				config := conf // Replace with the appropriate type and values.

				Convey("When UpdateCoreConfig is called", func() {
					err := UpdateCoreConfig(conn, mockStore, config, mockRunner)

					Convey("Then no error should be returned", func() {
						So(err, ShouldBeNil)
					})

					// Add more assertions as per your requirements.
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

	fetchConfig()
	meta := "{\n\t\"metadata\": {\n\t\t\"metadataversion\": \"unity-cs-0.1\",\n\t\t\"exectarget\": \"act\",\n\t\t\"deploymentname\": \"managementdashboard\",\n\t\t\"services\": [\n\t\t\t{\"name\":\"ryantestdeploy\",\"source\":\"unity-sds/unity-sps-prototype\",\"version\":\"xxx\",\"branch\":\"main\"}\n\t\t],\n\t\t\"extensions\":{\n\t\t\t\"kubernetes\":{\n\t\t\t\t\"clustername\":\"unity-sps-managementdashboard\",\n\t\t\t\t\"owner\":\"ryan\",\n\t\t\t\t\"projectname\":\"testproject\",\n\t\t\t\t\"nodegroups\":{\n\t\t\t\t\t\"group1\": {\n\t\t\t\t\t\t\"instancetype\": \"m5.xlarge\",\n\t\t\t\t\t\t\"nodecount\":\"1\"\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t}\n\t}\n}"
	err := InstallMarketplaceApplication(nil, meta, conf, "", r)
	log.Errorf("Error: %v", err)
}

func TestRun(t *testing.T) {
	r := ActRunnerImpl{}
	mockStore := &MockStore{}

	ng1 := marketplace.Install_Extensions_Nodegroups{
		Name:         "my ng",
		Instancetype: "m5.xlarge",
		Nodecount:    "5",
	}

	narr := []*marketplace.Install_Extensions_Nodegroups{&ng1}

	eks := marketplace.Install_Extensions_Eks{
		Clustername: "test cluster",
		Owner:       "tom",
		Projectname: "testing",
		Nodegroups:  narr,
	}

	msg := marketplace.Install_Extensions{
		Eks: &eks,
	}
	m := marketplace.Install_Applications{
		Name:      "unity-eks",
		Version:   "0.1",
		Variables: nil,
	}

	c := marketplace.Install{
		Applications: &m,
		Extensions:   &msg,
	}

	TriggerInstall(nil, mockStore, c, conf, r)
}
