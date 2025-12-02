package parser

import (
	"github.com/caelondev/lento/src/ast"
	"github.com/caelondev/lento/src/lexer"
)

func parseStatement(p *parser) ast.Statement {
	if p.isEOF() {
		return nil
	}

	statementFunction, exists := statementLU[p.currentTokenType()]

	if exists {
		return statementFunction(p)
	}

	expression := parseExpression(p, DEFAULT_BP)

	p.expect(lexer.SEMICOLON)

	return &ast.ExpressionStatement{
		Expression: expression,
		Line: p.line,
	}
}
