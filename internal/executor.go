package internal

import (
	"fmt"
	"regexp"
	"strings"
)

type Executor struct {
	ins     []Instruction
	vars    map[string]string
	scratch map[ScratchPosition]bool
	scores  map[string]int
	err     error
}

func NewExecutor(ins []Instruction, vars map[string]string) *Executor {
	return &Executor{
		ins:     ins,
		vars:    vars,
		scratch: map[ScratchPosition]bool{},
		scores:  map[string]int{},
	}
}

// nolint:gocyclo
func (i *Executor) Execute() {
	pos := 0
Loop:
	for pos < len(i.ins) {
		ins := i.ins[pos]
		switch ins.Operation {
		case OperationIsEqual:
			i.scratch[ins.Ret] = i.operandsEqual(ins.Operand1, ins.Operand2)
		case OperationIsNotEqual:
			i.scratch[ins.Ret] = !i.operandsEqual(ins.Operand1, ins.Operand2)
		case OperationIsGreaterThan:
			i.scratch[ins.Ret] = i.intFromOperand(ins.Operand1) > i.intFromOperand(ins.Operand2)
		case OperationIsGreaterThanOrEqual:
			i.scratch[ins.Ret] = i.intFromOperand(ins.Operand1) >= i.intFromOperand(ins.Operand2)
		case OperationIsLessThan:
			i.scratch[ins.Ret] = i.intFromOperand(ins.Operand1) < i.intFromOperand(ins.Operand2)
		case OperationIsLessThanOrEqual:
			i.scratch[ins.Ret] = i.intFromOperand(ins.Operand1) <= i.intFromOperand(ins.Operand2)
		case OperationContains:
			i.scratch[ins.Ret] = strings.Contains(i.stringFromOperand(ins.Operand1), i.stringFromOperand(ins.Operand2))
		case OperationDoesNotContain:
			i.scratch[ins.Ret] = !strings.Contains(i.stringFromOperand(ins.Operand1), i.stringFromOperand(ins.Operand2))
		case OperationMatches:
			i.scratch[ins.Ret] = i.regexpFromOperand(ins.Operand2).MatchString(i.stringFromOperand(ins.Operand1))
		case OperationDoesNotMatch:
			i.scratch[ins.Ret] = !i.regexpFromOperand(ins.Operand2).MatchString(i.stringFromOperand(ins.Operand1))
		case OperationJumpIfZero:
			sv := i.scratchVarFromOperand(ins.Operand1)
			i.scratch[ins.Ret] = !sv
			if !sv {
				pos = i.instructionPositionFromOperand(ins.Operand2)
				continue
			}
		case OperationJumpIfNotZero:
			sv := i.scratchVarFromOperand(ins.Operand1)
			i.scratch[ins.Ret] = sv
			if sv {
				pos = i.instructionPositionFromOperand(ins.Operand2)
				continue
			}
		case OperationAddScore:
			name := i.scoreNameFromOperand(ins.Operand1)
			val := i.intFromOperand(ins.Operand2)
			i.scores[name] += val
		case OperationSubScore:
			name := i.scoreNameFromOperand(ins.Operand1)
			val := i.intFromOperand(ins.Operand2)
			i.scores[name] -= val
		case OperationSetScore:
			name := i.scoreNameFromOperand(ins.Operand1)
			val := i.intFromOperand(ins.Operand2)
			i.scores[name] = val
		case OperationNegate:
			val := i.scratchVarFromOperand(ins.Operand1)
			i.scratch[ins.Ret] = !val
		case OperationExit:
			break Loop
		case OperationNoop:
			// Nothing
		default:
			i.setErr(fmt.Errorf("unexpected operation %v", ins.Operation))
		}
		if i.err != nil {
			break Loop
		}
		pos++
	}
}

func (i *Executor) Scores() map[string]int {
	return i.scores
}

func (i *Executor) Err() error {
	return i.err
}

func (i *Executor) setErr(err error) {
	if i.err == nil {
		i.err = err
	}
}

func (i *Executor) operandsEqual(op1, op2 Operand) bool {
	switch o := op1.(type) {
	case IntOperand:
		return o.Value == i.intFromOperand(op2)
	case StringOperand:
		return o.Value == i.stringFromOperand(op2)
	case ScoreOperand:
		return i.scores[o.Name] == i.intFromOperand(op2)
	case VarOperand:
		return i.vars[o.Name] == i.stringFromOperand(op2)
	default:
		i.setErr(fmt.Errorf("unexpected operand of type %T for equality check", op1))
	}
	return false
}

func (i *Executor) intFromOperand(op Operand) (v int) {
	switch o := op.(type) {
	case IntOperand:
		v = o.Value
	case ScoreOperand:
		v = i.scores[o.Name]
	default:
		i.setErr(fmt.Errorf("could not coerce operand of type %T into int", op))
	}
	return
}

func (i *Executor) stringFromOperand(op Operand) (s string) {
	switch o := op.(type) {
	case StringOperand:
		s = o.Value
	case VarOperand:
		s = i.vars[o.Name]
	default:
		i.setErr(fmt.Errorf("could not coerce operand of type %T into string", op))
	}
	return
}

func (i *Executor) regexpFromOperand(op Operand) (r *regexp.Regexp) {
	switch o := op.(type) {
	case RegexpOperand:
		r = o.Value
	default:
		i.setErr(fmt.Errorf("could not coerce operand of type %T into Regexp", op))
	}
	return
}

func (i *Executor) instructionPositionFromOperand(op Operand) (p int) {
	switch o := op.(type) {
	case InstructionPositionOperand:
		p = o.Pos
	default:
		i.setErr(fmt.Errorf("could not coerce operand of type %T into instruction position", op))
	}
	return
}

func (i *Executor) scratchVarFromOperand(op Operand) (b bool) {
	switch o := op.(type) {
	case ScratchOperand:
		b = i.scratch[o.Pos]
	default:
		i.setErr(fmt.Errorf("could not coerce operand of type %T into scratch variable", op))
	}
	return
}

func (i *Executor) scoreNameFromOperand(op Operand) (s string) {
	switch o := op.(type) {
	case ScoreOperand:
		s = o.Name
	default:
		i.setErr(fmt.Errorf("could not coerce operand of type %T into score", op))
	}
	return
}
