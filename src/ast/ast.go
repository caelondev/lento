package ast

type Statement interface {
	Statement()
	GetLine() uint
}

type Expression interface {
	Expression()
	GetLine() uint
}
