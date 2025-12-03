package runtime

import (
	"fmt"

	"github.com/caelondev/lento/src/ast"
)

type ValueTypes string

const (
	BOOLEAN_VALUE ValueTypes = "boolean"
	NIL_VALUE    ValueTypes = "nil"
	NUMBER_VALUE ValueTypes = "number"
	STRING_VALUE ValueTypes = "string"
	ARRAY_VALUE ValueTypes = "array"
	FUNCTION_VALUE        ValueTypes = "function"
	NATIVE_FUNCTION_VALUE ValueTypes = "native_function"
)

type RuntimeValue interface {
	Type() ValueTypes
	String() string
}

type NilValue struct{}

func (n *NilValue) Type() ValueTypes {
	return NIL_VALUE
}

func (n *NilValue) String() string {
	return "nil"
}

type NativeFunctionValue struct {
	Name string
	Call func(args []RuntimeValue, env Environment, i *Interpreter) RuntimeValue
}

func (n *NativeFunctionValue) Type() ValueTypes {
	return NATIVE_FUNCTION_VALUE
}

func (n *NativeFunctionValue) String() string {
	return fmt.Sprintf("[ native function '%s' ]", n.Name)
}


type StringValue struct {
	Value string
}

func (n *StringValue) Type() ValueTypes {
	return STRING_VALUE
}

func (n *StringValue) String() string {
	return fmt.Sprintf("\"%v\"", n.Value)
}

type NumberValue struct {
	Value float64
}

func (n *NumberValue) Type() ValueTypes {
	return NUMBER_VALUE
}

func (n *NumberValue) String() string {
	return fmt.Sprintf("%v", n.Value)
}

type FunctionValue struct {
	Name string
	Parameters []string
	Body ast.Statement
	Environment Environment
}

func (n *FunctionValue) Type() ValueTypes {
	return BOOLEAN_VALUE
}

func (n *FunctionValue) String() string {
	return fmt.Sprintf("[ function '%s' ]", n.Name)
}

type BooleanValue struct {
	Value bool
}

func (n *BooleanValue) Type() ValueTypes {
	return BOOLEAN_VALUE
}

func (n *BooleanValue) String() string {
	return fmt.Sprintf("%v", n.Value)
}

func NIL() *NilValue {
	return &NilValue{}
}

func BOOLEAN(value bool) *BooleanValue {
	return &BooleanValue{ Value: value }
}

func NATIVE_FUNCTION(name string, call func([]RuntimeValue, Environment, *Interpreter) RuntimeValue) *NativeFunctionValue {
	return &NativeFunctionValue{
		Name: name,
		Call: call,
	}
}

