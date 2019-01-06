package internal

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecutor_Execute(t *testing.T) {
	testCases := []struct {
		name     string
		ins      []Instruction
		expected map[string]int
	}{
		{
			name: "is equal check, pass",
			ins: []Instruction{
				{Operation: OperationIsEqual, Ret: 1, Operand1: StringOperand{Value: "a"}, Operand2: StringOperand{Value: "a"}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{"x": 1},
		},
		{
			name: "is equal check, fail",
			ins: []Instruction{
				{Operation: OperationIsEqual, Ret: 1, Operand1: StringOperand{Value: "a"}, Operand2: StringOperand{Value: "b"}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{},
		},
		{
			name: "is not equal check, pass",
			ins: []Instruction{
				{Operation: OperationIsNotEqual, Ret: 1, Operand1: StringOperand{Value: "a"}, Operand2: StringOperand{Value: "b"}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{"x": 1},
		},
		{
			name: "is equal check, fail",
			ins: []Instruction{
				{Operation: OperationIsNotEqual, Ret: 1, Operand1: StringOperand{Value: "a"}, Operand2: StringOperand{Value: "a"}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{},
		},
		{
			name: "is greater than check, pass",
			ins: []Instruction{
				{Operation: OperationIsGreaterThan, Ret: 1, Operand1: IntOperand{Value: 2}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{"x": 1},
		},
		{
			name: "is greater than check, fail",
			ins: []Instruction{
				{Operation: OperationIsGreaterThan, Ret: 1, Operand1: IntOperand{Value: 2}, Operand2: IntOperand{Value: 2}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{},
		},
		{
			name: "is greater than or equal check, pass (gt)",
			ins: []Instruction{
				{Operation: OperationIsGreaterThanOrEqual, Ret: 1, Operand1: IntOperand{Value: 2}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{"x": 1},
		},
		{
			name: "is greater than or equal check, pass (eq)",
			ins: []Instruction{
				{Operation: OperationIsGreaterThanOrEqual, Ret: 1, Operand1: IntOperand{Value: 2}, Operand2: IntOperand{Value: 2}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{"x": 1},
		},
		{
			name: "is greater than or equal check, fail",
			ins: []Instruction{
				{Operation: OperationIsGreaterThanOrEqual, Ret: 1, Operand1: IntOperand{Value: 1}, Operand2: IntOperand{Value: 2}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{},
		},

		{
			name: "is less than check, pass",
			ins: []Instruction{
				{Operation: OperationIsLessThan, Ret: 1, Operand1: IntOperand{Value: 1}, Operand2: IntOperand{Value: 2}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{"x": 1},
		},
		{
			name: "is less than check, fail",
			ins: []Instruction{
				{Operation: OperationIsLessThan, Ret: 1, Operand1: IntOperand{Value: 2}, Operand2: IntOperand{Value: 2}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{},
		},
		{
			name: "is less than or equal check, pass (lt)",
			ins: []Instruction{
				{Operation: OperationIsLessThanOrEqual, Ret: 1, Operand1: IntOperand{Value: 1}, Operand2: IntOperand{Value: 2}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{"x": 1},
		},
		{
			name: "is less than or equal check, pass (eq)",
			ins: []Instruction{
				{Operation: OperationIsLessThanOrEqual, Ret: 1, Operand1: IntOperand{Value: 2}, Operand2: IntOperand{Value: 2}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{"x": 1},
		},
		{
			name: "is less than or equal check, fail",
			ins: []Instruction{
				{Operation: OperationIsLessThanOrEqual, Ret: 1, Operand1: IntOperand{Value: 2}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{},
		},
		{
			name: "contains check, pass",
			ins: []Instruction{
				{Operation: OperationContains, Ret: 1, Operand1: StringOperand{Value: "abc"}, Operand2: StringOperand{Value: "b"}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{"x": 1},
		},
		{
			name: "contains check, fail",
			ins: []Instruction{
				{Operation: OperationContains, Ret: 1, Operand1: StringOperand{Value: "abc"}, Operand2: StringOperand{Value: "d"}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{},
		},
		{
			name: "matches check, pass",
			ins: []Instruction{
				{Operation: OperationMatches, Ret: 1, Operand1: StringOperand{Value: "abc"}, Operand2: RegexpOperand{Value: regexp.MustCompile("b")}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{"x": 1},
		},
		{
			name: "matches check, fail",
			ins: []Instruction{
				{Operation: OperationMatches, Ret: 1, Operand1: StringOperand{Value: "abc"}, Operand2: RegexpOperand{Value: regexp.MustCompile("d")}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{},
		},
		{
			name: "does not match check, pass",
			ins: []Instruction{
				{Operation: OperationDoesNotMatch, Ret: 1, Operand1: StringOperand{Value: "abc"}, Operand2: RegexpOperand{Value: regexp.MustCompile("d")}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{"x": 1},
		},
		{
			name: "does not match check, fail",
			ins: []Instruction{
				{Operation: OperationDoesNotMatch, Ret: 1, Operand1: StringOperand{Value: "abc"}, Operand2: RegexpOperand{Value: regexp.MustCompile("b")}},
				{Operation: OperationJumpIfZero, Operand1: ScratchOperand{Pos: 1}, Operand2: InstructionPositionOperand{Pos: 3}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{},
		},
		{
			name: "score adjustments",
			ins: []Instruction{
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationAddScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationAddScore, Operand1: ScoreOperand{Name: "y"}, Operand2: IntOperand{Value: 1}},
			},
			expected: map[string]int{"x": 2, "y": 1},
		},
		{
			name: "scratch negate",
			ins: []Instruction{
				{Operation: OperationIsEqual, Ret: 1, Operand1: StringOperand{Value: "a"}, Operand2: StringOperand{Value: "a"}},
				{Operation: OperationNegate, Ret: 2, Operand1: ScratchOperand{Pos: 1}},
				{Operation: OperationJumpIfNotZero, Operand1: ScratchOperand{Pos: 2}, Operand2: InstructionPositionOperand{Pos: 4}},
				{Operation: OperationSetScore, Operand1: ScoreOperand{Name: "x"}, Operand2: IntOperand{Value: 1}},
				{Operation: OperationNoop},
			},
			expected: map[string]int{"x": 1},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			ex := NewExecutor(tc.ins, map[string]string{})
			ex.Execute()
			assert.NoError(tt, ex.Err())
			assert.Equal(tt, tc.expected, ex.Scores())
		})
	}
}
