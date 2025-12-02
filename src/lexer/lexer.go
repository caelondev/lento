package lexer

import (
	"fmt"
	"strconv"

	errorhandler "github.com/caelondev/lento/src/error-handler"
)

type Lexer struct {
	SourceCode   []rune
	Tokens       []*Token
	ErrorHandler *errorhandler.ErrorHandler

	Start   int
	Current int
	Line    uint
}

func NewLexer(sourceCode string, errorHandler *errorhandler.ErrorHandler) *Lexer {
	return &Lexer{
		SourceCode:   []rune(sourceCode),
		Tokens:       make([]*Token, 0, len(sourceCode)+1),
		ErrorHandler: errorHandler,
		Start:        0,
		Current:      0,
		Line:         1,
	}
}

func (l *Lexer) Tokenize() []*Token {
	for !l.isEOF() {
		l.Start = l.Current
		l.AnalyzeTokens()
	}

	l.Tokens = append(l.Tokens, &Token{
		TokenType: EOF,
		Lexeme:    "END_OF_FILE",
		Line:      l.Line,
	})

	return l.Tokens
}

func (l *Lexer) isEOF() bool {
	return l.Current >= len(l.SourceCode)
}

func (l *Lexer) AnalyzeTokens() {
	char := l.advance()
	switch char {
	case '(':
		l.addToken(LEFT_PARENTHESIS)
	case ')':
		l.addToken(RIGHT_PARENTHESIS)
	case '{':
		l.addToken(LEFT_BRACE)
	case '}':
		l.addToken(RIGHT_BRACE)
	case '[':
		l.addToken(LEFT_BRACKET)
	case ']':
		l.addToken(RIGHT_BRACKET)
	case '.':
		l.addToken(DOT)
	case ',':
		l.addToken(COMMA)
	case '*':
		l.handleCompound(STAR, STAR_EQUALS)
	case '%':
		l.handleCompound(MODULO, MODULO_EQUALS)
	case '-':
		l.handleMinus()
	case '+':
		l.handlePlus()
	case '/':
		l.handleSlash()
	case '"', '\'':
		l.handleString(char)
	case '`':
		l.handleMultilineString()

	case ' ', '\r', '\t':
		// ignore ---
		break
	case '\n':
		l.Line++

	default:
		if isNumber(char) {
			l.handleNumbers()
		} else if isAlphabet(char) { // Handle identifiers and keywords
			l.handleIdentifier()
		} else {
			l.ErrorHandler.ReportError(
				"Lexer-Tokenizer",
				fmt.Sprintf("Unrecognized Token found '%c'", char),
				int(l.Line),
				65,
			)
		}
	}
}

func (l *Lexer) handleMultilineString() {
	startLine := l.Line
	for !l.isEOF() && l.peek() != '`' {
		if l.peek() == '\n' {
			l.Line++
		}
		l.advance()
	}
	if l.isEOF() {
		l.ErrorHandler.ReportError(
			"Lexer-Tokenizer",
			"Unterminated multiline string",
			int(l.Line),
			65,
			)
	}
	l.match('`') // Check and eat '`' closing string --- 
	literal := string(l.SourceCode[l.Start+1 : l.Current-1])
	l.addTokenWithLiteral(STRING, literal, startLine)
}

func (l *Lexer) handleString(char rune) {
	for !l.isEOF() && l.peek() != char /* closing char string */ {
		if l.peek() == '\n' {
			break
		}
		l.advance()
	}

	if l.isEOF() || l.peek() == '\n' {
		l.ErrorHandler.ReportError(
			"Lexer-Tokenizer",
			"Unterminated non-multiline string",
			int(l.Line),
			65,
		)
	}

	l.match(char)

	literal := string(l.SourceCode[l.Start +1 : l.Current -1])
	l.addTokenWithLiteral(STRING, literal, 0)
}

func (l *Lexer) handleSlash() {
	if l.match('/') { // Oneline comment ---
		for l.peek() != '\n' && !l.isEOF() {
			l.advance() // Keep eating tokens until EOF or new line ---
		}
	} else if l.match('*') { // Multiline comment ---

		for !l.isEOF() && !(l.peek() == '*' && l.peekNext() == '/') {
			if l.peek() == '\n' {
				l.Line++
			}

			l.advance()
		}

		l.match('*')
		l.match('/')

	} else {
		l.handleCompound(SLASH, SLASH_EQUALS)
	}
}

func (l *Lexer) handleIdentifier() {
	for isAlphanumeric(l.peek()) && !l.isEOF() {
		l.advance() // Eat all tokens
	}

	value := string(l.SourceCode[l.Start:l.Current])
	keyword, exists := RESERVED_KEYWORDS[value]
	if exists {
		l.addToken(keyword)
	} else {
		l.addTokenWithLiteral(IDENTIFIER, value, 0)
	}
}

func (l *Lexer) handleMinus() {
	if l.peek() == '-' {
		l.advance() // Eat '-' token ---
		l.addToken(MINUS_MINUS)
	} else {
		l.handleCompound(MINUS, MINUS_EQUALS)
	}
}

func (l *Lexer) handlePlus() {
	if l.peek() == '+' {
		l.advance() // Eat '+' token ---
		l.addToken(PLUS_PLUS)
	} else {
		l.handleCompound(PLUS, PLUS_EQUALS)
	}
}

func (l *Lexer) handleCompound(regular, compound TokenType) {
	if l.peek() == '=' {
		l.advance() // Eat '=' token ---
		l.addToken(compound)
	} else {
		l.addToken(regular)
	}
}

func (l *Lexer) handleNumbers() {
	for isNumber(l.peek()) && !l.isEOF() {
		l.advance() // Keep consuming numbers ---
	}

	// Handle floats
	if l.peek() == '.' {
		if !isNumber(l.peekNext()) {
			l.ErrorHandler.ReportError(
				"Lexer-Tokenizer",
				"Expected number after '.' operator...",
				int(l.Line),
				65,
			)
		}

		l.match('.')

		for isNumber(l.peek()) {
			l.advance()
		}
	}

	value := string(l.SourceCode[l.Start:l.Current])
	parsedNumber, error := strconv.ParseFloat(value, 64)
	if error != nil {
		l.ErrorHandler.Report(
			int(l.Line),
			fmt.Sprintf("Failed to parse %s into a float whilst tokenizing", value),
		)
	}

	l.addTokenWithLiteral(NUMBER, parsedNumber, 0)
}

func (l *Lexer) advance() rune {
	token := l.SourceCode[l.Current]
	l.Current++
	return token
}

func (l *Lexer) peek() rune {
	return l.SourceCode[l.Current]
}

func (l *Lexer) peekNext() rune {
	if l.Current+1 >= len(l.SourceCode) {
		return 0
	}
	return l.SourceCode[l.Current+1]
}

func (l *Lexer) match(expected rune) bool {
	if l.isEOF() {
		return false
	}
	if l.SourceCode[l.Current] != expected {
		return false
	}

	l.Current++
	return true
}

func (l *Lexer) addToken(tokenType TokenType) {
	l.addTokenWithLiteral(tokenType, nil, 0)
}

func (l *Lexer) addTokenWithLiteral(tokenType TokenType, literal any, line uint) {
	if line == 0 {
		l.Tokens = append(l.Tokens, &Token{
			TokenType: tokenType,
			Lexeme:    string(l.SourceCode[l.Start:l.Current]),
			Literal:   literal,
			Line:      l.Line,
		})
	} else {
		l.Tokens = append(l.Tokens, &Token{
			TokenType: tokenType,
			Lexeme:    string(l.SourceCode[l.Start:l.Current]),
			Literal:   literal,
			Line:      line,
		})
	}
}

func isNumber(char rune) bool {
	return char >= '0' && char <= '9'
}

func isAlphabet(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		char == '_'
}

func isAlphanumeric(char rune) bool {
	return isAlphabet(char) || isNumber(char)
}

func (l *Lexer) expectError(expected rune, errorMessage string) {
	if l.peek() != expected {
		l.ErrorHandler.ReportError(
			"Lexer-Tokenizer",
			errorMessage,
			int(l.Line),
			65,
		)
	}
}
