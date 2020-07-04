package brulee

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
)

var (
	program Program
	vars    map[string]string
	scores  map[string]int
)

func theProgram(p *messages.PickleStepArgument_PickleDocString) error {
	var err error
	program, err = Compile(strings.NewReader(p.Content))
	return err
}

func variables(table *messages.PickleStepArgument_PickleTable) error {
	for _, row := range table.Rows[1:] {
		vars[row.Cells[0].Value] = row.Cells[1].Value
	}
	return nil
}

func theProgramIsRun() error {
	var err error
	scores, err = program.Run(vars)
	return err
}

func theScoreOutputIs(table *messages.PickleStepArgument_PickleTable) error {
	if len(table.Rows)-1 != len(scores) {
		return fmt.Errorf("row count mismatch, expected %d, actual %d", len(table.Rows)-1, len(scores))
	}
	for _, row := range table.Rows[1:] {
		name := row.Cells[0].Value
		value, err := strconv.Atoi(row.Cells[1].Value)
		if err != nil {
			return err
		}
		if scores[name] != value {
			return fmt.Errorf("score mismatch for %s, expected %d, actual %d", name, value, scores[name])
		}
	}
	return nil
}

func theScoreOutputIsEmpty() error {
	if len(scores) != 0 {
		return fmt.Errorf("score output expected to be empty, actual %d", len(scores))
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(_ *godog.Scenario) {
		program = Program{}
		vars = map[string]string{}
		scores = map[string]int{}
	})
	ctx.Step(`^the program:$`, theProgram)
	ctx.Step(`^variables:$`, variables)
	ctx.Step(`^the program is run$`, theProgramIsRun)
	ctx.Step(`^the score output is:$`, theScoreOutputIs)
	ctx.Step(`^the score output is empty$`, theScoreOutputIsEmpty)
}

func TestMain(m *testing.M) {
	opts := godog.Options{
		Format:    "progress",
		Paths:     []string{"features"},
		Randomize: time.Now().UTC().UnixNano(),
	}

	status := godog.TestSuite{
		Name: "brulee",
		ScenarioInitializer:  InitializeScenario,
		Options: &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
