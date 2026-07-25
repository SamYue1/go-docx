package features

import (
	"os"
	"testing"

	"github.com/SamYue1/go-docx/test/features/steps"
	"github.com/cucumber/godog"
)

func TestMain(m *testing.M) {
	opts := godog.Options{
		Format: "progress",
		Paths:  []string{"."},
	}
	status := godog.TestSuite{
		Name: "go-docx",
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			steps.RegisterSteps(ctx)
		},
		Options: &opts,
	}.Run()
	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
