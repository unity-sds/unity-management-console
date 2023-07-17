package processes

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/unity-sds/unity-management-console/backend/internal/action"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/database/models"
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

func (m *MockStore) FetchSSMParams() ([]models.SSMParameters, error) {
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
	action.ActRunner
}

func (m *MockActRunner) RunAct(path string, inputs, env, secrets map[string]string, conn *websocket.Conn, appConfig config.AppConfig) error {
	return nil
}

func TestUpdateCoreConfig(t *testing.T) {
	Convey("Given a mock sql connection and a websocket connection", t, func() {

		Convey("And a mock implementation of fetchCoreParams", func() {

			Convey("When UpdateCoreConfig is called", func() {
				// temporarily point the global DB variable to our mock
				//mockStore := new(MockStore)
				//
				//mockData := []models.SSMParameters{{
				//	Key:   "/unity/core/project",
				//	Value: "project-value",
				//	Type:  "String",
				//}, {
				//	Key:   "/unity/core/venue",
				//	Value: "venue-value",
				//	Type:  "String",
				//}}
				//
				//mockStore.On("FetchSSMParams").Return(mockData, nil)
				//
				//mockRunner := &MockActRunner{}
				//
				//conn := new(websocket.Conn) // You might want to mock this as well if your RunAct function interacts with it.
				fetchConfig()
				config := conf // Replace with the appropriate type and values.

				Convey("When UpdateCoreConfig is called", func() {
					err := UpdateCoreConfig(&config)

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
