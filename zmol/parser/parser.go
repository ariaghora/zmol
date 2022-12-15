package parser

import (
	"strconv"

	"github.com/ariaghora/zmol/zmol/ast"
	"github.com/ariaghora/zmol/zmol/lexer"
)

const (
	_ int = iota
	PrecLowest
	PrecEquals // ==
	PrecGtLt   // > or <
	PrecAddSub // +
	PrecProd   // *
	PrecPrefix // -X or !X
	PrecCall   // myFunction(X)
)

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
	case lexer.TokAt:
		return p.parseVarAssignStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVarAssignStatement() *ast.VarrAssignmentStatement {
	statement := &ast.VarrAssignmentStatement{Token: p.curTok}

	if !p.expectPeek(lexer.TokIdent) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Text}

	if !p.expectPeek(lexer.TokAssign) {
		return nil
	}

	// TODO: We're skipping the expressions until we
	// encounter a line break or EOF
	for !(p.curTok.Type == lexer.TokNewLine || p.curTok.Type == lexer.TokEOF) {
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

func (p *Parser) Errors() []string {
	return p.errors
}
