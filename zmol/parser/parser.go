package parser

import (
	"github.com/ariaghora/zmol/zmol/ast"
	"github.com/ariaghora/zmol/zmol/lexer"
)

type Parser struct {
	l *lexer.ZLex

	curTok  lexer.ZTok
	peekTok lexer.ZTok
	tokIdx  int

	errors []string
}

func NewParser(l *lexer.ZLex) *Parser {
	p := &Parser{
		l:      l,
		tokIdx: 0,
	}

	// Read two tokens, so curTok and peekTok are both set
	p.nextToken()

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
		return nil
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

func (p *Parser) Errors() []string {
	return p.errors
}
