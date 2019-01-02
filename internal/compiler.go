package internal

import "fmt"

type Compiler struct {
	ins        *InstructionsBuffer
	scratchPos uint
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
		panic("could not resolve score change or rule")
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
		panic("could not resolve condition or expression")
	}
}

func (c *Compiler) evaluateCondition(cond Condition, res ScratchPosition) {
	c.ins.Append(Instruction{
		Operation: operationFromEqualityString(cond.Op),
		Ret:       res,
		Operand1:  operandFromMixedValue(cond.LeftValue),
		Operand2:  operandFromMixedValue(cond.RightValue),
	})
}

func (c *Compiler) evaluateConsequences(cons Consequences) {
	for _, cs := range cons.Consequences {
		c.evaluateConsequent(cs)
	}
}

func (c *Compiler) evaluateConsequent(cons Consequent) {
	switch {
	case cons.ScoreChange != nil:
		c.evaluateScoreChange(*cons.ScoreChange)
	case cons.SubRule != nil:
		c.evaluateRule(*cons.SubRule)
	default:
		panic("could not resolve score change or sub rule")
	}
}

func (c *Compiler) evaluateScoreChange(sc ScoreChange) {
	c.ins.Append(Instruction{
		Operation: operationFromScoreChange(sc),
		Operand1:  ScoreOperand{Name: sc.Score.Name},
		Operand2:  operandFromIntValue(sc.Value),
	})
}

func operationFromEqualityString(s string) (op Operation) {
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
		panic("unknown operation: " + s)
	}
	return
}

func operandFromMixedValue(mv MixedValue) (op Operand) {
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
		panic("unresolvable operand: " + fmt.Sprintf("%+v", mv))
	}
	return
}

func operandFromIntValue(iv IntValue) (op Operand) {
	switch {
	case iv.Int != nil:
		op = IntOperand{Value: *iv.Int}
	case iv.Score != nil:
		op = ScoreOperand{Name: (*iv.Score).Name}
	default:
		panic("unresolvable operand: " + fmt.Sprintf("%+v", iv))
	}
	return
}

func operationFromScoreChange(sc ScoreChange) (op Operation) {
	switch sc.Operator {
	case "+=":
		op = OperationAddScore
	case "-=":
		op = OperationSubScore
	case "=":
		op = OperationSetScore
	default:
		panic("unknown operation: " + sc.Operator)
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

