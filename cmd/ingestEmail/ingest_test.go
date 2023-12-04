package ingestemail

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/cucumber/godog"
)

const outFile = "jsonl/eml.jsonl"

func arrange() error {
	if _, err := os.Stat("fixtures/test1.eml"); errors.Is(err, os.ErrNotExist) {
		err := os.Link("fixtures/test1.eml", "input/test1.eml")
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat("fixtures/test2.eml"); errors.Is(err, os.ErrNotExist) {
		err = os.Link("fixtures/test2.eml", "input/test2.eml")
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat(outFile); err == nil {
		os.Remove(outFile)
	}
	return nil
}

func act() error {
	err := Ingest("input", outFile)
	return err
}

func assert() error {
	expected, err := os.ReadFile("fixtures/expectedEmlIngest.jsonl")
	if err != nil {
		return err
	}
	actual, err := os.ReadFile("jsonl/eml.jsonl")
	if err != nil {
		return err
	}

	if string(actual) != string(expected) {
		return fmt.Errorf("Expected\n%s\n but got\n%s", string(actual), string(expected))
	}
	return nil
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Given("Test emails exist in the input folder", arrange)
	ctx.When("I run the ingest function", act)
	ctx.Then("I should get email metadata in the output folder", assert)
}
