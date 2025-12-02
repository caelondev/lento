package lexer

import "fmt"

type Token struct {
	TokenType TokenType
	Lexeme string
	Literal any
	Line uint
} 

func NewToken(TokenType TokenType, Lexeme string, Literal any, Line uint) *Token {
	return &Token{
		TokenType,
		Lexeme,
		Literal,
		Line,
	}
}

func (t *Token) String() {
	fmt.Printf("\nTokenType: %v\nLexeme: %s\nLiteral: %v\nLine: %d\n", TokenTypeString[t.TokenType], t.Lexeme, t.Literal,t.Line)
}
