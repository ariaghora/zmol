package parser

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ariaghora/zmol/pkg/ast"
	"github.com/ariaghora/zmol/pkg/lexer"
)

// Following constants are used to determine the precedence of the operators
const (
	_ int = iota
	PrecLowest
	PrecAssign   // 2: =
	PrecTernary  // 2: ?:
	PrecPipeline // 2: |> >-
	PrecOr       // 3: ||
	PrecAnd      // 4: &&
	PrecEquals   // 8: ==
	PrecGtLt     // 9: > or <
	PrecAddSub   // 11: +
	PrecProd     // 12: *
	PrecPrefix   // 14: -X or !X
	PrecCall     // 17: myFunction(X)
	PrecIndex    // 18: array[index]
)

var precedences = map[lexer.TokType]int{
	lexer.TokAssign: PrecAssign,

	// Boolean operators
	lexer.TokEq:    PrecEquals,
	lexer.TokNotEq: PrecEquals,
	lexer.TokGt:    PrecGtLt,
	lexer.TokLt:    PrecGtLt,
	lexer.TokGTE:   PrecGtLt,
	lexer.TokLTE:   PrecGtLt,

	// Logical operators
	lexer.TokAnd: PrecAnd,
	lexer.TokOr:  PrecOr,

	// Pipeline operator
	lexer.TokPipe:   PrecPipeline,
	lexer.TokFilter: PrecPipeline,

	// Arithmetic operators
	lexer.TokPlus:   PrecAddSub,
	lexer.TokMinus:  PrecAddSub,
	lexer.TokSlash:  PrecProd,
	lexer.TokAster:  PrecProd,
	lexer.TokMod:    PrecProd,
	lexer.TokLParen: PrecCall,

	// Index operator
	lexer.TokLBrac: PrecIndex,

	// Ternary operator
	lexer.TokQuestion: PrecTernary,
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
	p.registerPrefix(lexer.TokTrue, p.parseBooleanLiteral)
	p.registerPrefix(lexer.TokFalse, p.parseBooleanLiteral)
	p.registerPrefix(lexer.TokPlus, p.parserPrefixExpression)
	p.registerPrefix(lexer.TokMinus, p.parserPrefixExpression)
	p.registerPrefix(lexer.TokAt, p.parseFuncLiteral)
	p.registerPrefix(lexer.TokLParen, p.parseGroupedExpression)
	p.registerPrefix(lexer.TokLBrac, p.parseListLiteral)

	p.infixParseFns = make(map[lexer.TokType]infixParseFn)
	p.registerInfix(lexer.TokPlus, p.parseInfixExpression)
	p.registerInfix(lexer.TokMinus, p.parseInfixExpression)
	p.registerInfix(lexer.TokAster, p.parseInfixExpression)
	p.registerInfix(lexer.TokSlash, p.parseInfixExpression)
	p.registerInfix(lexer.TokAssign, p.parseInfixExpression)
	p.registerInfix(lexer.TokMod, p.parseInfixExpression)

	// List access
	p.registerInfix(lexer.TokLBrac, p.parseIndexExpression)

	// Boolean operators
	p.registerInfix(lexer.TokEq, p.parseInfixExpression)
	p.registerInfix(lexer.TokNotEq, p.parseInfixExpression)
	p.registerInfix(lexer.TokLt, p.parseInfixExpression)
	p.registerInfix(lexer.TokLTE, p.parseInfixExpression)
	p.registerInfix(lexer.TokGt, p.parseInfixExpression)
	p.registerInfix(lexer.TokGTE, p.parseInfixExpression)

	// Logical operators
	p.registerInfix(lexer.TokAnd, p.parseInfixExpression)
	p.registerInfix(lexer.TokOr, p.parseInfixExpression)

	// Pipeline operator
	p.registerInfix(lexer.TokPipe, p.parseInfixExpression)
	p.registerInfix(lexer.TokFilter, p.parseInfixExpression)

	p.registerInfix(lexer.TokLParen, p.parseCallExpression)
	p.registerInfix(lexer.TokQuestion, p.parserTernaryExpression)
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
		if p.curTok.Type == lexer.TokNewLine {
			p.nextToken()
			continue
		}
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
	case lexer.TokIter:
		return p.parseIter()
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

func (p *Parser) parseIter() *ast.IterStatement {
	// Iterate over a list
	statement := &ast.IterStatement{Token: p.curTok}

	p.nextToken()

	list := p.parseExpression(PrecLowest)

	if !p.expectPeek(lexer.TokAs) {
		return nil
	}

	// Get the identifier
	p.nextToken()

	if p.curTok.Type != lexer.TokIdent {
		return nil
	}

	statement.Ident = &ast.Identifier{Token: p.curTok, Value: p.curTok.Text}

	if !p.expectPeek(lexer.TokColon) {
		return nil
	}

	p.nextToken()

	body := p.parseBlockStatement()

	statement.List = list
	statement.Body = body

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

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curTok, Value: p.curTok.Type == lexer.TokTrue}
}

func (p *Parser) parseListLiteral() ast.Expression {
	list := &ast.ListLiteral{Token: p.curTok}
	list.Elements = p.parseExpressionList(lexer.TokRBrac)
	return list
}

func (p *Parser) parseExpressionList(end lexer.TokType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTok.Type == end {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(PrecLowest))

	for p.peekTok.Type == lexer.TokComma {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(PrecLowest))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
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
		// lit.Multiline = true
		lit.Body = p.parseBlockStatement()
	} else {
		p.nextToken()
		// lit.Multiline = false
		stmt := p.parseExpressionStatement()
		lit.Body = &ast.BlockStatement{Token: p.curTok, Statements: []ast.Statement{stmt}}
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
		p.skipLinebreak()
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	if p.curTok.Type == lexer.TokEOF {
		msg := "Unexpected end of file while parsing block statement"
		p.errors = append(p.errors, msg)
		fmt.Println(msg)
		os.Exit(1)
	}

	return block
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    p.curTok,
		Function: function,
	}

	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parserTernaryExpression(condition ast.Expression) ast.Expression {
	exp := &ast.TernaryExpression{
		Token:     p.curTok,
		Condition: condition,
	}

	p.nextToken()
	exp.Consequence = p.parseExpression(PrecLowest)

	if !p.expectPeek(lexer.TokColon) {
		return nil
	}

	p.nextToken()
	exp.Alternative = p.parseExpression(PrecLowest)

	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTok.Type == lexer.TokRParen {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(PrecLowest))

	for p.peekTok.Type == lexer.TokComma {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(PrecLowest))
	}

	if !p.expectPeek(lexer.TokRParen) {
		return nil
	}

	return args
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(PrecLowest)

	if !p.expectPeek(lexer.TokRParen) {
		return nil
	}

	return exp
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{
		Token: p.curTok,
		Left:  left,
	}

	p.nextToken()
	exp.Index = p.parseExpression(PrecLowest)

	if !p.expectPeek(lexer.TokRBrac) {
		return nil
	}

	return exp
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
