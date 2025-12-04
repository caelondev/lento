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
	Operand  Expression
	Line     uint
}

func (node *UnaryExpression) Expression() {}
func (node *UnaryExpression) GetLine() uint {
	return node.Line
}

type AssignmentExpression struct {
	Assignee Expression
	Value    Expression
	Line     uint
}

func (node *AssignmentExpression) Expression() {}
func (node *AssignmentExpression) GetLine() uint {
	return node.Line
}

type CallExpression struct {
	Caller    Expression
	Arguments []Expression
	Line      uint
}

func (c *CallExpression) Expression() {}
func (c *CallExpression) GetLine() uint {
	return c.Line
}

type ArrayExpression struct {
	Elements []Expression
	Line     uint
}

func (a *ArrayExpression) Expression() {}
func (a *ArrayExpression) GetLine() uint {
	return a.Line
}

type IndexExpression struct {
	Array Expression
	Index  Expression
	Line   uint
}

func (i *IndexExpression) Expression() {}
func (i *IndexExpression) GetLine() uint {
	return i.Line
}

type ObjectProperty struct {
	Key string
	Value Expression
}

type ObjectExpression struct {
	Properties []ObjectProperty
	Line uint
}

func (i *ObjectExpression) Expression() {}
func (i *ObjectExpression) GetLine() uint {
	return i.Line
}


