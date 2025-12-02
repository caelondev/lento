package errorhandler

import (
	"fmt"
)

type ErrorHandler struct {
	HadError bool
	ErrorCode int
}

func New() *ErrorHandler {
	return &ErrorHandler{}
}

func (e *ErrorHandler) Report(line int, errorMessage string) {
	e.ReportError("Report", errorMessage, line, 65)
}

func (e *ErrorHandler) ReportError(reporter, errorMessage string, line, code int) {
	fmt.Printf("%s::Error on line %d: %s\n", reporter, int(line), errorMessage)
	e.ErrorCode = code
	e.HadError = true
}
