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

type FloatLiteral struct {
	Token lexer.ZTok
	Value float64
}

func (fl *FloatLiteral) expressionNode() {}
func (ie *FloatLiteral) Literal() string { return ie.Token.Text }
func (ie *FloatLiteral) Str() string     { return ie.Token.Text }

type BooleanLiteral struct {
	Token lexer.ZTok
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) Literal() string { return bl.Token.Text }
func (bl *BooleanLiteral) Str() string     { return bl.Token.Text }

type ListLiteral struct {
	Token    lexer.ZTok // the '[' token
	Elements []Expression
}

func (ll *ListLiteral) expressionNode() {}
func (ll *ListLiteral) Literal() string { return ll.Token.Text }
func (ll *ListLiteral) Str() string {
	out := "["
	for _, e := range ll.Elements {
		out += e.Str() + ", "
	}
	out += "]"
	return out
}

type IndexExpression struct {
	Token lexer.ZTok // the '[' token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) Literal() string { return ie.Token.Text }
func (ie *IndexExpression) Str() string {
	return "(" + ie.Left.Str() + "[" + ie.Index.Str() + "])"
}

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

type BlockStatement struct {
	Token      lexer.ZTok // the '{' token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()  {}
func (bs *BlockStatement) Literal() string { return bs.Token.Text }
func (bs *BlockStatement) Str() string {
	out := ""
	for _, s := range bs.Statements {
		out += s.Str()
	}
	return out
}

type FuncLiteral struct {
	Token      lexer.ZTok // the 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FuncLiteral) expressionNode() {}
func (fl *FuncLiteral) Literal() string { return fl.Token.Text }
func (fl *FuncLiteral) Str() string {
	params := ""
	for _, p := range fl.Parameters {
		params += p.Str() + ", "
	}
	return fl.Token.Text + "(" + params + ")" + fl.Body.Str()
}

type CallExpression struct {
	Token     lexer.ZTok // the '(' token
	Function  Expression // Identifier or FuncLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) Literal() string { return ce.Token.Text }
func (ce *CallExpression) Str() string {
	args := ""
	for _, a := range ce.Arguments {
		args += a.Str() + ", "
	}
	return ce.Function.Str() + "(" + args + ")"
}

type IfStatement struct {
	Token       lexer.ZTok // the 'if' token
	Condition   Expression
	Consequence Statement
	Alternative Statement
}

func (is *IfStatement) statementNode()  {}
func (is *IfStatement) Literal() string { return is.Token.Text }
func (is *IfStatement) Str() string {
	out := "if " + is.Condition.Str() + " " + is.Consequence.Str()
	if is.Alternative != nil {
		out += "else " + is.Alternative.Str()
	}
	return out
}

type TernaryExpression struct {
	Token       lexer.ZTok // the '?' token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (te *TernaryExpression) expressionNode() {}
func (te *TernaryExpression) Literal() string { return te.Token.Text }
func (te *TernaryExpression) Str() string {
	return te.Condition.Str() + "?" + te.Consequence.Str() + ":" + te.Alternative.Str()
}

type IterStatement struct {
	Token lexer.ZTok // the 'iter' token
	List  Expression
	Body  *BlockStatement
	Ident *Identifier
}

func (is *IterStatement) statementNode()  {}
func (is *IterStatement) Literal() string { return is.Token.Text }
func (is *IterStatement) Str() string {
	return is.Token.Text + " " + is.List.Str() + " " + is.Body.Str()
}
