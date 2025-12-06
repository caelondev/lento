package errorhandler

type ErrorType string

const (
	ArgumentLengthError ErrorType = "ENGTH_ARGS_ERR"
	ReportError ErrorType = "REPORT_ERR"
	ExpectedTypeError ErrorType = "EXPECTED_TYPE_ERR"
	UnterminatedError ErrorType = "UNTERMINATED_ERR"
	UnknownTokenError ErrorType = "UNKNOWN_TOKEN_ERR"
	UnexpectedTokenError ErrorType = "UNEXPEXTED_TOKEN_ERR"
	InvalidArgumentError ErrorType = "INVALID_ARGUMENT_ERR"
	NonFunctionExpressionError ErrorType = "NON_FUNCTION_EXPR_ERR"
	NativeFunctionError ErrorType = "NATIVE_FUNCTION_ERR"
	ArrayIndexError ErrorType = "ARRAY_INDEXING_ERR"
	ObjectKeyError ErrorType = "OBJECT_KEY_ERR"
	MemberExpressionError ErrorType = "MEMBER_EXPR_ERR"
	VariableDeclarationError ErrorType = "VARIABLE_DECL_ERR"
	InvalidPostfixExpressionError ErrorType = "INVALID_POSTFIX_EXPR_ERR"
	IllegalStatementError ErrorType = "ILLEGAL_STATEMENT_ERR"
)
