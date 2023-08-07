package terraform

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/terraform-exec/tfexec"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-cs-manager/marketplace"
	"github.com/unity-sds/unity-management-console/backend/internal/application/config"
	ws "github.com/unity-sds/unity-management-console/backend/internal/websocket"
	"io"
	"os"
	"path/filepath"
	"strings"
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

type wsWriter struct {
	builder strings.Builder
	wsmgr   *ws.WebSocketManager
	userid  string
	level   string
}

func (w *wsWriter) Write(p []byte) (n int, err error) {
	// Here, we're just writing the bytes to a strings.Builder, but you could do anything you want.
	n, err = w.builder.Write(p)
	m := marketplace.LogLine{
		Line:      string(p),
		Level:     w.level,
		Timestamp: "",
		Type:      "",
	}
	mes := marketplace.UnityWebsocketMessage_Logs{Logs: &m}
	mess := &marketplace.UnityWebsocketMessage{Content: &mes}
	data, err := proto.Marshal(mess)
	if err != nil {
		log.WithError(err).Error("Failed to marshal log line")
		return
	}
	w.wsmgr.SendMessageToUserID(w.userid, data)

	return n, err
}

func RunTerraform(appconf *config.AppConfig, wsmgr *ws.WebSocketManager, id string, executor TerraformExecutor) {
	bucket := fmt.Sprintf("bucket=%s", appconf.BucketName)
	key := fmt.Sprintf("key=%s", "default")
	region := fmt.Sprintf("region=%s", appconf.AWSRegion)

	p := filepath.Join(appconf.Workdir, "workspace")
	tf, err := executor.NewTerraform(p, "/usr/bin/terraform")
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}
	wwsWriter := &wsWriter{
		userid: id,
		wsmgr:  wsmgr,
		level:  "INFO",
	}

	wwserrWriter := &wsWriter{
		userid: id,
		wsmgr:  wsmgr,
		level:  "ERROR",
	}
	writerStdout := io.MultiWriter(os.Stdout, wwsWriter)
	writerStderr := io.MultiWriter(os.Stderr, wwserrWriter)

	tf.SetStdout(writerStdout)
	tf.SetStderr(writerStderr)
	tf.SetLogger(log.StandardLogger())

	//if wsmgr != nil {
	//	stopCapture := CaptureOutput(id, wsmgr)
	//	defer stopCapture() // Ensure that we stop capturing even if an error occurs
	//}
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

//func CaptureOutput(id string, wsmgr *ws.WebSocketManager) (stopCapture func()) {
//
//	stdoutReader, stdoutWriter, _ := os.Pipe()
//	stderrReader, stderrWriter, _ := os.Pipe()
//
//	oldStdout := os.Stdout
//	oldStderr := os.Stderr
//	os.Stdout = stdoutWriter
//	os.Stderr = stderrWriter
//
//	done := make(chan bool, 2)
//	readAndCapture := func(reader *os.File, buffer *bytes.Buffer) {
//		buf := make([]byte, 1024)
//		for {
//			n, err := reader.Read(buf)
//			if err != nil {
//				if err != io.EOF {
//					// log.Println("read error:", err)
//				}
//				break
//			}
//
//			if _, err := buffer.Write(buf[:n]); err != nil {
//				// log.Println("buffer write:", err)
//				return
//			}
//			m := marketplace.LogLine{
//				Line:      string(buf[:n]),
//				Level:     "",
//				Timestamp: "",
//				Type:      "",
//			}
//			mes := marketplace.UnityWebsocketMessage_Logs{Logs: &m}
//			mess := &marketplace.UnityWebsocketMessage{Content: &mes}
//			data, err := proto.Marshal(mess)
//			if err != nil {
//				log.WithError(err).Error("Failed to marshal log line")
//				return
//			}
//			wsmgr.SendMessageToUserID(id, data)
//		}
//		done <- true
//	}
//
//	go readAndCapture(stdoutReader, &stdoutBuffer)
//	go readAndCapture(stderrReader, &stderrBuffer)
//
//	return func() {
//		stdoutWriter.Close()
//		stderrWriter.Close()
//		os.Stdout = oldStdout
//		os.Stderr = oldStderr
//		<-done
//		<-done
//	}
//}
