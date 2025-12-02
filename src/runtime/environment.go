package runtime

import (
	"fmt"
	"slices"
	errorhandler "github.com/caelondev/lento/src/error-handler"
)

type Environment interface {
	DeclareVariable(line uint, variableName string, value RuntimeValue, isConstant bool, isNative bool) RuntimeValue
	AssignVariable(line uint, variableName string, value RuntimeValue) RuntimeValue
	LookupVariable(line uint, variableName string) RuntimeValue
	ResolveVariable(line uint, variableName string) Environment
	IsNative(variableName string) bool
}

type EnvironmentStruct struct {
	parent       Environment
	variables    map[string]RuntimeValue
	constants    []string
	natives      []string
	errorHandler *errorhandler.ErrorHandler
}

func NewEnvironment(parent Environment, errorHandler *errorhandler.ErrorHandler) Environment {
	env := &EnvironmentStruct{
		parent:       parent,
		variables:    make(map[string]RuntimeValue),
		constants:    make([]string, 0),
		natives:      make([]string, 0),
		errorHandler: errorHandler,
	}

	if parent == nil { // Global Env ---
		DeclareGlobalVariables(env)
	}
	return env
}

func DeclareGlobalVariables(env Environment) {
	isConstant := true
	isNative := true

	env.DeclareVariable(0, "nil", NIL(), isConstant, isNative)
	env.DeclareVariable(0, "true", BOOLEAN(true), isConstant, isNative)
	env.DeclareVariable(0, "false", BOOLEAN(false), isConstant, isNative)
}

func (e *EnvironmentStruct) DeclareVariable(line uint, variableName string, value RuntimeValue, isConstant bool, isNative bool) RuntimeValue {
	_, exists := e.variables[variableName]

	if exists {
		e.errorHandler.Report(
			int(line),
			fmt.Sprintf("Cannot declare variable '%s' as it already exists", variableName),
		)
		return NIL()
	}

	if isConstant {
		e.constants = append(e.constants, variableName)
	}

	if isNative {
		e.natives = append(e.natives, variableName)
	}

	e.variables[variableName] = value

	return value
}

func (e *EnvironmentStruct) AssignVariable(line uint, variableName string, value RuntimeValue) RuntimeValue {
	env := e.ResolveVariable(line, variableName)

	if env == nil {
		e.errorHandler.Report(int(line), fmt.Sprintf(
			"Invalid left-hand assignment: '%s' is not defined",
			variableName,
		))
		return NIL()
	}

	envStruct := env.(*EnvironmentStruct)

	isNative := slices.Contains(envStruct.natives, variableName)
	isConstant := slices.Contains(envStruct.constants, variableName)

	if isNative {
		e.errorHandler.Report(int(line), fmt.Sprintf(
			"Cannot reassign keyword '%s'",
			variableName,
		))
		return NIL()
	}

	if isConstant {
		e.errorHandler.Report(int(line), fmt.Sprintf(
			"Cannot reassign constant '%s'",
			variableName,
		))
		return NIL()
	}

	envStruct.variables[variableName] = value
	return value
}

func (e *EnvironmentStruct) LookupVariable(line uint, variableName string) RuntimeValue {
	env := e.ResolveVariable(line, variableName)

	if env == nil {
		e.errorHandler.Report(int(line), fmt.Sprintf(
			"Cannot resolve variable '%s'",
			variableName,
		))
		return NIL()
	}

	envStruct := env.(*EnvironmentStruct)
	return envStruct.variables[variableName]
}

func (e *EnvironmentStruct) ResolveVariable(line uint, variableName string) Environment {
	if _, exists := e.variables[variableName]; exists {
		return e
	}

	if e.parent == nil {
		return nil
	}

	return e.parent.ResolveVariable(line, variableName)
}

func (e *EnvironmentStruct) IsNative(variableName string) bool {
	if slices.Contains(e.natives, variableName) {
		return true
	}

	if e.parent != nil {
		return e.parent.IsNative(variableName)
	}

	return false
}
