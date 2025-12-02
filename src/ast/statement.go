package ast

type BlockStatement struct {
	Body []Statement
}

func (node *BlockStatement) Statement() {}

type ExpressionStatement struct {
	Expression Expression
}

func (node *ExpressionStatement) Statement() {}
