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

type IfStatement struct {
    Condition  Expression
    Consequent Statement
    Alternate  Statement
    Line       uint
}

func (i *IfStatement) GetLine() uint {
    return i.Line
}

func (i *IfStatement) Statement() {}

type FunctionDeclarationStatement struct {
    Name       string
    Parameters []string
    Body       Statement
    Line       uint
}

func (f *FunctionDeclarationStatement) Statement() {}
func (f *FunctionDeclarationStatement) GetLine() uint {
    return f.Line
}

