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
		err := c.compileInfixExpression(node)
		if err != nil {
			return err
		}
	case *ast.PrefixExpression:
		err := c.compilePrefixExpression(node)
		if err != nil {
			return err
		}
	case *ast.IntegerLiteral:
		intVal := val.INT(node.Value)
		c.emit(bytecode.OpConstant, c.addConstant(intVal))
	case *ast.FloatLiteral:
		floatVal := val.FLOAT(node.Value)
		c.emit(bytecode.OpConstant, c.addConstant(floatVal))
	case *ast.StringLiteral:
		stringVal := val.STRING(node.Value)
		c.emit(bytecode.OpConstant, c.addConstant(stringVal))
	case *ast.BooleanLiteral:
		boolVal := val.BOOL(node.Value).Value
		if boolVal {
			c.emit(bytecode.OpTrue)
		} else {
			c.emit(bytecode.OpFalse)
		}
	}

	return nil
}

func (c *Compiler) compilePrefixExpression(node *ast.PrefixExpression) error {
	err := c.Compile(node.Right)
	if err != nil {
		return err
	}

	fmt.Println(node.Operator)

	switch node.Operator {
	case "-":
		c.emit(bytecode.OpNeg)
		return nil
	}
	return fmt.Errorf("unknown operator %s", node.Operator)
}

func (c *Compiler) compileInfixExpression(node *ast.InfixExpression) error {
	if err := c.Compile(node.Left); err != nil {
		return err
	}

	if err := c.Compile(node.Right); err != nil {
		return err
	}

	op, found := map[string]bytecode.Opcode{
		"+":  bytecode.OpAdd,
		"-":  bytecode.OpSub,
		"*":  bytecode.OpMul,
		"/":  bytecode.OpDiv,
		"%":  bytecode.OpMod,
		"==": bytecode.OpEqual,
		"!=": bytecode.OpNotEqual,
		">":  bytecode.OpGreaterThan,
		"<":  bytecode.OpLessThan,
		">=": bytecode.OpGreaterThanEqual,
		"<=": bytecode.OpLessThanEqual,
	}[node.Operator]

	if !found {
		return fmt.Errorf("unknown operator %s", node.Operator)
	}
	c.emit(op)
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
