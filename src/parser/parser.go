package parser

import (
	"fmt"
	"strings"

	"github.com/caelondev/lento/src/ast"
	errorhandler "github.com/caelondev/lento/src/error-handler"
	"github.com/caelondev/lento/src/lexer"
)

type parser struct {
	tokens   []*lexer.Token
	position int
	line     uint

	errorHandler *errorhandler.ErrorHandler
}

func ProduceAST(tokens []*lexer.Token, errorHandler *errorhandler.ErrorHandler) ast.BlockStatement {
	body := make([]ast.Statement, 0)
	parser := instantiateParser(tokens, errorHandler)

	for !parser.isEOF() {
		body = append(body, parseStatement(parser))
	}

	return ast.BlockStatement{
		Body: body,
		Line: parser.line,
	}
}

func instantiateParser(tokens []*lexer.Token, errorHandler *errorhandler.ErrorHandler) *parser {
	createTokenLookups()

	return &parser{
		tokens:   tokens,
		position: 0,
		line:     1,

		errorHandler: errorHandler,
	}
}

func (p *parser) currentTokenType() lexer.TokenType {
	return p.currentToken().TokenType
}
func (p *parser) currentToken() *lexer.Token {
	return p.tokens[p.position]
}

func (p *parser) isEOF() bool {
	return p.position >= len(p.tokens) || p.currentTokenType() == lexer.EOF
}

func (p *parser) expect(expectedTypes ...lexer.TokenType) *lexer.Token {
	return p.expectError("", expectedTypes...)
}

func (p *parser) expectError(err string, expectedTypes ...lexer.TokenType) *lexer.Token {
	token := p.currentToken()
	tokenType := token.TokenType

	matched := false
	for _, t := range expectedTypes {
		if tokenType == t {
			matched = true
			break
		}
	}

	if !matched {
		if err == "" {
			var names []string
			for _, t := range expectedTypes {
				names = append(names, lexer.TokenTypeString[t])
			}
			err = fmt.Sprintf("Expected %s but got %s instead", strings.Join(names, "/"), lexer.TokenTypeString[tokenType])
		}

		p.errorHandler.ReportError(
			"Parser-Expect",
			err,
			p.line,
			errorhandler.UnexpectedTokenError,
		)
	}

	return p.advance()
}

func (p *parser) advance() *lexer.Token {
	token := p.currentToken()
	if token.Line > p.line {
		p.line = token.Line
	}
	p.position++
	return token
}

func (p *parser) synchronize() {
	for !p.isEOF() {
		if p.currentTokenType() == lexer.SEMICOLON {
			p.advance()
			return
		}
		p.advance()
	}
}

func (p *parser) ignore(tokenType lexer.TokenType) {
	if p.currentTokenType() == tokenType {
		p.advance()
	}
}
