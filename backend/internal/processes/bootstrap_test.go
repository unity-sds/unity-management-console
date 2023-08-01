package processes

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"path/filepath"
	"testing"
)

func TestWriteInitTemplate(t *testing.T) {
	// Use an in-memory filesystem for testing
	fs := afero.NewMemMapFs()

	// Given
	workdir := "/some/test/dir"
	appConfig := &config.AppConfig{Workdir: workdir}

	// When
	writeInitTemplate(fs, appConfig)

	// Then
	filePath := filepath.Join(workdir, "provider.tf")

	// Check if file exists
	exists, err := afero.Exists(fs, filePath)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, exists, "File 'provider.tf' should have been created")

	// Check file content
	expectedContent := `terraform {
  backend "s3" {
    dynamodb_table = "terraform_state"
  }
}`
	content, err := afero.ReadFile(fs, filePath)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedContent, string(content), "File content does not match")
}
