package act

import (
	"context"
	"fmt"
	"sync"

	"github.com/nektos/act/pkg/model"
	"github.com/nektos/act/pkg/runner"
)

type Config struct {
}
type Runner struct {
	name string

	cfg *Config

	envs map[string]string

	runningTasks sync.Map
}

func runAct() {
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

	rr, err := runner.New(runnerConfig)
	if err != nil {
		fmt.Printf("err %v", err)
	} else {
		fmt.Printf("%v", rr)
	}

	executor := rr.NewPlanExecutor(plan).Finally(func(ctx context.Context) error {
		fmt.Printf("%v", "here")
		return nil
	})
	err = executor(context.TODO())
	if err != nil {
		fmt.Printf("%v", err)
	}
}
