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
		stmt := parseStatement(parser)
		
		if errorHandler.HadError {
			return ast.BlockStatement{
				Body: nil,
				Line: parser.line,
			}
		}
		
		if stmt != nil {
			body = append(body, stmt)
		}
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
	if p.position >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
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

		return token
	}

	return p.advance()
}

func (p *parser) advance() *lexer.Token {
	token := p.currentToken()
	if token.Line > p.line {
		p.line = token.Line
	}
	
	if p.position < len(p.tokens)-1 {
		p.position++
	}
	
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
