package act

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/nektos/act/pkg/common"
	"github.com/nektos/act/pkg/model"
	"github.com/nektos/act/pkg/runner"
	"github.com/sirupsen/logrus"
)

type Config struct {
}
type Runner struct {
	name string

	cfg *Config

	envs map[string]string

	runningTasks sync.Map
}

// Define a struct
type MyLogger struct {
	Output bytes.Buffer
}

// Implement the JobLoggerFactory interface on the struct
func (m *MyLogger) WithJobLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Info("MyLogger: WithJobLogger() invoked")
	logger.SetOutput(&m.Output)
	return logger
}
func RunAct(conn *websocket.Conn) {
	var plan *model.Plan
	planner, err := model.NewWorkflowPlanner("/home/barber/Projects/unity-cs-infra/.github/workflows/deploy_eks.yml", false)

	//plan, plannerErr := planner.PlanJob("deploy_eks")
	plan, plannerErr := planner.PlanEvent("workflow_dispatch")
	if plan == nil && plannerErr != nil {
		fmt.Printf("%v", plannerErr)
	}
	runnerConfig := &runner.Config{
		Workdir:     "/home/barber/Projects/unity-cs-infra",
		BindWorkdir: false,

		ReuseContainers:  false,
		ForcePull:        false,
		ForceRebuild:     false,
		LogOutput:        true,
		JSONLogger:       false,
		Env:              map[string]string{},
		Secrets:          map[string]string{},
		GitHubInstance:   "github.com",
		AutoRemove:       true,
		NoSkipCheckout:   true,
		ContainerOptions: "",
		Privileged:       false,
		Platforms:        map[string]string{"ubuntu-latest": "catthehacker/ubuntu:act-latest"},
	}
	myLogger := &MyLogger{}
	var loggerFactory runner.JobLoggerFactory
	loggerFactory = myLogger
	runner.WithJobLoggerFactory(context.Background(), loggerFactory)
	var logBuffer bytes.Buffer

	// Create a new logger
	logger := logrus.New()

	// Set the logger output to the bytes.Buffer
	logger.SetOutput(&logBuffer)
	common.WithLogger(context.Background(), logger)
	rr, err := runner.New(runnerConfig)

	if err != nil {
		//	fmt.Printf("err %v", err)
	} else {
		//	fmt.Printf("%v", rr)
		fmt.Println("DONE")
		fmt.Printf("%v", loggerFactory.WithJobLogger().Out)
	}

	//	executor := rr.NewPlanExecutor(plan).Finally(func(ctx context.Context) error {
	//fmt.Printf("%v", "here")
	//	return nil
	//	})
	//err = executor(context.Background())
	//if err != nil {
	//fmt.Printf("%v", err)
	//	}

	// Create a pipe for stdout and stderr
	stdoutReader, stdoutWriter, _ := os.Pipe()
	stderrReader, stderrWriter, _ := os.Pipe()

	// Replace os.Stdout and os.Stderr with the writer end of the pipe
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	os.Stdout = stdoutWriter
	os.Stderr = stderrWriter

	var wg sync.WaitGroup
	wg.Add(3)

	// Create buffers to hold the output
	var stdoutBuffer, stderrBuffer bytes.Buffer

	// Start a goroutine to read from the pipe and write to the buffer
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

			if err := conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
				//log.Println("write:", err)
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		io.Copy(&stderrBuffer, stderrReader)
	}()

	// Now you can run your code
	go func() {
		defer wg.Done()

		// Your code goes here
		executor := rr.NewPlanExecutor(plan).Finally(func(ctx context.Context) error {
			//fmt.Printf("%v", "here")
			return nil
		})
		err := executor(context.Background())
		if err != nil {
			//fmt.Printf("%v", err)
		}

		// Once done, close the writers and restore the original stdout and stderr
		stdoutWriter.Close()
		stderrWriter.Close()
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	// Now you can print the captured output
	fmt.Println("stdout:", stdoutBuffer.String())
	fmt.Println("stderr:", stderrBuffer.String())
	fmt.Println("over done")
	fmt.Printf("%v", myLogger.Output.String())
	fmt.Printf("%v", logBuffer.String())

}
