package internal

import (
	"fmt"

	"github.com/pkg/errors"
)

type Compiler struct {
	ins        *InstructionsBuffer
	scratchPos uint
	err        error
}

func NewCompiler() *Compiler {
	return &Compiler{
		ins: &InstructionsBuffer{},
	}
}

func (c *Compiler) Compile(root Root) {
	for _, s := range root.Statements {
		c.evaluateStatement(s)
	}
	return
}

func (c *Compiler) evaluateStatement(s Statement) {
	switch {
	case s.ScoreChange != nil:
		c.evaluateScoreChange(*s.ScoreChange)
	case s.Rule != nil:
		c.evaluateRule(*s.Rule)
	default:
		c.setErr(fmt.Errorf("could not resolve score change or rule from %+v", s))
	}
}

func (c *Compiler) evaluateRule(rule Rule) {
	scratch := c.nextScratchPos()
	c.evaluateExpression(rule.Expression, scratch)
	pos := c.ins.Reserve()
	c.evaluateConsequences(rule.Consequences)
	c.ins.Replace(pos, Instruction{
		Operation: OperationJumpIfZero,
		Operand1:  ScratchOperand{Pos: scratch},
		Operand2:  InstructionPositionOperand{Pos: c.ins.Head()},
	})
	c.ins.Append(Instruction{
		Operation: OperationNoop,
	})
}

func (c *Compiler) evaluateExpression(e Expression, res ScratchPosition) {
	if len(e.Or) > 1 {
		reserved := map[int]ScratchPosition{}
		for _, or := range e.Or {
			inner := c.nextScratchPos()
			c.evaluateOrExpression(or, inner)
			pos := c.ins.Reserve()
			reserved[pos] = inner
		}
		for p, pos := range reserved {
			c.ins.Replace(p, Instruction{
				Operation: OperationJumpIfNotZero,
				Ret:       res,
				Operand1:  ScratchOperand{Pos: pos},
				Operand2:  InstructionPositionOperand{Pos: c.ins.Head()},
			})
		}
	} else {
		c.evaluateOrExpression(e.Or[0], res)
	}
}

func (c *Compiler) evaluateOrExpression(or OrExpression, res ScratchPosition) {
	if len(or.And) > 1 {
		reserved := map[int]ScratchPosition{}
		for _, coe := range or.And {
			inner := c.nextScratchPos()
			c.evaluateConditionOrExpression(coe, inner)
			pos := c.ins.Reserve()
			reserved[pos] = inner
		}
		for p, pos := range reserved {
			c.ins.Replace(p, Instruction{
				Operation: OperationJumpIfZero,
				Ret:       res,
				Operand1:  ScratchOperand{Pos: pos},
				Operand2:  InstructionPositionOperand{Pos: c.ins.Head()},
			})
		}
		c.ins.Append(Instruction{
			Operation: OperationNegate,
			Ret:       res,
			Operand1:  ScratchOperand{Pos: res},
		})
	} else {
		c.evaluateConditionOrExpression(or.And[0], res)
	}
}

func (c *Compiler) evaluateConditionOrExpression(coe ConditionOrExpression, res ScratchPosition) {
	switch {
	case coe.Condition != nil:
		c.evaluateCondition(*coe.Condition, res)
	case coe.Expression != nil:
		c.evaluateExpression(*coe.Expression, res)
	default:
		c.setErr(fmt.Errorf("could not resolve condition or expression from %+v", coe))
	}
}

func (c *Compiler) evaluateCondition(cond Condition, res ScratchPosition) {
	op, err := operationFromEqualityString(cond.Op)
	if err != nil {
		c.setErr(errors.Wrap(err, "failed to map condition operation"))
		return
	}
	operand1, err := operandFromMixedValue(cond.LeftValue)
	if err != nil {
		c.setErr(errors.Wrap(err, "failed to map first operand"))
		return
	}
	operand2, err := operandFromMixedValue(cond.RightValue)
	if err != nil {
		c.setErr(errors.Wrap(err, "failed to map second operand"))
		return
	}
	c.ins.Append(Instruction{
		Operation: op,
		Ret:       res,
		Operand1:  operand1,
		Operand2:  operand2,
	})
}

func (c *Compiler) evaluateConsequences(cons Consequences) {
	for _, s := range cons.Consequences {
		c.evaluateStatement(s)
	}
}

func (c *Compiler) evaluateScoreChange(sc ScoreChange) {
	op, err := operationFromScoreChange(sc)
	if err != nil {
		c.setErr(errors.Wrap(err, "failed to map score change operation"))
		return
	}
	operand1 := ScoreOperand{Name: sc.Score.Name}
	operand2, err := operandFromIntValue(sc.Value)
	if err != nil {
		c.setErr(errors.Wrap(err, "failed to map second operand"))
		return
	}
	c.ins.Append(Instruction{
		Operation: op,
		Operand1:  operand1,
		Operand2:  operand2,
	})
}

func operationFromEqualityString(s string) (op Operation, err error) {
	switch s {
	case "==":
		op = OperationIsEqual
	case "!=":
		op = OperationIsNotEqual
	case ">":
		op = OperationIsGreaterThan
	case ">=":
		op = OperationIsGreaterThanOrEqual
	case "<":
		op = OperationIsLessThan
	case "<=":
		op = OperationIsLessThanOrEqual
	case "contains":
		op = OperationContains
	case "doesnotcontain":
		op = OperationDoesNotContain
	default:
		err = fmt.Errorf("unknown operation %s", s)
	}
	return
}

func operandFromMixedValue(mv MixedValue) (op Operand, err error) {
	switch {
	case mv.Var != nil:
		op = VarOperand{Name: *mv.Var}
	case mv.String != nil:
		op = StringOperand{Value: *mv.String}
	case mv.Int != nil:
		op = IntOperand{Value: *mv.Int}
	case mv.Score != nil:
		op = ScoreOperand{Name: (*mv.Score).Name}
	default:
		err = fmt.Errorf("unresolvable mixed value %+v", mv)
	}
	return
}

func operandFromIntValue(iv IntValue) (op Operand, err error) {
	switch {
	case iv.Int != nil:
		op = IntOperand{Value: *iv.Int}
	case iv.Score != nil:
		op = ScoreOperand{Name: (*iv.Score).Name}
	default:
		err = fmt.Errorf("unresolvable int value %+v", iv)
	}
	return
}

func operationFromScoreChange(sc ScoreChange) (op Operation, err error) {
	switch sc.Operator {
	case "+=":
		op = OperationAddScore
	case "-=":
		op = OperationSubScore
	case "=":
		op = OperationSetScore
	default:
		err = fmt.Errorf("unknown operation %+v", sc)
	}
	return
}

func (c *Compiler) nextScratchPos() ScratchPosition {
	c.scratchPos++
	return ScratchPosition(c.scratchPos)
}

func (c *Compiler) Instructions() []Instruction {
	return c.ins.Instructions()
}

func (c *Compiler) Err() error {
	return c.err
}

func (c *Compiler) setErr(err error) {
	if c.err == nil {
		c.err = err
	}
}
