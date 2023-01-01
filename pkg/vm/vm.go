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
		}
	}
	return nil
}

func (vm *VM) LastPoppedStackElem() val.ZValue {
	if len(vm.stack) == 0 {
		return nil
	}

	return vm.stack[vm.sp-1]
}

func (vm *VM) push(v val.ZValue) error {
	if vm.sp == StackSize {
		return errors.New("stack overflow")
	}

	vm.stack[vm.sp] = v
	vm.sp++
	return nil
}
