package ast

import "github.com/caelondev/lento/src/lexer"

type NumberExpression struct {
	Value float64
	Line  uint
}

func (node *NumberExpression) Expression() {}
func (node *NumberExpression) GetLine() uint {
	return node.Line
}

type StringExpression struct {
	Value string
	Line  uint
}

func (node *StringExpression) Expression() {}
func (node *StringExpression) GetLine() uint {
	return node.Line
}

type SymbolExpression struct {
	Value string
	Line  uint
}

func (node *SymbolExpression) Expression() {}
func (node *SymbolExpression) GetLine() uint {
	return node.Line
}

type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator *lexer.Token
	Line     uint
}

func (node *BinaryExpression) Expression() {}
func (node *BinaryExpression) GetLine() uint {
	return node.Line
}

type UnaryExpression struct {
	Operator *lexer.Token
	Operand Expression
	Line uint
}

func (node *UnaryExpression) Expression() {}
func (node *UnaryExpression) GetLine() uint {
	return node.Line
}
