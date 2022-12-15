package ast

import "github.com/ariaghora/zmol/zmol/lexer"

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
	return ls.Token.Text + " " + ls.Name.Str() + " = " + ls.Value.Str() + " "
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
