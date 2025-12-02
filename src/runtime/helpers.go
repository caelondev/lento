package runtime

func isTruthy(value RuntimeValue) bool {
	switch v := value.(type) {
	case *NilValue:
		return false
	case *BooleanValue:
		return v.Value
	case *NumberValue:
		return v.Value != 0
	case *StringValue:
		return v.Value != ""
	default:
		return true
	}
}

