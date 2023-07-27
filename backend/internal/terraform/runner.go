package terraform

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	ws "github.com/unity-sds/unity-management-console/backend/internal/websocket"
	"io"
	"os"
)

var stdoutBuffer bytes.Buffer
var stderrBuffer bytes.Buffer

func RunTerraform(appconf *config.AppConfig, wsmgr *ws.WebSocketManager, id string) {
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

	if wsmgr != nil {
		stopCapture := CaptureOutput(id, wsmgr)
		defer stopCapture() // Ensure that we stop capturing even if an error occurs
	}
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

func CaptureOutput(id string, wsmgr *ws.WebSocketManager) (stopCapture func()) {

	stdoutReader, stdoutWriter, _ := os.Pipe()
	stderrReader, stderrWriter, _ := os.Pipe()

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	os.Stdout = stdoutWriter
	os.Stderr = stderrWriter

	done := make(chan bool, 2)
	readAndCapture := func(reader *os.File, buffer *bytes.Buffer) {
		buf := make([]byte, 1024)
		for {
			n, err := reader.Read(buf)
			if err != nil {
				if err != io.EOF {
					// log.Println("read error:", err)
				}
				break
			}

			if _, err := buffer.Write(buf[:n]); err != nil {
				// log.Println("buffer write:", err)
				return
			}
			wsmgr.SendMessageToUserID(id, buf[:n])
		}
		done <- true
	}

	go readAndCapture(stdoutReader, &stdoutBuffer)
	go readAndCapture(stderrReader, &stderrBuffer)

	return func() {
		stdoutWriter.Close()
		stderrWriter.Close()
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		<-done
		<-done
	}
}
