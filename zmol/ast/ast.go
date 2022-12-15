package ast

import "github.com/ariaghora/zmol/zmol/lexer"

type Node interface {
	Literal() string
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

type VarrAssignmentStatement struct {
	Token lexer.ZTok
	Name  *Identifier
	Value Expression
}

func (ls *VarrAssignmentStatement) statementNode()  {}
func (ls *VarrAssignmentStatement) Literal() string { return ls.Token.Text }

type Identifier struct {
	Token lexer.ZTok
	Value string
}

func (i *Identifier) Literal() string { return i.Value }
func (i *Identifier) expressionNode() {}
