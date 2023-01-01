package compiler

import (
	"fmt"

	"github.com/ariaghora/zmol/pkg/ast"
	"github.com/ariaghora/zmol/pkg/bytecode"
	"github.com/ariaghora/zmol/pkg/val"
)

type Bytecode struct {
	Instructions bytecode.Instructions
	Constants    []val.ZValue
}

type Compiler struct {
	instructions bytecode.Instructions
	constants    []val.ZValue
}

func NewCompiler() *Compiler {
	return &Compiler{
		instructions: bytecode.Instructions{},
		constants:    []val.ZValue{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(bytecode.OpPop)
	case *ast.InfixExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(bytecode.OpAdd)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.IntegerLiteral:
		intVal := val.INT(node.Value)
		c.emit(bytecode.OpConstant, c.addConstant(intVal))
	}

	return nil
}

func (c *Compiler) emit(op bytecode.Opcode, operands ...int) int {
	instruction := bytecode.Make(op, operands...)
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, instruction...)

	return posNewInstruction
}

func (c *Compiler) addConstant(val val.ZValue) int {
	c.constants = append(c.constants, val)
	return len(c.constants) - 1
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}
