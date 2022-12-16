package ast

import "github.com/ariaghora/zmol/pkg/lexer"

type Node interface {
	Literal() string
	Str() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) Literal() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].Literal()
	}
	return ""
}

func (p *Program) Str() string {
	out := ""
	for _, s := range p.Statements {
		out += s.Str()
	}
	return out
}

type VarrAssignmentStatement struct {
	Token lexer.ZTok
	Name  *Identifier
	Value Expression
}

func (ls *VarrAssignmentStatement) statementNode()  {}
func (ls *VarrAssignmentStatement) Literal() string { return ls.Token.Text }
func (ls *VarrAssignmentStatement) Str() string {
	s := ls.Token.Text +
		ls.Name.Str() +
		" = "
	if ls.Value != nil {
		s += ls.Value.Str()
	}
	return s
}

type Identifier struct {
	Token lexer.ZTok
	Value string
}

func (i *Identifier) Literal() string { return i.Value }
func (i *Identifier) expressionNode() {}
func (i *Identifier) Str() string     { return i.Value }

type ExpressionStatement struct {
	Token      lexer.ZTok // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()  {}
func (es *ExpressionStatement) Literal() string { return es.Token.Text }
func (es *ExpressionStatement) Str() string     { return es.Expression.Str() }

type IntegerLiteral struct {
	Token lexer.ZTok
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) Literal() string { return il.Token.Text }
func (il *IntegerLiteral) Str() string     { return il.Token.Text }

type InfixExpression struct {
	Token    lexer.ZTok // the operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) Literal() string { return ie.Token.Text }
func (ie *InfixExpression) Str() string {
	return "(" + ie.Left.Str() + " " + ie.Operator + " " + ie.Right.Str() + ")"
}

type PrefixExpression struct {
	Token    lexer.ZTok // the operator token, e.g. -
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) Literal() string { return pe.Token.Text }
func (pe *PrefixExpression) Str() string {
	return "(" + pe.Operator + pe.Right.Str() + ")"
}
