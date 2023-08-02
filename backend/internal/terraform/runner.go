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

type TerraformExecutor interface {
	NewTerraform(string, string) (*tfexec.Terraform, error)
	Init(context.Context, ...tfexec.InitOption) error
	Plan(context.Context, ...tfexec.PlanOption) (bool, error)
	Apply(context.Context, ...tfexec.ApplyOption) error
	SetStdout(io.Writer)
	SetStderr(io.Writer)
	SetLogger(*log.Logger)
}

type RealTerraformExecutor struct {
	tf *tfexec.Terraform
}

func (r *RealTerraformExecutor) NewTerraform(dir string, execPath string) (*tfexec.Terraform, error) {
	tf, err := tfexec.NewTerraform(dir, execPath)
	if err != nil {
		return nil, err
	}

	r.tf = tf
	return tf, nil
}

func (r *RealTerraformExecutor) Init(ctx context.Context, opts ...tfexec.InitOption) error {
	return r.tf.Init(ctx, opts...)
}

func (r *RealTerraformExecutor) Plan(ctx context.Context, opts ...tfexec.PlanOption) (bool, error) {
	return r.tf.Plan(ctx, opts...)
}

func (r *RealTerraformExecutor) Apply(ctx context.Context, opts ...tfexec.ApplyOption) error {
	return r.tf.Apply(ctx, opts...)
}

func (r *RealTerraformExecutor) SetStdout(w io.Writer) {
	r.tf.SetStdout(w)
}

func (r *RealTerraformExecutor) SetStderr(w io.Writer) {
	r.tf.SetStderr(w)
}

func (r *RealTerraformExecutor) SetLogger(l *log.Logger) {
	r.tf.SetLogger(l)
}
func RunTerraform(appconf *config.AppConfig, wsmgr *ws.WebSocketManager, id string, executor TerraformExecutor) {
	bucket := fmt.Sprintf("bucket=%s", appconf.BucketName)
	key := fmt.Sprintf("key=%s", "default")
	region := fmt.Sprintf("region=%s", appconf.AWSRegion)

	tf, err := executor.NewTerraform(appconf.Workdir, "/usr/bin/terraform")
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
	err = executor.Init(context.Background(), tfexec.Upgrade(true), tfexec.BackendConfig(bucket), tfexec.BackendConfig(key), tfexec.BackendConfig(region))

	if err != nil {
		log.WithError(err).Error("error initialising terraform")
	}

	change, err := executor.Plan(context.Background())

	if err != nil {
		log.WithError(err).Error("error running plan")
	}

	fmt.Printf("change: %v", change)

	if change {
		err = executor.Apply(context.Background())

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
