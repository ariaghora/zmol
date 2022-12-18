package parser

import (
	"strconv"

	"github.com/ariaghora/zmol/pkg/ast"
	"github.com/ariaghora/zmol/pkg/lexer"
)

// Following constants are used to determine the precedence of the operators
const (
	_ int = iota
	PrecLowest
	PrecAssign // =
	PrecEquals // ==
	PrecGtLt   // > or <
	PrecAddSub // +
	PrecProd   // *
	PrecPrefix // -X or !X
	PrecCall   // myFunction(X)
)

var precedences = map[lexer.TokType]int{
	lexer.TokAssign: PrecAssign,
	lexer.TokEq:     PrecEquals,
	lexer.TokPlus:   PrecAddSub,
	lexer.TokMinus:  PrecAddSub,
	lexer.TokSlash:  PrecProd,
	lexer.TokAster:  PrecProd,
	lexer.TokLt:     PrecGtLt,
	lexer.TokGt:     PrecGtLt,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)
type Parser struct {
	l *lexer.ZLex

	curTok  lexer.ZTok
	peekTok lexer.ZTok
	tokIdx  int

	infixParseFns  map[lexer.TokType]infixParseFn
	prefixParseFns map[lexer.TokType]prefixParseFn
	errors         []string
}

func NewParser(l *lexer.ZLex) *Parser {
	p := &Parser{
		l:      l,
		tokIdx: 0,
	}

	// Read two tokens, so curTok and peekTok are both set
	p.nextToken()

	p.prefixParseFns = make(map[lexer.TokType]prefixParseFn)
	p.registerPrefix(lexer.TokIdent, p.parseIdentifier)
	p.registerPrefix(lexer.TokInt, p.parseIntegerLiteral)
	p.registerPrefix(lexer.TokFloat, p.parseFloatLiteral)
	p.registerPrefix(lexer.TokPlus, p.parserPrefixExpression)
	p.registerPrefix(lexer.TokMinus, p.parserPrefixExpression)
	p.registerPrefix(lexer.TokAt, p.parseFuncLiteral)

	p.infixParseFns = make(map[lexer.TokType]infixParseFn)
	p.registerInfix(lexer.TokPlus, p.parseInfixExpression)
	p.registerInfix(lexer.TokMinus, p.parseInfixExpression)
	p.registerInfix(lexer.TokAster, p.parseInfixExpression)
	p.registerInfix(lexer.TokSlash, p.parseInfixExpression)
	// p.registerInfix(lexer.TokEq, p.parseInfixExpression)

	return p
}

func (p *Parser) nextToken() {
	p.curTok = p.l.Tokens[p.tokIdx]
	if p.tokIdx < len(p.l.Tokens)-1 {
		p.tokIdx++
	}
	p.peekTok = p.l.Tokens[p.tokIdx]
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curTok.Type != lexer.TokEOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curTok.Type {
	case lexer.TokLet:
		return p.parseVarAssign()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVarAssign() *ast.VarrAssignmentStatement {
	statement := &ast.VarrAssignmentStatement{Token: p.curTok}

	if !p.expectPeek(lexer.TokIdent) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Text}

	if !p.expectPeek(lexer.TokAssign) {
		return nil
	}

	p.nextToken()

	statement.Value = p.parseExpression(PrecLowest)

	if p.peekTok.Type == lexer.TokNewLine {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curTok}

	stmt.Expression = p.parseExpression(PrecLowest)

	if p.peekTok.Type == lexer.TokNewLine {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(prec int) ast.Expression {
	prefix := p.prefixParseFns[p.curTok.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	for !(p.peekTok.Type == lexer.TokNewLine || p.peekTok.Type == lexer.TokEOF) && prec < p.peekPrec() {
		infix := p.infixParseFns[p.peekTok.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curTok, Value: p.curTok.Text}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curTok}

	value, err := strconv.ParseInt(p.curTok.Text, 0, 64)
	if err != nil {
		msg := "Could not parse %q as integer"
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curTok}

	value, err := strconv.ParseFloat(p.curTok.Text, 64)
	if err != nil {
		msg := "Could not parse %q as float"
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parserPrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curTok,
		Operator: p.curTok.Text,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PrecPrefix)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curTok,
		Operator: p.curTok.Text,
		Left:     left,
	}

	prec := p.curPrec()
	p.nextToken()
	expression.Right = p.parseExpression(prec)

	return expression
}

func (p *Parser) parseFuncLiteral() ast.Expression {
	lit := &ast.FuncLiteral{Token: p.curTok}

	if !p.expectPeek(lexer.TokLParen) {
		return nil
	}

	lit.Parameters = p.parseFuncParameters()

	if !p.expectPeek(lexer.TokColon) {
		return nil
	}

	// if p.expectPeek(lexer.TokNewLine) {
	if p.peekTok.Type == lexer.TokNewLine {
		p.skipLinebreak()
		lit.Multiline = true
		lit.Body = p.parseBlockStatement()
	} else {
		p.nextToken()
		lit.Multiline = false
		lit.BodySingle = p.parseExpressionStatement()
	}

	return lit
}

func (p *Parser) parseFuncParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTok.Type == lexer.TokRParen {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curTok, Value: p.curTok.Text}
	identifiers = append(identifiers, ident)

	for p.peekTok.Type == lexer.TokComma {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curTok, Value: p.curTok.Text}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(lexer.TokRParen) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curTok}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !(p.curTok.Type == lexer.TokEnd || p.curTok.Type == lexer.TokEOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) peekError(t lexer.TokType) {
	msg := "Expected next token to be %s, got %s instead"
	p.errors = append(p.errors, msg)
}

func (p *Parser) expectPeek(t lexer.TokType) bool {
	if p.peekTok.Type == t {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) registerPrefix(tokType lexer.TokType, fn prefixParseFn) {
	p.prefixParseFns[tokType] = fn
}

func (p *Parser) registerInfix(tokType lexer.TokType, fn infixParseFn) {
	p.infixParseFns[tokType] = fn
}

func (p *Parser) curPrec() int {
	if p, ok := precedences[p.curTok.Type]; ok {
		return p
	}

	return PrecLowest
}

func (p *Parser) peekPrec() int {
	if p, ok := precedences[p.peekTok.Type]; ok {
		return p
	}

	return PrecLowest
}

func (p *Parser) skipLinebreak() {
	for p.peekTok.Type == lexer.TokNewLine {
		p.nextToken()
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}
