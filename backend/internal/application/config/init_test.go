package config

import (
	"errors"
	"github.com/smartystreets/goconvey/convey"
	"github.com/thomaspoignant/go-feature-flag/ffuser"
	"testing"
)

type MockFFClient struct {
	mockGetValue func(key string) (string, error)
}

func (m *MockFFClient) BoolVariation(s string, u ffuser.User, f bool) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockFFClient) Close() {
	// Add any behavior you need here, or leave it as a no-op.
	return
}
func TestInitFeatureFlags(t *testing.T) {
	convey.Convey("Given a feature flag configuration", t, func() {
		mockClient := &MockFFClient{
			mockGetValue: func(key string) (string, error) {
				if key == "expectedKey" {
					return "expectedValue", nil
				}
				return "", errors.New("unexpected key")
			},
		}
		FFClient = mockClient // Let's assume FFClient is the variable of the interface FeatureFlagClient type in your config package

		convey.Convey("When initFeatureFlags is called", func() {
			f := initFeatureFlags()

			convey.Convey("Then a new feature flag client should be created", func() {
				convey.ShouldNotBeNil(f)
			})
		})
	})
}
