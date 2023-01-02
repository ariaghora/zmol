package vm

import (
	"errors"

	"github.com/ariaghora/zmol/pkg/bytecode"
	"github.com/ariaghora/zmol/pkg/compiler"
	"github.com/ariaghora/zmol/pkg/val"
)

const StackSize = 2048

type VM struct {
	constants    []val.ZValue
	instructions bytecode.Instructions
	stack        []val.ZValue
	sp           int // Points to the next value. Top of stack is stack[sp-1]
}

func NewVM(bytecode *compiler.Bytecode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,
		stack:        make([]val.ZValue, StackSize),
		sp:           0,
	}
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		switch op := bytecode.Opcode(vm.instructions[ip]); op {
		case bytecode.OpConstant:
			constIndex := bytecode.ReadUInt16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case bytecode.OpAdd, bytecode.OpSub, bytecode.OpMul, bytecode.OpDiv, bytecode.OpMod:
			vm.executeBinaryOperation(op)
		case bytecode.OpPop:
			vm.pop()
		}
	}
	return nil
}

// handle arithmetic operations
func (vm *VM) executeBinaryOperation(op bytecode.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	leftOperand, ok := left.(val.ZArithOperand)
	if !ok {
		return errors.New("unsupported operand type: " + string(leftOperand.Type()))
	}

	var result val.ZValue = nil
	switch op {
	case bytecode.OpAdd:
		result = leftOperand.Add(right)
	case bytecode.OpSub:
		result = leftOperand.Sub(right)
	case bytecode.OpMul:
		result = leftOperand.Mul(right)
	case bytecode.OpDiv:
		result = leftOperand.Div(right)
	case bytecode.OpMod:
		result = leftOperand.Mod(right)
	}

	return vm.push(result)
}

func (vm *VM) StackTop() val.ZValue {
	if len(vm.stack) == 0 {
		return nil
	}

	return vm.stack[vm.sp-1]
}

func (vm *VM) LastPoppedStackElem() val.ZValue {
	return vm.stack[vm.sp]
}

func (vm *VM) pop() val.ZValue {
	if vm.sp == 0 {
		return nil
	}

	vm.sp--
	return vm.stack[vm.sp]
}

func (vm *VM) push(v val.ZValue) error {
	if vm.sp == StackSize {
		return errors.New("stack overflow")
	}

	vm.stack[vm.sp] = v
	vm.sp++
	return nil
}

func (vm *VM) Stack() []val.ZValue {
	return vm.stack
}

func (vm *VM) Sp() int {
	return vm.sp
}
