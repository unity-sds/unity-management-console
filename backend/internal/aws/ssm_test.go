package aws

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestReadSSMParameter(t *testing.T) {
	Convey("Given an AWS account and an AWS connection", t, func() {
		Convey("When ReadSSMParameter is called", func() {

			Convey("Then a parameter and value should be returned", func() {

			})
		})
	})
}

func TestReadMissingSSMParameter(t *testing.T) {
	Convey("Given an AWS account and an AWS connection", t, func() {
		Convey("When ReadSSMParameter is called but the parameter doesn't exist", func() {

			Convey("Then the error should be handled", func() {

			})
		})
	})
}

func TestReadSSMParameters(t *testing.T) {
	Convey("Given an AWS account and an AWS connection", t, func() {
		Convey("When ReadSSMParameters is called", func() {

			Convey("Then a list of SSM Parameters should be returned", func() {

			})
		})
	})
}
