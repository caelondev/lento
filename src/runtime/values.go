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
	OBJECT_VALUE ValueTypes = "object"
	FUNCTION_VALUE        ValueTypes = "function"
	NATIVE_FUNCTION_VALUE ValueTypes = "native_function"
)

type RuntimeValue interface {
	Type() ValueTypes
	String() string
}

type ArrayValue struct {
	Elements []RuntimeValue
}

func (a *ArrayValue) Type() ValueTypes {
	return ARRAY_VALUE
}

func (a *ArrayValue) String() string {
	if len(a.Elements) == 0 {
		return "[]"
	}
	
	result := "["
	for i, elem := range a.Elements {
		if i > 0 {
			result += ", "
		}
		result += elem.String()
	}
	result += "]"
	return result
}

type ObjectPropertyValue struct {
	Key string
	Value RuntimeValue
}

type ObjectValue struct {
	Properties []ObjectPropertyValue
}

func (n *ObjectValue) Type() ValueTypes {
	return OBJECT_VALUE
}

func (n *ObjectValue) String() string {
	return n.stringWithIndent(0)
}

func (n *ObjectValue) stringWithIndent(depth int) string {
	if len(n.Properties) == 0 {
		return "{}"
	}

	indent := ""
	for range depth {
		indent += "  "
	}
	nextIndent := indent + "  "

	result := "{\n"
	for _, prop := range n.Properties {
		result += nextIndent + prop.Key + ": "
		
		// Handle nested objects with proper indentation
		if obj, ok := prop.Value.(*ObjectValue); ok {
			result += obj.stringWithIndent(depth + 1)
		} else {
			result += prop.Value.String()
		}
		
		result += ",\n"
	}
	result += indent + "}"
	return result
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

func OBJECT(properties []ObjectPropertyValue) *ObjectValue {
	return &ObjectValue{
		Properties: properties,
	}
}

func ARRAY(elements []RuntimeValue) *ArrayValue {
	return &ArrayValue{Elements: elements}
}
