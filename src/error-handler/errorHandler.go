package errorhandler

import (
	"fmt"
)

type ErrorHandler struct {
	HadError bool
	ErrorMessage ErrorType
}

func New() *ErrorHandler {
	return &ErrorHandler{}
}

func (e *ErrorHandler) Report(line uint, errorMessage string) {
	e.ReportError("Report", errorMessage, line, ReportError)
}

func (e *ErrorHandler) ReportError(reporter, errorMessage string, line uint, code ErrorType) {
	if e.HadError { // Avoid multiple confusing error messages ---
		return
	}
	fmt.Printf("%s::Error on line %d: %s\n", reporter, int(line), errorMessage)
	e.HadError = true
}
