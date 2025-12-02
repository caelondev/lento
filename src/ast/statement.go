package ast

type BlockStatement struct {
	Body []Statement
	Line uint
}

func (node *BlockStatement) Statement() {}
func (node *BlockStatement) GetLine() uint {
	return node.Line
}

type ExpressionStatement struct {
	Expression Expression
	Line       uint
}

func (node *ExpressionStatement) Statement() {}
func (node *ExpressionStatement) GetLine() uint {
	return node.Line
}

type VariableDeclarationStatement struct {
	IsConstant bool
	Identifier string
	Value Expression
	Line       uint
}

func (node *VariableDeclarationStatement) Statement() {}
func (node *VariableDeclarationStatement) GetLine() uint {
	return node.Line
}
