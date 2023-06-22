package act

import (
	"bytes"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/unity-sds/unity-control-plane/backend/internal/application/config"

	"github.com/gorilla/websocket"
	"github.com/nektos/act/pkg/common"
	"github.com/nektos/act/pkg/model"
	"github.com/nektos/act/pkg/runner"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type ActRunner struct {
	Workflow      string
	Inputs        map[string]string
	Env           map[string]string
	Secrets       map[string]string
	Conn          *websocket.Conn
	Workdir       string
	RunnerConfig  *runner.Config
	Plan          *model.Plan
	LoggerFactory runner.JobLoggerFactory
	Logger        *logrus.Logger
	StdoutBuffer  bytes.Buffer
	StderrBuffer  bytes.Buffer
	PlanExecutor  common.Executor
	AppConfig     config.AppConfig
}

type MyLogger struct {
	Output bytes.Buffer
}

// WithJobLogger creates a new logrus Logger, sets an information level log
// indicating the invocation of this method, and directs its output to the
// buffer of the MyLogger instance. It then returns the newly created logger.
// This method can be useful when you want to capture the logs in a buffer
// for further processing or testing.
//
// Example usage:
//
//	myLogger := &MyLogger{}
//	logger := myLogger.WithJobLogger()
//	logger.Info("This log will be stored in myLogger's buffer.")
//
// Returns:
//
//	*logrus.Logger : a pointer to the newly created logrus Logger instance
func (m *MyLogger) WithJobLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Info("MyLogger: WithJobLogger() invoked")
	logger.SetOutput(&m.Output)
	return logger
}
func NewActRunner(workflow string, inputs, env, secrets map[string]string, conn *websocket.Conn, appConfig config.AppConfig) *ActRunner {
	// setup the default ActRunner here
	return &ActRunner{Workflow: workflow, Inputs: inputs, Env: env, Secrets: secrets, Conn: conn, AppConfig: appConfig, Workdir: appConfig.Workdir}
}

// CreateWorkflowPlan is a method of ActRunner that creates a workflow plan based
// on the provided workflow in the ActRunner instance. It uses the model.NewWorkflowPlanner
// function to initialize a workflow planner with the given workflow and a non-forked flag.
// The planner is then used to generate a plan for the "workflow_dispatch" event.
// The resulting plan is stored in the ActRunner's Plan field.
//
// This method will return an error if any errors are encountered during the
// planning process. The errors may come from the model.NewWorkflowPlanner function or
// from the planner.PlanEvent method, if they are unable to create the plan as expected.
//
// Returns:
//
//	error : an error object if an error occurs during the planning process, otherwise nil
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

// CreateRunnerConfig is a method of ActRunner that creates a runner configuration
// for the ActRunner instance. It sets various configuration options such as working
// directory, GitHub token, environment variables, secrets, inputs and others. These
// settings are based on the current ActRunner instance's values and constants.
// The runner configuration is stored in the ActRunner's RunnerConfig field.
//
// Note: This function currently always returns nil as it does not contain any operations
// that could cause an error. Future enhancements might change this behavior.
//
// Returns:
//
//	error : an error object if an error occurs during the configuration setup, otherwise nil
func (ar *ActRunner) CreateRunnerConfig() error {
	ar.RunnerConfig = &runner.Config{
		Workdir:          ar.Workdir,
		EventName:        "workflow_dispatch",
		BindWorkdir:      false,
		Token:            ar.AppConfig.GithubToken,
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

// SetupLogger is a method of ActRunner that sets up a logger for the ActRunner instance.
// It initializes the LoggerFactory field with a new instance of MyLogger and sets
// it as the job logger factory for the runner package.
// A new instance of logrus.Logger is created and its output is set to the
// ActRunner's StdoutBuffer to capture the log output. This Logger instance is set as
// the logger for the common package. The method does not return any value.
func (ar *ActRunner) SetupLogger() {
	ar.LoggerFactory = &MyLogger{}
	runner.WithJobLoggerFactory(context.Background(), ar.LoggerFactory)
	ar.Logger = logrus.New()
	ar.Logger.SetOutput(&ar.StdoutBuffer)
	common.WithLogger(context.Background(), ar.Logger)
}

// RunWorkflow is a method of ActRunner that executes the workflow using the ActRunner's RunnerConfig.
// It initializes a new runner instance with the RunnerConfig, then creates a new plan executor from the
// runner instance and the ActRunner's Plan. The Finally method of the plan executor is set up to ensure that
// the executor is executed even if an error occurs. The workflow is executed by calling the PlanExecutor with the
// background context. The method returns an error if there is an error during the execution.
func (ar *ActRunner) RunWorkflow() error {
	rr, err := runner.New(ar.RunnerConfig)
	if err != nil {
		return err
	}
	ar.PlanExecutor = rr.NewPlanExecutor(ar.Plan).Finally(func(ctx context.Context) error {
		return nil
	})
	log.Infof("Final token : %v", ar.RunnerConfig.Token)
	return ar.PlanExecutor(context.Background())
}

// CaptureOutput is a method of ActRunner that captures the standard output and standard error streams during
// the execution of the workflow. It creates pipe readers and writers for the standard output and standard error,
// then replaces the original standard output and error streams with the writers to redirect the output. The method
// launches three goroutines to read from the pipe readers and write the output to the corresponding buffers. If a
// websocket connection is available, the captured output is also sent through the connection. After the workflow
// execution, the original standard output and error streams are restored, and the pipe writers are closed. This
// method blocks until all goroutines are finished.
func (ar *ActRunner) CaptureOutput() (stopCapture func()) {
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
			if ar.Conn != nil {
				if err := ar.Conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
					// log.Println("write:", err)
					return
				}
			}
		}
		done <- true
	}

	go readAndCapture(stdoutReader, &ar.StdoutBuffer)
	go readAndCapture(stderrReader, &ar.StderrBuffer)

	return func() {
		stdoutWriter.Close()
		stderrWriter.Close()
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		<-done
		<-done
	}
}

// PrintOutput is a method of ActRunner that prints the captured standard output and standard error from the
// workflow execution. It prints the contents of the ActRunner's StdoutBuffer and StderrBuffer to the console
// using the fmt.Println function.
func (ar *ActRunner) PrintOutput() {
	fmt.Println("stdout:", ar.StdoutBuffer.String())
	fmt.Println("stderr:", ar.StderrBuffer.String())
}

// RunAct is a function that executes an Act workflow with the provided parameters. It first creates a new
// ActRunner instance using the provided workflow, inputs, environment variables, secrets, and WebSocket connection.
// Then it creates the workflow plan, runner configuration, and sets up the logger for the ActRunner instance.
// After this setup, it runs the workflow, captures the output and prints the captured output.
// If an error occurs at any point during this process, it is returned.
func RunAct(workflow string, inputs map[string]string, env map[string]string, secrets map[string]string, conn *websocket.Conn, appConfig config.AppConfig) error {
	log.Info("Creating runner")
	ar := NewActRunner(workflow, inputs, env, secrets, conn, appConfig)

	log.Info("Creating plan")
	if err := ar.CreateWorkflowPlan(); err != nil {
		return err
	}
	log.Info("Create config")
	err := ar.CreateRunnerConfig()
	if err != nil {
		return err
	}

	log.Info("Configuring logger")
	ar.SetupLogger()

	stopCapture := ar.CaptureOutput()
	defer stopCapture() // Ensure that we stop capturing even if an error occurs
	log.Info("Running workflow")

	if err := ar.RunWorkflow(); err != nil {
		return err
	}
	ar.PrintOutput()

	return nil
}
