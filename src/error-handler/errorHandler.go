package errorhandler

import (
	"fmt"
	"os"
)

type ErrorHandler struct {
	hadError bool
}

func New() *ErrorHandler {
	return &ErrorHandler{}
}

func (e *ErrorHandler) Report(line int, errorMessage string) {
	e.ReportError("Report", errorMessage, line, 65)
}

func (e *ErrorHandler) ReportError(reporter, errorMessage string, line, code int) {
	fmt.Printf("%s::Error on line %d: %s\n", reporter, int(line), errorMessage)
	os.Exit(code)
}
