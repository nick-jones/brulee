package internal

import (
	"fmt"
	"regexp"
)

type Operation int8

const (
	OperationNoop Operation = iota
	OperationIsEqual
	OperationIsNotEqual
	OperationIsGreaterThan
	OperationIsGreaterThanOrEqual
	OperationIsLessThan
	OperationIsLessThanOrEqual
	OperationContains
	OperationDoesNotContain
	OperationMatches
	OperationDoesNotMatch
	OperationJumpIfZero
	OperationJumpIfNotZero
	OperationAddScore
	OperationSubScore
	OperationSetScore
	OperationNegate
	OperationExit
)

var operationToStringMap = map[Operation]string{
	OperationNoop:                 "NOOP",
	OperationIsEqual:              "IS_EQUAL",
	OperationIsNotEqual:           "IS_NOT_EQUAL",
	OperationIsGreaterThan:        "IS_GREATER_THAN",
	OperationIsGreaterThanOrEqual: "IS_GREATER_THAN_OR_EQUAL",
	OperationIsLessThan:           "IS_LESS_THAN",
	OperationIsLessThanOrEqual:    "IS_LESS_THAN_OR_EQUAL",
	OperationContains:             "CONTAINS",
	OperationDoesNotContain:       "DOES_NOT_CONTAIN",
	OperationMatches:              "MATCHES",
	OperationDoesNotMatch:         "DOES_NOT_MATCH",
	OperationJumpIfZero:           "JUMP_IF_ZERO",
	OperationJumpIfNotZero:        "JUMP_IF_NOT_ZERO",
	OperationAddScore:             "ADD_SCORE",
	OperationSubScore:             "SUB_SCORE",
	OperationSetScore:             "SET_SCORE",
	OperationNegate:               "NEGATE",
	OperationExit:                 "EXIT",
}

func (o Operation) String() string {
	return operationToStringMap[o]
}

type Operand interface {
	String() string
}

type RegexpOperand struct {
	Value *regexp.Regexp
}

func (ro RegexpOperand) String() string {
	return fmt.Sprintf("regexp(/%s/)", ro.Value)
}

type VarOperand struct {
	Name string
}

func (vo VarOperand) String() string {
	return fmt.Sprintf("var(%s)", vo.Name)
}

type IntOperand struct {
	Value int
}

func (io IntOperand) String() string {
	return fmt.Sprintf("int(%d)", io.Value)
}

type StringOperand struct {
	Value string
}

func (so StringOperand) String() string {
	return fmt.Sprintf(`string("%s")`, so.Value)
}

type ScratchOperand struct {
	Pos ScratchPosition
}

func (so ScratchOperand) String() string {
	return so.Pos.String()
}

type ScoreOperand struct {
	Name string
}

func (so ScoreOperand) String() string {
	return fmt.Sprintf("score(%s)", so.Name)
}

type InstructionPositionOperand struct {
	Pos int
}

func (so InstructionPositionOperand) String() string {
	return fmt.Sprintf(`->%d`, so.Pos)
}

type ScratchPosition uint

func (sv ScratchPosition) String() string {
	return fmt.Sprintf("$%d", sv)
}

type Instruction struct {
	Operation Operation
	Ret       ScratchPosition
	Operand1  Operand
	Operand2  Operand
}

func (i Instruction) StringSlice() []string {
	parts := make([]string, 4)
	parts[0] = i.Operation.String()
	if i.Ret != 0 {
		parts[1] = i.Ret.String()
	}
	if i.Operand1 != nil {
		parts[2] = i.Operand1.String()
	}
	if i.Operand2 != nil {
		parts[3] = i.Operand2.String()
	}
	return parts
}

type InstructionsBuffer struct {
	ins []Instruction
}

func (i *InstructionsBuffer) Append(in Instruction) {
	i.ins = append(i.ins, in)
}

func (i *InstructionsBuffer) Reserve() int {
	pos := len(i.ins)
	i.Append(Instruction{})
	return pos
}

func (i *InstructionsBuffer) Head() int {
	return len(i.ins)
}

func (i *InstructionsBuffer) Replace(pos int, in Instruction) {
	i.ins[pos] = in
}

func (i *InstructionsBuffer) Instructions() []Instruction {
	return i.ins
}
