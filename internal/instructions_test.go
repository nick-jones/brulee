package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstruction_StringSlice(t *testing.T) {
	testCases := []struct {
		name     string
		ins      Instruction
		expected []string
	}{
		{
			name: "just operation",
			ins: Instruction{
				Operation: OperationNoop,
			},
			expected: []string{"NOOP", "", "", ""},
		},
		{
			name: "operation and return",
			ins: Instruction{
				Operation: OperationNoop,
				Ret:       ScratchPosition(1),
			},
			expected: []string{"NOOP", "$1", "", ""},
		},
		{
			name: "operation, return and operand 1",
			ins: Instruction{
				Operation: OperationNoop,
				Ret:       ScratchPosition(1),
				Operand1:  IntOperand{Value: 1},
			},
			expected: []string{"NOOP", "$1", "int(1)", ""},
		},
		{
			name: "operation, return and both operands",
			ins: Instruction{
				Operation: OperationNoop,
				Ret:       ScratchPosition(1),
				Operand1:  IntOperand{Value: 1},
				Operand2:  IntOperand{Value: 2},
			},
			expected: []string{"NOOP", "$1", "int(1)", "int(2)"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			assert.Equal(tt, tc.expected, tc.ins.StringSlice())
		})
	}
}
