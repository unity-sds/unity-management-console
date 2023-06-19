package act

import (
	"github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRunner(t *testing.T) {
	//RunAct("/home/barber/Projects/unity-cs-infra/.github/workflows/deploy_eks.yml", nil)
}

func TestActRunner(t *testing.T) {
	Convey("Given a new ActRunner instance", t, func() {
		ar := NewActRunner("../../../test-data/workflows/demo-workflow.yml", map[string]string{}, map[string]string{}, map[string]string{}, &websocket.Conn{})

		Convey("When the instance is created", func() {
			Convey("Then the instance should not be nil", func() {
				So(ar, ShouldNotBeNil)
			})
		})

		Convey("When CreateWorkflowPlan is called", func(){
			err := ar.CreateWorkflowPlan()

			Convey("Then no error should be returned", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When CreateRunnerConfig is called", func() {
			err := ar.CreateRunnerConfig()

			Convey("Then no error should be returned", func(){
				So(err, ShouldBeNil)
			})
		})

		Convey("When SetupLogger is called", func() {
			ar.SetupLogger()

			Convey("Then LoggerFactory should not be nil", func() {
				So(ar.LoggerFactory, ShouldNotBeNil)
			})
		})
	})
}