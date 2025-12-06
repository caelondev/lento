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
	Init      Statement
	Condition Expression
	Increment Expression
	Body      Statement
	Line      uint
}

func (f *ForStatement) Statement() {}
func (f *ForStatement) GetLine() uint {
	return f.Line
}

type ReturnStatement struct {
	Value Expression
	Line  uint
}

func (r *ReturnStatement) GetLine() uint { return r.Line }
func (r *ReturnStatement) Statement() {}

type BreakStatement struct {
	Line uint
}

func (b *BreakStatement) GetLine() uint { return b.Line }
func (b *BreakStatement) Statement() {}

type ContinueStatement struct {
	Line uint
}

func (c *ContinueStatement) GetLine() uint { return c.Line }
func (c *ContinueStatement) Statement() {}
