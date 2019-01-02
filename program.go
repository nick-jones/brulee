package brulee

import (
	"github.com/olekukonko/tablewriter"
	"io"
	"strconv"
)

type Scores map[string]int

type Program struct {
	ins []Instruction
}

func (p Program) Run(vars map[string]string) (Scores, error) {
	i := &Interpreter{
		ins:     p.ins,
		vars:    vars,
		scratch: make(map[ScratchPosition]bool),
		scores:  make(Scores),
	}
	i.Interpret()
	if err := i.Err(); err != nil {
		return Scores{}, err
	}
	return i.Scores(), nil
}

func (p Program) Dump(w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Pos", "Op", "Ret", "Operand1", "Operand2"})

	for i, in := range p.ins {
		s := []string{strconv.Itoa(i)}
		s = append(s, in.StringSlice()...)
		table.Append(s)
	}

	table.Render()
}