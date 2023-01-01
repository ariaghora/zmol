package bytecode

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
	}

	for _, test := range tests {
		instruction := Make(test.op, test.operands...)

		if len(instruction) != len(test.expected) {
			t.Errorf("wrong length. want=%d, got=%d", len(test.expected), len(instruction))
		}

		for i, b := range test.expected {
			if instruction[i] != b {
				t.Errorf("wrong byte at position %d. want=%d, got=%d", i, b, instruction[i])
			}
		}
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
	}

	expected := `0000 OpAdd
0001 OpConstant 2
0004 OpConstant 65535
`

	concatenated := Instructions{}
	for _, ins := range instructions {
		concatenated = append(concatenated, ins...)
	}

	if concatenated.String() != expected {
		t.Fatalf("wrong instructions string. want=%q, got=%q", expected, concatenated.String())
	}

}

func TestOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65534}, 2},
	}

	for _, test := range tests {
		instruction := Make(test.op, test.operands...)

		definition, ok := definitions[test.op]
		if !ok {
			t.Fatalf("opcode %d not defined", test.op)
		}

		operands, n := ReadOperands(definition, instruction[1:])
		if n != test.bytesRead {
			t.Fatalf("wrong number of bytes read. want=%d, got=%d", test.bytesRead, n)
		}

		for i, operand := range operands {
			if operand != test.operands[i] {
				t.Errorf("wrong operand at position %d. want=%d, got=%d", i, test.operands[i], operand)
			}
		}
	}

}
