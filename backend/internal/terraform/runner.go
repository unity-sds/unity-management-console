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
	"time"
)

var stdoutBuffer bytes.Buffer
var stderrBuffer bytes.Buffer

type TerraformExecutor interface {
	NewTerraform(string, string) (*tfexec.Terraform, error)
	Init(context.Context, ...tfexec.InitOption) error
	Plan(context.Context, ...tfexec.PlanOption) (bool, error)
	Apply(context.Context, ...tfexec.ApplyOption) error
	Destroy(context.Context, ...tfexec.DestroyOption) error
	SetStdout(io.Writer)
	SetStderr(io.Writer)
	SetLogger(*log.Logger)
	WorkspaceNew(context.Context, string) error
	WorkspaceSelect(context.Context, string) error
	WorkspaceDelete(context.Context, string) error
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

func (r *RealTerraformExecutor) Destroy(ctx context.Context, opts ...tfexec.DestroyOption) error {
	return r.tf.Destroy(ctx, opts...)
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

func (r *RealTerraformExecutor) WorkspaceNew(ctx context.Context, name string) error {
	return r.tf.WorkspaceNew(ctx, name)
}

func (r *RealTerraformExecutor) WorkspaceSelect(ctx context.Context, name string) error {
	return r.tf.WorkspaceSelect(ctx, name)
}

func (r *RealTerraformExecutor) WorkspaceDelete(ctx context.Context, name string) error {
	return r.tf.WorkspaceDelete(ctx, name)
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
	data, err := proto.Marshal(&marketplace.UnityWebsocketMessage{Content: &mes})
	if err != nil {
		log.WithError(err).Error("Failed to marshal log line")
		return
	}
	w.wsmgr.SendMessageToUserID(w.userid, data)

	return n, err
}

func RunTerraform(appconf *config.AppConfig, wsmgr *ws.WebSocketManager, id string, executor TerraformExecutor, target string) error {
	bucket := fmt.Sprintf("bucket=%s", appconf.BucketName)

	key := fmt.Sprintf("key=%s-%s-tfstate", appconf.Project, appconf.Venue)
	region := fmt.Sprintf("region=%s", appconf.AWSRegion)

	p := filepath.Join(appconf.Workdir, "workspace")
	tf, err := executor.NewTerraform(p, "/usr/local/bin/terraform")
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
	writerStderr := io.MultiWriter(os.Stderr, os.Stdout, wwserrWriter)

	tf.SetStdout(writerStdout)
	tf.SetStderr(writerStderr)
	tf.SetLogger(log.StandardLogger())

	err = executor.Init(context.Background(), tfexec.Upgrade(true), tfexec.BackendConfig(bucket), tfexec.BackendConfig(key), tfexec.BackendConfig(region))

	if err != nil {
		log.WithError(err).Error("error initialising terraform")
		message := marketplace.SimpleMessage{
			Operation: "terraform",
			Payload:   "failed",
		}
		om := marketplace.UnityWebsocketMessage_Simplemessage{Simplemessage: &message}
		wrap := marketplace.UnityWebsocketMessage{Content: &om}
		b, _ := proto.Marshal(&wrap)
		wsmgr.SendMessageToUserID(id, b)
		return err
	}
	writerStdout.Write([]byte("Waiting 30 seconds for the state to settle"))
	time.Sleep(30 * time.Second)
	change := false
	if target != "" {
		log.Infof("Running terraform with target: %s", target)
		change, err = executor.Plan(context.Background(), tfexec.Target(target))
	} else {
		change, err = executor.Plan(context.Background())
	}

	if err != nil {
		log.WithError(err).Error("error running plan")
		message := marketplace.SimpleMessage{
			Operation: "terraform",
			Payload:   "failed",
		}
		om := marketplace.UnityWebsocketMessage_Simplemessage{Simplemessage: &message}
		wrap := marketplace.UnityWebsocketMessage{Content: &om}
		b, _ := proto.Marshal(&wrap)
		wsmgr.SendMessageToUserID(id, b)
		return err
	}

	if change {
		if target != "" {
			log.Infof("Running terraform with target: %s", target)
			err = executor.Apply(context.Background(), tfexec.Target(target))
		} else {
			err = executor.Apply(context.Background())
		}

		if err != nil {
			message := marketplace.SimpleMessage{
				Operation: "terraform",
				Payload:   "failed",
			}
			om := marketplace.UnityWebsocketMessage_Simplemessage{Simplemessage: &message}
			wrap := marketplace.UnityWebsocketMessage{Content: &om}
			b, _ := proto.Marshal(&wrap)
			wsmgr.SendMessageToUserID(id, b)
			log.WithError(err).Error("error running apply")
			return err
		}

	}
	message := marketplace.SimpleMessage{
		Operation: "terraform",
		Payload:   "completed",
	}
	om := marketplace.UnityWebsocketMessage_Simplemessage{Simplemessage: &message}
	wrap := marketplace.UnityWebsocketMessage{Content: &om}
	b, _ := proto.Marshal(&wrap)
	wsmgr.SendMessageToUserID(id, b)
	return nil
}

func RunTerraformLogOutToFile(appconf *config.AppConfig, logfile string, executor TerraformExecutor, target string) error {
	bucket := fmt.Sprintf("bucket=%s", appconf.BucketName)
	key := fmt.Sprintf("key=%s-%s-tfstate", appconf.Project, appconf.Venue)
	region := fmt.Sprintf("region=%s", appconf.AWSRegion)

	p := filepath.Join(appconf.Workdir, "workspace")
	tf, err := executor.NewTerraform(p, "/usr/local/bin/terraform")
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	// Open the log file in append mode
	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening log file: %s", err)
	}
	defer file.Close()

	// Create a multi-writer that writes to both file and stdout
	writerStdout := io.MultiWriter(os.Stdout, file)
	writerStderr := io.MultiWriter(os.Stderr, file)

	tf.SetStdout(writerStdout)
	tf.SetStderr(writerStderr)
	tf.SetLogger(log.StandardLogger())

	// Create a workspace for this installation
	workspaceName := fmt.Sprintf("%s-%s", target, time.Now().Format("20060102150405"))
	err = tf.WorkspaceNew(context.Background(), workspaceName)
	if err != nil {
		log.WithError(err).Error("error creating workspace")
		file.WriteString(fmt.Sprintf("Error creating workspace: %s\n", err))
		return err
	}

	file.WriteString("Starting Terraform")

	err = executor.Init(context.Background(), tfexec.Upgrade(true), tfexec.BackendConfig(bucket), tfexec.BackendConfig(key), tfexec.BackendConfig(region))

	if err != nil {
		log.WithError(err).Error("error initialising terraform")
		file.WriteString(fmt.Sprintf("Error initialising terraform: %s\n", err))
		return err
	}

	file.WriteString("Waiting 60 seconds for the state to settle\n")
	time.Sleep(60 * time.Second)

	change := false
	if target != "" {
		log.Infof("Running terraform with target: %s", target)
		file.WriteString(fmt.Sprintf("Running terraform with target: %s\n", target))
		change, err = executor.Plan(context.Background(), tfexec.Target(target))
	} else {
		change, err = executor.Plan(context.Background())
	}

	if err != nil {
		log.WithError(err).Error("error running plan")
		file.WriteString(fmt.Sprintf("Error running plan: %s\n", err))
		return err
	}

	if change {
		if target != "" {
			log.Infof("Running terraform with target: %s", target)
			file.WriteString(fmt.Sprintf("Running terraform with target: %s\n", target))
			err = executor.Apply(context.Background(), tfexec.Target(target))
		} else {
			err = executor.Apply(context.Background())
		}

		if err != nil {
			log.WithError(err).Error("error running apply")
			file.WriteString(fmt.Sprintf("Error running apply: %s\n", err))
			return err
		}
	}

	file.WriteString("Terraform operation completed successfully\n")
	return nil
}

func DestroyTerraformModule(appconf *config.AppConfig, logfile string, executor TerraformExecutor, moduleName string) error {
	bucket := fmt.Sprintf("bucket=%s", appconf.BucketName)
	key := fmt.Sprintf("key=%s-%s-tfstate", appconf.Project, appconf.Venue)
	region := fmt.Sprintf("region=%s", appconf.AWSRegion)

	p := filepath.Join(appconf.Workdir, "workspace")
	tf, err := executor.NewTerraform(p, "/usr/local/bin/terraform")
	if err != nil {
		log.Fatalf("error running Terraform: %s", err)
	}

	// Open the log file in append mode
	logfileHandle, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening log file: %s", err)
	}
	defer logfileHandle.Close()

	// Create a multi-writer that writes to both file and stdout
	writerStdout := io.MultiWriter(os.Stdout, logfileHandle)
	writerStderr := io.MultiWriter(os.Stderr, logfileHandle)

	tf.SetStdout(writerStdout)
	tf.SetStderr(writerStderr)
	tf.SetLogger(log.StandardLogger())

	// Switch to the workspace for this module
	workspaceName := fmt.Sprintf("%s-%s", moduleName, time.Now().Format("20060102150405"))
	err = tf.WorkspaceSelect(context.Background(), workspaceName)
	if err != nil {
		log.WithError(err).Error("error selecting workspace")
		logfileHandle.WriteString(fmt.Sprintf("Error selecting workspace: %s\n", err))
		return err
	}

	logfileHandle.WriteString("Starting Terraform")

	err = executor.Init(context.Background(), tfexec.Upgrade(true), tfexec.BackendConfig(bucket), tfexec.BackendConfig(key), tfexec.BackendConfig(region))

	if err != nil {
		log.WithError(err).Error("error initialising terraform")
		logfileHandle.WriteString(fmt.Sprintf("Error initialising terraform: %s\n", err))
		return err
	}

	logfileHandle.WriteString("Waiting 60 seconds for the state to settle\n")
	time.Sleep(60 * time.Second)

	err = executor.Destroy(context.Background(), tfexec.Target(fmt.Sprintf("module.%s", moduleName)))
	if err != nil {
		log.WithError(err).Error("error running terraform destroy")
		logfileHandle.WriteString(fmt.Sprintf("Error running destroy: %s\n", err))
		return err
	}

	// Delete the workspace after successful destroy
	err = tf.WorkspaceDelete(context.Background(), workspaceName)
	if err != nil {
		log.WithError(err).Error("error deleting workspace")
		logfileHandle.WriteString(fmt.Sprintf("Error deleting workspace: %s\n", err))
		// Don't return error here as the main operation (destroy) was successful
	}

	logfileHandle.WriteString("Terraform destroy completed successfully\n")
	return nil
}

func DestroyAllTerraform(appconf *config.AppConfig, wsmgr *ws.WebSocketManager, id string, executor TerraformExecutor) error {
	p := filepath.Join(appconf.Workdir, "workspace")

	tf, err := executor.NewTerraform(p, "/usr/local/bin/terraform")
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	writerStdout := io.MultiWriter(os.Stdout)
	writerStderr := io.MultiWriter(os.Stderr)
	if wsmgr != nil {
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
		writerStdout = io.MultiWriter(os.Stdout, wwsWriter)
		writerStderr = io.MultiWriter(os.Stderr, wwserrWriter)
	}

	tf.SetStdout(writerStdout)
	tf.SetStderr(writerStderr)
	tf.SetLogger(log.StandardLogger())

	err = tf.Destroy(context.Background())
	if err != nil {
		log.WithError(err).Error("error running terraform destroy")
		if wsmgr != nil {
			message := marketplace.SimpleMessage{
				Operation: "terraform",
				Payload:   "failed",
			}
			om := marketplace.UnityWebsocketMessage_Simplemessage{Simplemessage: &message}
			wrap := marketplace.UnityWebsocketMessage{Content: &om}
			b, _ := proto.Marshal(&wrap)
			wsmgr.SendMessageToUserID(id, b)
			return err
		}
	}

	if wsmgr != nil {
		message := marketplace.SimpleMessage{
			Operation: "terraform",
			Payload:   "completed",
		}
		om := marketplace.UnityWebsocketMessage_Simplemessage{Simplemessage: &message}
		wrap := marketplace.UnityWebsocketMessage{Content: &om}
		b, _ := proto.Marshal(&wrap)
		wsmgr.SendMessageToUserID(id, b)
	}
	return nil
}
