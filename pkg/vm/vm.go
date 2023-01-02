package vm

import (
	"errors"
	"fmt"

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
		case bytecode.OpTrue, bytecode.OpFalse:
			err := vm.push(val.BOOL(op == bytecode.OpTrue))
			if err != nil {
				return err
			}
		case bytecode.OpNeg:
			err := vm.executeUnaryOperation(op)
			if err != nil {
				return err
			}
		case bytecode.OpAdd, bytecode.OpSub, bytecode.OpMul, bytecode.OpDiv, bytecode.OpMod:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case bytecode.OpEqual, bytecode.OpNotEqual, bytecode.OpGreaterThan, bytecode.OpLessThan, bytecode.OpGreaterThanEqual, bytecode.OpLessThanEqual:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}
		case bytecode.OpPop:
			vm.pop()
		}
	}
	return nil
}

func (vm *VM) executeUnaryOperation(op bytecode.Opcode) error {
	right := vm.pop()

	fmt.Printf("right: %T\n", right)
	rightOperand, ok := right.(val.ZArithOperand)
	if !ok {
		return errors.New("unsupported operand type")
	}

	var result val.ZValue = nil
	switch op {
	case bytecode.OpNeg:
		result = rightOperand.Neg()
	}

	vm.push(result)
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

// handle comparison operations
func (vm *VM) executeComparison(op bytecode.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	leftOperand, ok := left.(val.ZComparable)
	if !ok {
		return errors.New("unsupported operand type: " + string(left.Type()))
	}

	var result val.ZValue = nil
	switch op {
	case bytecode.OpEqual:
		result = leftOperand.Equal(right)
	case bytecode.OpNotEqual:
		result = leftOperand.NotEqual(right)
	case bytecode.OpGreaterThan:
		result = leftOperand.GreaterThan(right)
	case bytecode.OpLessThan:
		result = leftOperand.LessThan(right)
	case bytecode.OpGreaterThanEqual:
		result = leftOperand.GreaterThanEqual(right)
	case bytecode.OpLessThanEqual:
		result = leftOperand.LessThanEqual(right)
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
