package runtime

import (
	"fmt"
	"strconv"
	"strings"

	errorhandler "github.com/caelondev/lento/src/error-handler"
)

func NATIVE_PRINT_FUNCTION(args []RuntimeValue, env Environment, i *Interpreter) RuntimeValue {
	for _, arg := range args {
		if arg.Type() == "string" {
			strValue := arg.String()[1 : len(arg.String())-1]
			fmt.Print(strValue)
		} else {
			fmt.Print(arg)
		}
	}
	fmt.Println()
	return NIL()
}

func NATIVE_LEN_FUNCTION(args []RuntimeValue, env Environment, i *Interpreter) RuntimeValue {
	if len(args) != 1 {
		i.errorHandler.ReportError(
			"Interpreter-Native-Function",
			"len() expects exactly one argument",
			i.line,
			errorhandler.ArgumentLengthError,
		)
		return NIL()
	}

	arg := args[0]
	switch arg.Type() {
	case STRING_VALUE:
		strValue := arg.String()[1 : len(arg.String())-1]
		return &NumberValue{Value: float64(len(strValue))}
	case ARRAY_VALUE:
		arr, _ := arg.(*ArrayValue)
		return &NumberValue{Value: float64(len(arr.Elements))}

	default:
		i.errorHandler.ReportError(
			"Interpreter-Native-Function",
			fmt.Sprintf("Could not use len() on unsupported argument type (%s)", arg.Type()),
			i.line,
			errorhandler.ArgumentLengthError,
		)
		return NIL()
	}
}

func NATIVE_TO_UPPER_FUNCTION(args []RuntimeValue, env Environment, i *Interpreter) RuntimeValue {
	if len(args) != 1 || args[0].Type() != "string" {
		i.errorHandler.ReportError("Interpreter-Native-Function", "toUpper() expects one string argument", i.line, errorhandler.ArgumentLengthError)
		return NIL()
	}
	strValue := args[0].String()[1 : len(args[0].String())-1]
	return &StringValue{Value: strings.ToUpper(strValue)}
}

func NATIVE_TO_LOWER_FUNCTION(args []RuntimeValue, env Environment, i *Interpreter) RuntimeValue {
	if len(args) != 1 || args[0].Type() != "string" {
		i.errorHandler.ReportError("Interpreter-Native-Function", "toLower() expects one string argument", i.line, errorhandler.ArgumentLengthError)
		return NIL()
	}
	strValue := args[0].String()[1 : len(args[0].String())-1]
	return &StringValue{Value: strings.ToLower(strValue)}
}

func NATIVE_STR_FUNCTION(args []RuntimeValue, env Environment, i *Interpreter) RuntimeValue {
	if len(args) != 1 {
		i.errorHandler.ReportError("Interpreter-Native-Function", "str() expects exactly one argument", i.line, errorhandler.ArgumentLengthError)
		return NIL()
	}
	val := args[0]
	if val.Type() == "string" {
		return val
	}
	return &StringValue{Value: val.String()}
}

func NATIVE_NUM_FUNCTION(args []RuntimeValue, env Environment, i *Interpreter) RuntimeValue {
	if len(args) != 1 {
		i.errorHandler.ReportError("Interpreter-Native-Function", "num() expects exactly one argument", i.line, errorhandler.ArgumentLengthError)
		return NIL()
	}

	val := args[0]
	switch val.Type() {
	case "number":
		n, _ := strconv.ParseFloat(val.String(), 64)
		return &NumberValue{Value: n}
	case "string":
		str := val.String()[1 : len(val.String())-1]
		num, err := strconv.Atoi(str)
		if err != nil {
			i.errorHandler.ReportError("Interpreter-Native-Function", "num() invalid string to convert", i.line, errorhandler.NativeFunctionError)
			return NIL()
		}
		return &NumberValue{Value: float64(num)}
	default:
		i.errorHandler.ReportError("Interpreter-Native-Function", "num() can only convert number or string", i.line, "ERR_INT_TYPE")
		return NIL()
	}
}
