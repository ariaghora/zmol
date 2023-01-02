package bytecode

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte
type Opcode byte

const (
	OpConstant Opcode = iota
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpNeg
	OpPop
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
	OpAdd:      {"OpAdd", []int{}},
	OpSub:      {"OpSub", []int{}},
	OpMul:      {"OpMul", []int{}},
	OpDiv:      {"OpDiv", []int{}},
	OpMod:      {"OpMod", []int{}},
	OpNeg:      {"OpNeg", []int{}},
	OpPop:      {"OpPop", []int{}},
}

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, ins.formatInstruction(def, operands))

		i += 1 + read

	}
	return out.String()
}

func (ins Instructions) formatInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand length mismatch. want=%d, got=%d", operandCount, len(operands))
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d not defined", op)
	}
	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1 // 1 byte for the opcode
	for _, width := range def.OperandWidths {
		instructionLen += width
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)
	offset := 1

	for i, operand := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(operand))
		}
		offset += width
	}

	return instruction
}

func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0
	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUInt16(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

func ReadUInt16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}
