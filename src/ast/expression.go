package ast

import "github.com/caelondev/lento/src/lexer"


type NumberExpression struct {
	Value float64
}

func (node *NumberExpression) Expression() {}

type StringExpression struct {
	Value string
}

func (node *StringExpression) Expression() {}

type SymbolExpression struct {
	Value string
}

func (node *SymbolExpression) Expression() {}


type BinaryExpression struct {
	Left Expression
	Right Expression
	Operator *lexer.Token
}

func (node *BinaryExpression) Expression() {}
