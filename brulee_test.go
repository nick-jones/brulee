package brulee

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

var (
	program Program
	vars    map[string]string
	scores  map[string]int
)

func theProgram(p *gherkin.DocString) error {
	var err error
	program, err = Compile(strings.NewReader(p.Content))
	return err
}

func variables(table *gherkin.DataTable) error {
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

func theScoreOutputIs(table *gherkin.DataTable) error {
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

func FeatureContext(s *godog.Suite) {
	s.BeforeScenario(func(i interface{}) {
		program = Program{}
		vars = map[string]string{}
		scores = map[string]int{}
	})
	s.Step(`^the program:$`, theProgram)
	s.Step(`^variables:$`, variables)
	s.Step(`^the program is run$`, theProgramIsRun)
	s.Step(`^the score output is:$`, theScoreOutputIs)
	s.Step(`^the score output is empty$`, theScoreOutputIsEmpty)
}
