package internal

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
)

type InstructionsGenerator struct {
	buf        *InstructionsBuffer
	scratchPos uint
	err        error
}

func NewInstructionsGenerator() *InstructionsGenerator {
	return &InstructionsGenerator{
		buf: &InstructionsBuffer{},
	}
}

func (ig *InstructionsGenerator) Generate(root Root) {
	for _, s := range root.Statements {
		ig.evaluateStatement(s)
	}
	ig.buf.Append(Instruction{
		Operation: OperationNoop,
	})
}

func (ig *InstructionsGenerator) evaluateStatement(s Statement) {
	switch {
	case s.ScoreChange != nil:
		ig.evaluateScoreChange(*s.ScoreChange)
	case s.Rule != nil:
		ig.evaluateRule(*s.Rule)
	default:
		ig.setErr(fmt.Errorf("could not resolve score change or rule from %+v", s))
	}
}

func (ig *InstructionsGenerator) evaluateRule(rule Rule) {
	scratch := ig.nextScratchPos()
	ig.evaluateExpression(rule.Expression, scratch)
	pos := ig.buf.Reserve()
	ig.evaluateConsequences(rule.Consequences)
	ig.buf.Replace(pos, Instruction{
		Operation: OperationJumpIfZero,
		Operand1:  ScratchOperand{Pos: scratch},
		Operand2:  InstructionPositionOperand{Pos: ig.buf.Head()},
	})
}

func (ig *InstructionsGenerator) evaluateExpression(e Expression, res ScratchPosition) {
	if len(e.Or) > 1 {
		reserved := map[int]ScratchPosition{}
		for _, or := range e.Or {
			inner := ig.nextScratchPos()
			ig.evaluateOrExpression(or, inner)
			pos := ig.buf.Reserve()
			reserved[pos] = inner
		}
		for p, pos := range reserved {
			ig.buf.Replace(p, Instruction{
				Operation: OperationJumpIfNotZero,
				Ret:       res,
				Operand1:  ScratchOperand{Pos: pos},
				Operand2:  InstructionPositionOperand{Pos: ig.buf.Head()},
			})
		}
	} else {
		ig.evaluateOrExpression(e.Or[0], res)
	}
}

func (ig *InstructionsGenerator) evaluateOrExpression(or OrExpression, res ScratchPosition) {
	if len(or.And) > 1 {
		reserved := map[int]ScratchPosition{}
		for _, coe := range or.And {
			inner := ig.nextScratchPos()
			ig.evaluateConditionOrExpression(coe, inner)
			pos := ig.buf.Reserve()
			reserved[pos] = inner
		}
		for p, pos := range reserved {
			ig.buf.Replace(p, Instruction{
				Operation: OperationJumpIfZero,
				Ret:       res,
				Operand1:  ScratchOperand{Pos: pos},
				Operand2:  InstructionPositionOperand{Pos: ig.buf.Head()},
			})
		}
		ig.buf.Append(Instruction{
			Operation: OperationNegate,
			Ret:       res,
			Operand1:  ScratchOperand{Pos: res},
		})
	} else {
		ig.evaluateConditionOrExpression(or.And[0], res)
	}
}

func (ig *InstructionsGenerator) evaluateConditionOrExpression(coe ConditionOrExpression, res ScratchPosition) {
	switch {
	case coe.Condition != nil:
		ig.evaluateCondition(*coe.Condition, res)
	case coe.Expression != nil:
		ig.evaluateExpression(*coe.Expression, res)
	default:
		ig.setErr(fmt.Errorf("could not resolve condition or expression from %+v", coe))
	}
}

func (ig *InstructionsGenerator) evaluateCondition(cond Condition, res ScratchPosition) {
	op, err := operationFromEqualityString(cond.Op)
	if err != nil {
		ig.setErr(errors.Wrap(err, "failed to map condition operation"))
		return
	}
	operand1, err := operandFromMixedValue(cond.LeftValue)
	if err != nil {
		ig.setErr(errors.Wrap(err, "failed to map first operand"))
		return
	}
	operand2, err := operandFromMixedValue(cond.RightValue)
	if err != nil {
		ig.setErr(errors.Wrap(err, "failed to map second operand"))
		return
	}
	ig.buf.Append(Instruction{
		Operation: op,
		Ret:       res,
		Operand1:  operand1,
		Operand2:  operand2,
	})
}

func (ig *InstructionsGenerator) evaluateConsequences(cons Consequences) {
	for _, s := range cons.Consequences {
		ig.evaluateStatement(s)
	}
}

func (ig *InstructionsGenerator) evaluateScoreChange(sc ScoreChange) {
	op, err := operationFromScoreChange(sc)
	if err != nil {
		ig.setErr(errors.Wrap(err, "failed to map score change operation"))
		return
	}
	operand1 := ScoreOperand{Name: sc.Score.Name}
	operand2, err := operandFromIntValue(sc.Value)
	if err != nil {
		ig.setErr(errors.Wrap(err, "failed to map second operand"))
		return
	}
	ig.buf.Append(Instruction{
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
	case "matches":
		op = OperationMatches
	case "doesnotmatch":
		op = OperationDoesNotMatch
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
	case mv.Regexp != nil:
		op, err = buildRegexpOperand(*mv.Regexp)
	default:
		err = fmt.Errorf("unresolvable mixed value %+v", mv)
	}
	return
}

func buildRegexpOperand(s string) (RegexpOperand, error) {
	rg, err := regexp.Compile(s)
	if err != nil {
		return RegexpOperand{}, errors.Wrap(err, "regex compile failed")
	}
	return RegexpOperand{Value: rg}, nil
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

func (ig *InstructionsGenerator) nextScratchPos() ScratchPosition {
	ig.scratchPos++
	return ScratchPosition(ig.scratchPos)
}

func (ig *InstructionsGenerator) Instructions() []Instruction {
	return ig.buf.Instructions()
}

func (ig *InstructionsGenerator) Err() error {
	return ig.err
}

func (ig *InstructionsGenerator) setErr(err error) {
	if ig.err == nil {
		ig.err = err
	}
}
