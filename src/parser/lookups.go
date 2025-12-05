package parser

import (
	"github.com/caelondev/lento/src/ast"
	"github.com/caelondev/lento/src/lexer"
)

type BindingPower int

const (
	DEFAULT_BP BindingPower = iota
	COMMA
	ASSIGNMENT
	LOGICAL
	RELATIONAL
	ADDITIVE
	MULTIPLICATIVE
	UNARY
	POSTFIX
	CALL
	MEMBER
	PRIMARY
)

type StatementHandler func(p *parser) ast.Statement
type NudHandler func(p *parser) ast.Expression
type LedHandler func(p *parser, left ast.Expression, bp BindingPower) ast.Expression

type StatementLookup map[lexer.TokenType]StatementHandler
type NudLookup map[lexer.TokenType]NudHandler
type LedLookup map[lexer.TokenType]LedHandler
type BPLookup map[lexer.TokenType]BindingPower

var bindingPowerLU = BPLookup{}
var nudLU = NudLookup{}
var ledLU = LedLookup{}
var statementLU = StatementLookup{}

func led(tokenType lexer.TokenType, bp BindingPower, ledFunction LedHandler) {
	bindingPowerLU[tokenType] = bp
	ledLU[tokenType] = ledFunction
}

func nud(tokenType lexer.TokenType, nudFunction NudHandler) {
	bindingPowerLU[tokenType] = PRIMARY
	nudLU[tokenType] = nudFunction
}

func statement(tokenType lexer.TokenType, statementFunction StatementHandler) {
	bindingPowerLU[tokenType] = DEFAULT_BP
	statementLU[tokenType] = statementFunction
}

func createTokenLookups() {

	// LITERALS AND SYMBOLS ---
	nud(lexer.NUMBER, parsePrimaryExpression)
	nud(lexer.IDENTIFIER, parsePrimaryExpression)
	nud(lexer.STRING, parsePrimaryExpression)
	nud(lexer.LEFT_PARENTHESIS, parsePrimaryExpression)
	nud(lexer.LEFT_BRACE, parseObjectExpression)
	led(lexer.DOT, MEMBER, parseMemberExpression)

	// ARRAYS ---
	nud(lexer.LEFT_BRACKET, parseArrayExpression)
	led(lexer.LEFT_BRACKET, CALL, parseIndexExpression)

	// ADDITIVE & MULTIPLICATIVE ---
	led(lexer.PLUS, ADDITIVE, parseBinaryExpression)
	led(lexer.MINUS, ADDITIVE, parseBinaryExpression)
	led(lexer.SLASH, MULTIPLICATIVE, parseBinaryExpression)
	led(lexer.MODULO, MULTIPLICATIVE, parseBinaryExpression)
	led(lexer.STAR, MULTIPLICATIVE, parseBinaryExpression)

	// STATEMENTS ---
	led(lexer.ASSIGNMENT, ASSIGNMENT, parseAssignmentExpression)
	statement(lexer.VARIABLE, parseVariableDeclaration)
	statement(lexer.CONSTANT, parseVariableDeclaration)
	statement(lexer.IF, parseIfStatement)
	statement(lexer.FUNCTION, parseFunctionDeclaration)
	statement(lexer.WHILE, parseWhileStatement)

	// CALL EXPRESSION ---
	led(lexer.LEFT_PARENTHESIS, CALL, parseCallExpression)

	// UNARY OPERATORS ---
	nud(lexer.NOT, parseUnaryExpression)
	nud(lexer.MINUS, parseUnaryExpression)
	nud(lexer.PLUS, parseUnaryExpression)

	// COMPARISON OPERATORS --+
	led(lexer.GREATER, RELATIONAL, parseBinaryExpression)
	led(lexer.GREATER_EQUAL, RELATIONAL, parseBinaryExpression)
	led(lexer.LESS, RELATIONAL, parseBinaryExpression)
	led(lexer.LESS_EQUAL, RELATIONAL, parseBinaryExpression)
	led(lexer.EQUAL, RELATIONAL, parseBinaryExpression)
	led(lexer.NOT_EQUAL, RELATIONAL, parseBinaryExpression)
}
