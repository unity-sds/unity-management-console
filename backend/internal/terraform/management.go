package terraform

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclwrite"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"github.com/zclconf/go-cty/cty"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func AddApplicationToStack(appConfig *config.AppConfig, location string, meta *marketplace.MarketplaceMetadata) error {
	rand.Seed(time.Now().UnixNano())

	s := String(8)
	hclFile := hclwrite.NewEmptyFile()

	err := os.MkdirAll(filepath.Join(appConfig.Workdir, "workspace"), 0755)
	if err != nil {
		log.WithError(err).Error("Could not create workspace directory")
	}
	// create new file on system
	tfFile, err := os.Create(filepath.Join(appConfig.Workdir, "workspace", fmt.Sprintf("%v%v%v", meta.Name, s, ".tf")))
	if err != nil {
		log.WithError(err).Error("Problem creating tf file")
		return err
	}
	path := filepath.Join(location, meta.WorkDirectory)
	// initialize the body of the new file object
	rootBody := hclFile.Body()

	cloudenv := rootBody.AppendNewBlock("module", []string{meta.Name})
	cloudenvBody := cloudenv.Body()
	cloudenvBody.SetAttributeValue("source", cty.StringVal(path))

	_, err = tfFile.Write(hclFile.Bytes())
	if err != nil {
		log.WithError(err).Error("error writing hcl file")
		return err
	}
	return nil
}
