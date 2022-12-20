package lexer

import (
	"errors"
	"unicode"
)

type TokType string

const (
	TokPlus     TokType = "+"
	TokMinus            = "-"
	TokSlash            = "/"
	TokAster            = "*"
	TokMod              = "%"
	TokAnd              = "&&"
	TokOr               = "||"
	TokBitAnd           = "&"
	TokBitOr            = "|"
	TokEq               = "=="
	TokNotEq            = "!="
	TokGt               = ">"
	TokGTE              = ">="
	TokLt               = "<"
	TokLTE              = "<="
	TokDot              = "."
	TokComma            = ","
	TokLBrac            = "["
	TokRBrac            = "]"
	TokAssign           = "="
	TokAt               = "@"
	TokLParen           = "("
	TokRParen           = ")"
	TokColon            = ":"
	TokNot              = "!"
	TokQuestion         = "?"
	TokLet              = "let"
	TokNewLine          = "LINEBREAK"
	TokEOF              = "EOF"
	TokIdent            = "IDENT"
	TokIf               = "if"
	TokElse             = "else"
	TokTrue             = "true"
	TokFalse            = "false"
	TokInt              = "INT"
	TokFloat            = "FLOAT"
	TokEnd              = "END"
)

type ZTok struct {
	Type TokType
	Text string
}

type ZLex struct {
	i      int
	Tokens []ZTok
	code   string
}

var SingularTokOps = map[rune]TokType{
	'@':  TokAt,
	'+':  TokPlus,
	'-':  TokMinus,
	'*':  TokAster,
	'/':  TokSlash,
	'%':  TokMod,
	'&':  TokBitAnd,
	'|':  TokBitOr,
	'>':  TokGt,
	'<':  TokLt,
	'.':  TokDot,
	',':  TokComma,
	'[':  TokLBrac,
	']':  TokRBrac,
	'=':  TokAssign,
	'(':  TokLParen,
	')':  TokRParen,
	'\n': TokNewLine,
	':':  TokColon,
	'!':  TokNot,
	'?':  TokQuestion,
}

var KeywordTok = map[string]TokType{
	"let":   TokLet,
	"if":    TokIf,
	"else":  TokElse,
	"true":  TokTrue,
	"false": TokFalse,
	"end":   TokEnd,
}

func (z *ZLex) addIdent() {
	var nChar int
	start := z.i
	for z.i+nChar < len(z.code) && (unicode.IsLetter(rune(z.code[z.i+nChar])) ||
		unicode.IsDigit(rune(z.code[z.i+nChar])) ||
		z.code[z.i+nChar] == '_') {
		nChar++
	}
	if tokType, ok := KeywordTok[z.code[start:z.i+nChar]]; ok {
		z.addTok(tokType, nChar)
	} else {
		z.addTok(TokIdent, nChar)
	}
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
	if hasDot {
		z.addTok(TokFloat, nChar)
	} else {
		z.addTok(TokInt, nChar)
	}
}

func (z *ZLex) addTok(tokType TokType, nChar int) {
	z.Tokens = append(z.Tokens, ZTok{
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
		if unicode.IsSpace(rune(z.code[z.i])) && z.code[z.i] != '\n' {
			z.skipWhitespace()
		} else if tokType, ok := SingularTokOps[rune(z.code[z.i])]; ok {
			if z.code[z.i] == '>' && z.i+1 < len(z.code) && z.code[z.i+1] == '=' {
				z.addTok(TokGTE, 2)
			} else if z.code[z.i] == '<' && z.i+1 < len(z.code) && z.code[z.i+1] == '=' {
				z.addTok(TokLTE, 2)
			} else if z.code[z.i] == '=' && z.i+1 < len(z.code) && z.code[z.i+1] == '=' {
				z.addTok(TokEq, 2)
			} else if z.code[z.i] == '!' && z.i+1 < len(z.code) && z.code[z.i+1] == '=' {
				z.addTok(TokNotEq, 2)
			} else if z.code[z.i] == '&' && z.i+1 < len(z.code) && z.code[z.i+1] == '&' {
				z.addTok(TokAnd, 2)
			} else if z.code[z.i] == '|' && z.i+1 < len(z.code) && z.code[z.i+1] == '|' {
				z.addTok(TokOr, 2)
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
	z.Tokens = append(z.Tokens, ZTok{Type: TokEOF})
	return nil
}

// func (z *ZLex) NextToken() ZTok {
// 	if len(z.tokens) == 0 {
// 		z.Lex()
// 	}
// 	tok := z.tokens[0]
// 	z.tokens = z.tokens[1:]
// 	return tok
// }

func NewLexer(code string) *ZLex {
	return &ZLex{code: code}
}
