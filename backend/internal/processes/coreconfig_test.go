package processes

import (
	"github.com/stretchr/testify/assert"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/unity-sds/unity-management-console/backend/internal/database/models"
	"github.com/zclconf/go-cty/cty"
	"reflect"
	"strings"
	"testing"
)

type MockDatastore struct{}

type MockParam struct {
	Key   string
	Type  string
	Value string
}

// Mock database fetch operation
func (db *MockDatastore) FetchSSMParams() ([]models.SSMParameters, error) {
	// you can manipulate this part to simulate different situations, like returning an error or different data
	return []models.SSMParameters{
		{
			Key:   "testKey1",
			Type:  "testType1",
			Value: "testValue1",
		},
		{
			Key:   "testKey2",
			Type:  "testType2",
			Value: "testValue2",
		},
	}, nil
}

func (db *MockDatastore) FetchCoreParams() ([]models.CoreConfig, error) {
	// This can be modified to return specific values or errors for testing.
	return nil, nil
}

func (db *MockDatastore) StoreSSMParams(p []config.SSMParameter, owner string) error {
	return nil
}

var testAppConfig = &config.AppConfig{
	Workdir: "/test/workdir",
}

// Test values
var venue = "testVenue"
var project = "testProject"
var publicsubnets = "testPublicSubnets"
var privatesubnets = "testPrivateSubnets"
var ssmParameters = []cty.Value{
	cty.StringVal("testParameter1"),
	cty.StringVal("testParameter2"),
}

func TestGenerateFileStructure(t *testing.T) {
	hclFile := generateFileStructure(testAppConfig, venue, project, publicsubnets, privatesubnets, ssmParameters)

	// Convert the HCL file to a string to easily check the contents
	hclString := string(hclFile.Bytes())

	// Check that the HCL file contains the correct values
	if !strings.Contains(hclString, venue) || !strings.Contains(hclString, project) ||
		!strings.Contains(hclString, publicsubnets) || !strings.Contains(hclString, privatesubnets) {
		t.Errorf("HCL file does not contain correct values")
	}

	// Check that the HCL file contains the ssm parameters
	for _, param := range ssmParameters {
		if !strings.Contains(hclString, param.AsString()) {
			t.Errorf("HCL file does not contain ssm parameter: %s", param.AsString())
		}
	}
}

func TestGenerateSSMParameters(t *testing.T) {
	mockDB := &MockDatastore{}

	ssmParameters, err := generateSSMParameters(mockDB)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedParams := []cty.Value{
		cty.ObjectVal(map[string]cty.Value{
			"name":  cty.StringVal("testKey1"),
			"type":  cty.StringVal("testType1"),
			"value": cty.StringVal("testValue1"),
		}),
		cty.ObjectVal(map[string]cty.Value{
			"name":  cty.StringVal("testKey2"),
			"type":  cty.StringVal("testType2"),
			"value": cty.StringVal("testValue2"),
		}),
	}

	if !reflect.DeepEqual(ssmParameters, expectedParams) {
		t.Errorf("unexpected parameters: got %v, want %v", ssmParameters, expectedParams)
	}
}

func TestGetSSMParameterValueFromDatabase(t *testing.T) {
	db := &MockDatastore{}

	val, err := getSSMParameterValueFromDatabase("testKey1", db)
	assert.NoError(t, err)
	assert.Equal(t, "testValue1", val)

	val, err = getSSMParameterValueFromDatabase("nonexistentKey", db)
	assert.NoError(t, err)
	assert.Equal(t, "", val)
}
