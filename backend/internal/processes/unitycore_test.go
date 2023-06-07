package processes

import (
	"github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/unity-sds/unity-control-plane/backend/internal/database"
	"github.com/unity-sds/unity-control-plane/backend/internal/database/models"
	"testing"
)

type mockDB struct{}

type MockStore struct{}

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
	RunActFn func(path string, inputs, env, secrets map[string]string, conn *websocket.Conn) error
	UpdateCoreConfigFn func(conn *websocket.Conn, store database.Datastore) error
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
