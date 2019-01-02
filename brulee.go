package brulee

import (
	"io"
	"strconv"

	"github.com/nick-jones/brulee/internal"
	"github.com/olekukonko/tablewriter"
)

func Compile(r io.Reader) (Program, error) {
	program := Program{}
	root, err := internal.Parse(r)
	if err != nil {
		return program, err
	}
	ig := internal.NewInstructionsGenerator()
	ig.Generate(root)
	if err := ig.Err(); err != nil {
		return program, err
	}
	program.load(ig.Instructions())
	return program, nil
}

func MustCompile(r io.Reader) Program {
	program, err := Compile(r)
	if err != nil {
		panic(err)
	}
	return program
}

type Program struct {
	ins []internal.Instruction
}

func (p Program) Run(vars map[string]string) (map[string]int, error) {
	i := internal.NewExecutor(p.ins, vars)
	i.Execute()
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

func (p *Program) load(ins []internal.Instruction) {
	p.ins = ins
}
