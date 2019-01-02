package brulee

import (
	"github.com/nick-jones/brulee/internal"
	"github.com/olekukonko/tablewriter"
	"io"
	"strconv"
)

func Compile(r io.Reader) (Program, error) {
	root, err := internal.Parse(r)
	if err != nil {
		return Program{}, err
	}
	comp := internal.NewCompiler()
	comp.Compile(root)
	program := Program{ins: comp.Instructions()}
	return program, nil
}

type Program struct {
	ins []internal.Instruction
}

func (p Program) Run(vars map[string]string) (map[string]int, error) {
	i := internal.NewInterpreter(p.ins, vars)
	i.Run()
	if err := i.Err(); err != nil {
		return nil, err
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