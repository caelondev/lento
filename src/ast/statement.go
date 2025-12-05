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
	Value      Expression
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

type WhileLoopStatement struct {
	Condition Expression
	Body      Statement
	Line      uint
}

func (w *WhileLoopStatement) Statement() {}
func (w *WhileLoopStatement) GetLine() uint {
	return w.Line
}

type ForStatement struct {
	Init      Statement  // Initialization (e.g., var x = 0)
	Condition Expression // Loop condition (e.g., x < 10)
	Increment Expression  // Increment statement (e.g., x++)
	Body      Statement  // Usually a BlockStatement
	Line      uint
}

func (f *ForStatement) Statement() {}
func (f *ForStatement) GetLine() uint {
	return f.Line
}
