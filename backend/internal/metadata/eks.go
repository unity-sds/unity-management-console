package metadata

import (
	"github.com/golang/protobuf/proto"
	"github.com/unity-sds/unity-cs-manager/marketplace"
)

func GenerateApplicationMetadata(appname string, install *marketplace.Install, meta *marketplace.MarketplaceMetadata) ([]byte, error) {

	svc := marketplace.ActionMeta_Services{
		Name:    install.Applications.Name,
		Source:  meta.Package,
		Version: meta.Version,
		Branch:  "main",
	}
	actionmeta := &marketplace.ActionMeta{
		MetadataVersion: "unity-cs-0.1",
		Exectarget:      "act",
		DeploymentName:  appname,
		Services:        []*marketplace.ActionMeta_Services{&svc},
		Extensions:      nil,
	}

	return proto.Marshal(actionmeta)
}
