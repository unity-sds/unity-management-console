package main

import (
	"github.com/stretchr/testify/mock"
	"github.com/unity-sds/unity-control-plane/backend/internal/application/config"
	"github.com/unity-sds/unity-control-plane/backend/internal/database/models"
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

func TestGateway(t *testing.T) {

	appConfig := config.AppConfig{
		GithubToken:          "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",
		MarketplaceUrl:       "",
		WorkflowBasePath:     "/home/barber/Projects/unity-cs-infra/.github/workflows",
		AWSRegion:            "",
		BucketName:           "",
		Workdir:              "",
		DefaultSSMParameters: nil,
	}
	mockStore := new(MockStore)

	installGateway(mockStore, appConfig)
}
