package main

import (
	"github.com/stretchr/testify/mock"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/database/models"
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
