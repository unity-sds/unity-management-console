package terraform

import (
	"context"
	"github.com/hashicorp/terraform-exec/tfexec"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	ws "github.com/unity-sds/unity-management-console/backend/internal/websocket"
	"io"
	"testing"
)

type MockTerraformExecutor struct {
	mock.Mock
}

func (m *MockTerraformExecutor) NewTerraform(dir string, execPath string) (*tfexec.Terraform, error) {
	args := m.Called(dir, execPath)
	return args.Get(0).(*tfexec.Terraform), args.Error(1)
}

func (m *MockTerraformExecutor) Init(ctx context.Context, opts ...tfexec.InitOption) error {
	args := m.Called(ctx, opts)
	return args.Error(0)
}

func (m *MockTerraformExecutor) Plan(ctx context.Context, opts ...tfexec.PlanOption) (bool, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockTerraformExecutor) Apply(ctx context.Context, opts ...tfexec.ApplyOption) error {
	args := m.Called(ctx, opts)
	return args.Error(0)
}

func (m *MockTerraformExecutor) SetStdout(w io.Writer) {
	m.Called(w)
}

func (m *MockTerraformExecutor) SetStderr(w io.Writer) {
	m.Called(w)
}

func (m *MockTerraformExecutor) SetLogger(l *log.Logger) {
	m.Called(l)
}

func TestRunTerraform(t *testing.T) {
	appconf := &config.AppConfig{
		// populate with the appropriate values
	}

	id := "test"
	wsmgr := &ws.WebSocketManager{
		// populate with the appropriate values
	}

	mockExecutor := &MockTerraformExecutor{}
	mockExecutor.On("NewTerraform", mock.Anything, mock.Anything, mock.Anything).Return(&tfexec.Terraform{}, nil)
	mockExecutor.On("Init", mock.Anything, mock.Anything).Return(nil)
	mockExecutor.On("Plan", mock.Anything, mock.Anything).Return(true, nil)
	mockExecutor.On("Apply", mock.Anything, mock.Anything).Return(nil)

	RunTerraform(appconf, wsmgr, id, mockExecutor)

	mockExecutor.AssertExpectations(t)
}
