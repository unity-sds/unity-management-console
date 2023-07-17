package terraform

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	"os"
)

func RunTerraform(appconf *config.AppConfig) {
	bucket := fmt.Sprintf("bucket=%s", appconf.BucketName)
	key := fmt.Sprintf("key=%s", "default")
	region := fmt.Sprintf("region=%s", appconf.AWSRegion)

	tf, err := tfexec.NewTerraform(appconf.Workdir, "/usr/bin/terraform")
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	tf.SetStdout(os.Stdout)
	tf.SetStderr(os.Stderr)
	tf.SetLogger(log.StandardLogger())
	err = tf.Init(context.Background(), tfexec.Upgrade(true), tfexec.BackendConfig(bucket), tfexec.BackendConfig(key), tfexec.BackendConfig(region))

	if err != nil {
		log.WithError(err).Error("error initialising terraform")
	}

	change, err := tf.Plan(context.Background())

	if err != nil {
		log.WithError(err).Error("error running plan")
	}

	fmt.Printf("change: %v", change)

	if change {
		err = tf.Apply(context.Background())

		if err != nil {
			log.WithError(err).Error("error running apply")
		}

	}
}
