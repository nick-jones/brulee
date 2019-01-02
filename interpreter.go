package brulee

import (
	"fmt"
	"strings"
)

type Interpreter struct {
	ins     []Instruction
	vars    map[string]string
	scratch map[ScratchPosition]bool
	scores  Scores
	err     error
}

func (i *Interpreter) Interpret() {
	pos := 0
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
		case OperationJumpIfZero:
			sv := i.scratchVarFromOperand(ins.Operand1)
			if !sv {
				pos = i.instructionPositionFromOperand(ins.Operand2)
				i.scratch[ins.Ret] = true
				continue
			}
		case OperationJumpIfNotZero:
			sv := i.scratchVarFromOperand(ins.Operand1)
			if sv {
				pos = i.instructionPositionFromOperand(ins.Operand2)
				i.scratch[ins.Ret] = true
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
		case OperationNoop:
			// Nothing
		default:
			i.setErr(fmt.Errorf("unexpected operation %v", ins.Operation))
		}
		if i.err != nil {
			break
		}
		pos++
	}
}

func (i *Interpreter) Scores() Scores {
	return i.scores
}

func (i *Interpreter) Err() error {
	return i.err
}

func (i *Interpreter) setErr(err error) {
	if i.err == nil {
		i.err = err
	}
}

func (i *Interpreter) operandsEqual(op1, op2 Operand) bool {
	switch o := op1.(type) {
	case IntOperand:
		return o.value == i.intFromOperand(op2)
	case StringOperand:
		return o.value == i.stringFromOperand(op2)
	case ScoreOperand:
		return i.scores[o.name] == i.intFromOperand(op2)
	case VarOperand:
		return i.vars[o.name] == i.stringFromOperand(op2)
	default:
		i.setErr(fmt.Errorf("unexpected operand of type %T for equality check", op1))
	}
	return false
}

func (i *Interpreter) intFromOperand(op Operand) (v int) {
	switch o := op.(type) {
	case IntOperand:
		v = o.value
	case ScoreOperand:
		v = i.scores[o.name]
	default:
		i.setErr(fmt.Errorf("could not coerce operand of type %T into int", op))
	}
	return
}

func (i *Interpreter) stringFromOperand(op Operand) (s string) {
	switch o := op.(type) {
	case StringOperand:
		s = o.value
	case VarOperand:
		s = i.vars[o.name]
	default:
		i.setErr(fmt.Errorf("could not coerce operand of type %T into string", op))
	}
	return
}

func (i *Interpreter) instructionPositionFromOperand(op Operand) (p int) {
	switch o := op.(type) {
	case InstructionPositionOperand:
		p = o.pos
	default:
		i.setErr(fmt.Errorf("could not coerce operand of type %T into instruction position", op))
	}
	return
}

func (i *Interpreter) scratchVarFromOperand(op Operand) (b bool) {
	switch o := op.(type) {
	case ScratchOperand:
		b = i.scratch[o.pos]
	default:
		i.setErr(fmt.Errorf("could not coerce operand of type %T into scratch variable", op))
	}
	return
}

func (i *Interpreter) scoreNameFromOperand(op Operand) (s string) {
	switch o := op.(type) {
	case ScoreOperand:
		s = o.name
	default:
		i.setErr(fmt.Errorf("could not coerce operand of type %T into score", op))
	}
	return
}
