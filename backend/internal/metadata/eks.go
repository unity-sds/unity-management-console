package metadata

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-control-plane/backend/internal/marketplace"
	"google.golang.org/protobuf/encoding/protojson"
)

func GenerateEKSMetadata(extensions *marketplace.Install_Extensions) ([]byte, error) {

	m := protojson.MarshalOptions{
		UseProtoNames: true,
	}

	// Convert the protobuf message to JSON.
	jsonStr, err := m.Marshal(extensions)
	if err != nil {
		log.Errorf("unable to decode json: %s", jsonStr)
		return []byte{}, err
	}

	// We will unmarshal to a map of string keys and arbitrary types.
	var result map[string]interface{}
	err = json.Unmarshal(jsonStr, &result)
	if err != nil {
		log.Errorf("Error occurred during unmarshalling: %v", err)
		return []byte{}, err
	}

	// Extract the "eks" object.
	eks, ok := result["eks"]
	if !ok {
		log.Fatalf("'eks' key not found in the JSON")
		return []byte{}, err
	}

	// Convert the "eks" object to a map
	eksMap, ok := eks.(map[string]interface{})
	if !ok {
		log.Fatalf("'eks' value is not a map")
		return []byte{}, err
	}

	// Extract the "nodegroups" slice
	nodegroups, ok := eksMap["nodegroups"]
	if !ok {
		log.Fatalf("'nodegroups' key not found in the 'eks' map")
		return []byte{}, err
	}

	// Convert the "nodegroups" slice to a slice of maps
	nodegroupsSlice, ok := nodegroups.([]interface{})
	if !ok {
		log.Fatalf("'nodegroups' value is not a slice")
		return []byte{}, err
	}

	// Create a new map where each key is the "name" field from the slice elements
	newNodegroups := make(map[string]interface{})
	for _, ng := range nodegroupsSlice {
		ngMap, ok := ng.(map[string]interface{})
		if !ok {
			log.Fatalf("Element of 'nodegroups' slice is not a map")
			return []byte{}, err
		}
		name, ok := ngMap["name"].(string)
		if !ok {
			log.Fatalf("'name' key not found in the 'nodegroups' element map or is not a string")
			return []byte{}, err
		}
		delete(ngMap, "name") // remove the name key-value pair
		newNodegroups[name] = ngMap
	}

	// Replace the old "nodegroups" slice with the new map
	eksMap["nodegroups"] = newNodegroups

	// Marshal the "eks" object back to a JSON string.
	eksJson, err := json.Marshal(eksMap)
	if err != nil {
		log.Fatalf("Error occurred during marshalling: %v", err)
		return []byte{}, err
	}

	log.Infof("EKS Meta: %v", string(eksJson))

	return eksJson, nil
}
