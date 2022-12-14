package lexer

import (
	"errors"
	"unicode"
)

type TokType string

const (
	TokPlus   TokType = "+"
	TokMinus          = "-"
	TokGT             = ">"
	TokGTE            = ">="
	TokLT             = "<"
	TokLTE            = "<="
	TokDot            = "."
	TokComma          = ","
	TokLBrac          = "["
	TokRBrac          = "]"
	TokAssign         = "="
	TokEOF            = "EOF"
	TokIdent          = "IDENT"
	TokNumber         = "NUMBER"
)

type ZTok struct {
	Type TokType
	Text string
}

type ZLex struct {
	i      int
	tokens []ZTok
	code   string
}

var SingularTokOps = map[rune]TokType{
	'+': TokPlus,
	'-': TokMinus,
	'>': TokGT,
	'<': TokLT,
	'.': TokDot,
	',': TokComma,
	'[': TokLBrac,
	']': TokRBrac,
	'=': TokAssign,
}

func (z *ZLex) addIdent() {
	var nChar int
	for z.i+nChar < len(z.code) && (unicode.IsLetter(rune(z.code[z.i+nChar])) || unicode.IsDigit(rune(z.code[z.i+nChar]))) {
		nChar++
	}
	z.addTok(TokIdent, nChar)
}

func (z *ZLex) addNumber() {
	// handle integer and float. Make sure there is only one dot
	var nChar int
	var hasDot bool
	for z.i+nChar < len(z.code) && (unicode.IsDigit(rune(z.code[z.i+nChar])) || z.code[z.i+nChar] == '.') {
		if z.code[z.i+nChar] == '.' {
			if hasDot {
				panic(errors.New("invalid float"))
			}
			hasDot = true
		}
		nChar++
	}
	z.addTok(TokNumber, nChar)
}

func (z *ZLex) addTok(tokType TokType, nChar int) {
	z.tokens = append(z.tokens, ZTok{
		Type: tokType,
		Text: z.code[z.i : z.i+nChar],
	})
	z.i += nChar
}

func (z *ZLex) skipWhitespace() {
	for z.i < len(z.code) && unicode.IsSpace(rune(z.code[z.i])) {
		z.i++
	}
}

func (z *ZLex) Lex() error {
	for z.i < len(z.code) {
		if unicode.IsSpace(rune(z.code[z.i])) {
			z.skipWhitespace()
		} else if tokType, ok := SingularTokOps[rune(z.code[z.i])]; ok {
			if z.code[z.i] == '>' && z.i+1 < len(z.code) && z.code[z.i+1] == '=' {
				z.addTok(TokGTE, 2)
			} else if z.code[z.i] == '<' && z.i+1 < len(z.code) && z.code[z.i+1] == '=' {
				z.addTok(TokLTE, 2)
			} else {
				z.addTok(tokType, 1)
			}
		} else if unicode.IsLetter(rune(z.code[z.i])) {
			z.addIdent()
		} else if unicode.IsDigit(rune(z.code[z.i])) {
			z.addNumber()
		} else {
			return errors.New("Invalid token: " + string(z.code[z.i]))
		}
	}
	z.tokens = append(z.tokens, ZTok{Type: TokEOF})
	return nil
}

func NewLexer(code string) *ZLex {
	return &ZLex{code: code}
}
