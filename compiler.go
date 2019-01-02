package brulee

import "fmt"

func Compile(r Root) Program {
	c := &compiler{
		ins: &instructionsBuffer{},
	}
	c.evaluateRoot(r)
	return Program{
		ins: c.ins.Instructions(),
	}
}

type compiler struct {
	ins        *instructionsBuffer
	scratchPos uint
}

func (c *compiler) evaluateRoot(root Root) {
	for _, s := range root.Statements {
		c.evaluateRootStatement(s)
	}
	return
}

func (c *compiler) evaluateRootStatement(rs RootStatement) {
	switch {
	case rs.ScoreChange != nil:
		c.evaluateScoreChange(*rs.ScoreChange)
	case rs.Rule != nil:
		c.evaluateRule(*rs.Rule)
	default:
		panic("could not resolve score change or rule")
	}
}

func (c *compiler) evaluateRule(rule Rule) {
	scratch := c.nextScratchPos()
	c.evaluateExpression(rule.Expression, scratch)
	pos := c.ins.Reserve()
	c.evaluateConsequences(rule.Consequences)
	c.ins.Replace(pos, Instruction{
		Operation: OperationJumpIfZero,
		Operand1:  ScratchOperand{pos: scratch},
		Operand2:  InstructionPositionOperand{pos: c.ins.Head()},
	})
	c.ins.Append(Instruction{
		Operation: OperationNoop,
	})
}

func (c *compiler) evaluateExpression(e Expression, res ScratchPosition) {
	if len(e.Or) > 1 {
		positions := map[int]ScratchPosition{}
		for _, or := range e.Or {
			inner := c.nextScratchPos()
			c.evaluateOrExpression(or, inner)
			pos := c.ins.Reserve()
			positions[pos] = inner
		}
		for p, pos := range positions {
			c.ins.Replace(p, Instruction{
				Operation: OperationJumpIfNotZero,
				Ret:       res,
				Operand1:  ScratchOperand{pos: pos},
				Operand2:  InstructionPositionOperand{pos: c.ins.Head()},
			})
		}
	} else {
		c.evaluateOrExpression(e.Or[0], res)
	}
}

func (c *compiler) evaluateOrExpression(or OrExpression, res ScratchPosition) {
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
				Operand1:  ScratchOperand{pos: pos},
				Operand2:  InstructionPositionOperand{pos: c.ins.Head()},
			})
		}
		c.ins.Append(Instruction{
			Operation: OperationNegate,
			Ret:       res,
			Operand1:  ScratchOperand{pos: res},
		})
	} else {
		c.evaluateConditionOrExpression(or.And[0], res)
	}
}

func (c *compiler) evaluateConditionOrExpression(coe ConditionOrExpression, res ScratchPosition) {
	switch {
	case coe.Condition != nil:
		c.evaluateCondition(*coe.Condition, res)
	case coe.Expression != nil:
		c.evaluateExpression(*coe.Expression, res)
	default:
		panic("could not resolve condition or expression")
	}
}

func (c *compiler) evaluateCondition(cond Condition, res ScratchPosition) {
	c.ins.Append(Instruction{
		Operation: operationFromEqualityString(cond.Op),
		Ret:       res,
		Operand1:  operandFromMixedValue(cond.LeftValue),
		Operand2:  operandFromMixedValue(cond.RightValue),
	})
}

func (c *compiler) evaluateConsequences(cons Consequences) {
	for _, cs := range cons.Consequences {
		c.evaluateConsequent(cs)
	}
}

func (c *compiler) evaluateConsequent(cons Consequent) {
	switch {
	case cons.ScoreChange != nil:
		c.evaluateScoreChange(*cons.ScoreChange)
	case cons.SubRule != nil:
		c.evaluateRule(*cons.SubRule)
	default:
		panic("could not resolve score change or sub rule")
	}
}

func (c *compiler) evaluateScoreChange(sc ScoreChange) {
	c.ins.Append(Instruction{
		Operation: operationFromScoreChange(sc),
		Operand1:  ScoreOperand{name: sc.Score.Name},
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
		op = VarOperand{name: *mv.Var}
	case mv.String != nil:
		op = StringOperand{value: *mv.String}
	case mv.Int != nil:
		op = IntOperand{value: *mv.Int}
	case mv.Score != nil:
		op = ScoreOperand{name: (*mv.Score).Name}
	default:
		panic("unresolvable operand: " + fmt.Sprintf("%+v", mv))
	}
	return
}

func operandFromIntValue(iv IntValue) (op Operand) {
	switch {
	case iv.Int != nil:
		op = IntOperand{value: *iv.Int}
	case iv.Score != nil:
		op = ScoreOperand{name: (*iv.Score).Name}
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

func (c *compiler) nextScratchPos() ScratchPosition {
	c.scratchPos++
	return ScratchPosition(c.scratchPos)
}
