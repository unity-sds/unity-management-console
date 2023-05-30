package act

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/nektos/act/pkg/common"
	"github.com/nektos/act/pkg/model"
	"github.com/nektos/act/pkg/runner"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

type ActRunner struct {
	Workflow       string
	Inputs         map[string]string
	Env            map[string]string
	Secrets        map[string]string
	Conn           *websocket.Conn
	Workdir        string
	RunnerConfig   *runner.Config
	Plan           *model.Plan
	LoggerFactory  runner.JobLoggerFactory
	Logger         *logrus.Logger
	StdoutBuffer   bytes.Buffer
	StderrBuffer   bytes.Buffer
	PlanExecutor   common.Executor
}

type MyLogger struct {
	Output bytes.Buffer
}
func (m *MyLogger) WithJobLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Info("MyLogger: WithJobLogger() invoked")
	logger.SetOutput(&m.Output)
	return logger
}
func NewActRunner(workflow string, inputs, env, secrets map[string]string, conn *websocket.Conn) *ActRunner {
	// setup the default ActRunner here
	return &ActRunner{Workflow: workflow, Inputs: inputs, Env: env, Secrets: secrets, Conn: conn}
}

func (ar *ActRunner) CreateWorkflowPlan() error {
	planner, err := model.NewWorkflowPlanner(ar.Workflow, false)
	if err != nil {
		return err
	}
	ar.Plan, err = planner.PlanEvent("workflow_dispatch")
	if err != nil {
		return err
	}
	return nil
}

func (ar *ActRunner) CreateRunnerConfig() error {
	ar.RunnerConfig = &runner.Config{
		Workdir:          ar.Workdir,
		BindWorkdir:      false,
		Token:            os.Getenv("GITHUB_TOKEN"),
		ReuseContainers:  false,
		ForcePull:        false,
		LogOutput:        true,
		JSONLogger:       false,
		Env:              ar.Env,
		Secrets:          ar.Secrets,
		Inputs:           ar.Inputs,
		GitHubInstance:   "github.com",
		AutoRemove:       true,
		NoSkipCheckout:   true,
		ContainerOptions: "",
		Privileged:       false,
		Platforms:        map[string]string{"ubuntu-latest": "catthehacker/ubuntu:act-latest"},
	}
	return nil
}

func (ar *ActRunner) SetupLogger() {
	ar.LoggerFactory = &MyLogger{}
	runner.WithJobLoggerFactory(context.Background(), ar.LoggerFactory)
	ar.Logger = logrus.New()
	ar.Logger.SetOutput(&ar.StdoutBuffer)
	common.WithLogger(context.Background(), ar.Logger)
}

func (ar *ActRunner) RunWorkflow() error {
	rr, err := runner.New(ar.RunnerConfig)
	if err != nil {
		return err
	}
	ar.PlanExecutor = rr.NewPlanExecutor(ar.Plan).Finally(func(ctx context.Context) error {
		return nil
	})
	return ar.PlanExecutor(context.Background())
}

func (ar *ActRunner) CaptureOutput() {
	stdoutReader, stdoutWriter, _ := os.Pipe()
	stderrReader, stderrWriter, _ := os.Pipe()

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	os.Stdout = stdoutWriter
	os.Stderr = stderrWriter

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		buf := make([]byte, 1024)
		for {
			n, err := stdoutReader.Read(buf)
			if err != nil {
				if err != io.EOF {
					//log.Println("read error:", err)
				}
				break
			}

			if _, err := ar.StdoutBuffer.Write(buf[:n]); err != nil {
				//log.Println("buffer write:", err)
				return
			}
			if ar.Conn != nil {
				if err := ar.Conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
					//log.Println("write:", err)
					return
				}
			}
		}
	}()
	go func() {
		defer wg.Done()
		io.Copy(&ar.StderrBuffer, stderrReader)
	}()
	go func() {
		defer wg.Done()
		stdoutWriter.Close()
		stderrWriter.Close()
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()
	wg.Wait()
}

func (ar *ActRunner) PrintOutput() {
	fmt.Println("stdout:", ar.StdoutBuffer.String())
	fmt.Println("stderr:", ar.StderrBuffer.String())
}

func RunAct(workflow string, inputs map[string]string, env map[string]string, secrets map[string]string, conn *websocket.Conn) error {
	ar := NewActRunner(workflow, inputs, env, secrets, conn)

	if err := ar.CreateWorkflowPlan(); err != nil {
		return err
	}
	err := ar.CreateRunnerConfig()
	if err != nil {
		return err
	}
	ar.SetupLogger()
	if err := ar.RunWorkflow(); err != nil {
		return err
	}
	ar.CaptureOutput()
	ar.PrintOutput()

	return nil
}
