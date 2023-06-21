package metadata

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-control-plane/backend/internal/marketplace"
	"google.golang.org/protobuf/encoding/protojson"
)

func GenerateEKSMetadata(extensions *marketplace.Install_Extensions) ([]byte, error) {
	//{
	//	"clustername": "testng2",
	//	"owner": "tom",
	//	"projectname": "testproject",
	//	"nodegroups": {
	//		"group1": {
	//			"instancetype": "m5.xlarge",
	//			"nodecount": "1"
	//		},
	//		"group2": {
	//			"instancetype": "m5.xlarge",
	//			"nodecount": "1"
	//		}
	//	}
	//}

	m := protojson.MarshalOptions{
		UseProtoNames: true,
	}

	// Convert the protobuf message to JSON.
	jsonStr, err := m.Marshal(extensions)
	if err != nil {
		log.Errorf("unable to decode json: %s", jsonStr)
		return []byte{}, err
	}
	//eksstr := extensions.Eks.String()

	// We will unmarshal to a map of string keys and arbitrary types.
	var result map[string]interface{}
	err = json.Unmarshal(jsonStr, &result)
	if err != nil {
		log.Errorf("Error occurred during unmarshalling: %v", err)
		return []byte{}, err
	}

	// Convert the protobuf message to JSON.
	log.Infof("converting extensions to metadata, %s", jsonStr)

	// Extract the "eks" object.
	eks, ok := result["eks"]
	if !ok {
		log.Fatalf("'eks' key not found in the JSON")
		return []byte{}, err
	}

	// Marshal the "eks" object back to a JSON string.
	eksJson, err := json.Marshal(eks)
	if err != nil {
		log.Fatalf("Error occurred during marshalling: %v", err)
		return []byte{}, err
	}

	return eksJson, nil

}
