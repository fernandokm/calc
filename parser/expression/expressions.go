package expression

// Expression represents an evaluatable expression.
//
// grammar: expression
type Expression interface {
	// Evaluate returns an interface{} and an error.
	// If err is non-nil, value must be nil.
	Evaluate(Context) (value interface{}, err error)
}

// EvaluateNumber calls e.Evaluate(c). If an error occurs,
// value is set to nil and err is set to that error.
// Otherwise, EvaluateNumber attempts to convert the result
// of the expression to a *big.Rat. Upon failure, it returns
// value=nil and an error indicating the unexpected type. Upon success,
// the correct value is return with a err=nil.
func EvaluateNumber(e Expression, c Context) (value *big.Rat, err error) {
	val, err := e.Evaluate(c)
	if err != nil {
		return nil, err
	}
	n, ok := val.(*big.Rat)
	if !ok {
		return nil, UnexpectedType("*big.Rat", fmt.Sprintf("%T", val))
	}
	return n, nil
}

// EvaluateBoolean calls e.Evaluate(c). If an error occurs,
// value is set to false and err is set to that error.
// Otherwise, EvaluateBool attempts to convert the result
// of the expression to a bool. Upon failure, it returns
// value=false and an error indicating the unexpected type. Upon success,
// the correct value is return with a err=nil.
func EvaluateBoolean(e Expression, c Context) (value bool, err error) {
	val, err := e.Evaluate(c)
	if err != nil {
		return false, err
	}
	b, ok := val.(bool)
	if !ok {
		return false, UnexpectedType("bool", fmt.Sprintf("%T", val))
	}
	return b, nil
}

// EvaluateSet calls e.Evaluate(c). If an error occurs,
// value is set to nil and err is set to that error.
// Otherwise, EvaluateSet attempts to convert the result
// of the expression to a []Expression. Upon failure, it returns
// value=nil and an error indicating the unexpected type. Upon success,
// the correct value is return with a err=nil.
func EvaluateSet(e Expression, c Context) ([]Expression, error) {
	val, err := e.Evaluate(c)
	if err != nil {
		return nil, err
	}
	n, ok := val.([]Expression)
	if !ok {
		return nil, UnexpectedType("[]Expression", fmt.Sprintf("%T", val))
	}
	return n
}
