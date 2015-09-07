package expression

// Context represents an execution context
// for a set of expressions/statements.
// It stores variable and function definitions.
type Context struct {
	values map[string]Expression
}

const (
	constPrefix  = "c"
	varPrefix    = "v"
	prefixLength = 1
)

// NewContext creates a new empty Context.
func NewContext() Context {
	return Context{make(map[string]Expression)}
}

// GetVariables returns the variables contained by a Context.
func (c Context) GetVariables() map[string]Expression {
	vars := make(map[string]Expression)
	for k, v := range c.values {
		if k[:prefixLength] == varPrefix {
			vars[k[prefixLength:]] = v
		}
	}
	return vars
}

// GetConstants returns the constants contained by a Context.
func (c Context) GetConstants() map[string]Expression {
	consts := make(map[string]Expression)
	for k, v := range c.values {
		if k[:prefixLength] == constPrefix {
			consts[k[prefixLength:]] = v
		}
	}
	return consts
}

// SetVariable sets the value of a variable.
// If a constant with the specified name already exists,
// an error is returned.
func (c Context) SetVariable(name string, Expression value) error {
	if c.HasConstant(name) {
		return errors.New("Cannot change the value of a constant.")
	}
	c.values[varPrefix+name] = value
	return nil
}

// SetConstant creates a new constant. It can never be deleted or changed.
// If a constant or variable with the specified name already exists,
// an error is returned.
func (c Context) SetConstant(name string, Expression value) error {
	if c.HasConstant(name) {
		return errors.New("A constant with the name " + name + " already exists.")
	} else if c.HasVariable(name) {
		return errors.New(name + " is a variable. It cannot be turned into a constant")
	}
	c.values[constPrefix+name] = value
	return nil
}

// RemoveVariable removes a variable with the specified name.
// If no such variable exists, an error is returned.
// If the specified name refers to a constant, an error is returned.
func (c Context) RemoveVariable(name string) error {
	if c.HasVariable(name) {
		delete(c.values, varPrefix+name)
		return nil
	} else if c.HasConstant(name) {
		return errors.New(name + " is a constant. It cannot be deleted.")
	}
	return errors.New("No variable " + name + ".")
}

// Has checks whether a Context contains a variable or constant
// with the specified name.
func (c Context) Has(name string) bool {
	return c.HasVariable(name) || c.HasConstant(name)
}

// HasVariable checks whether a Context has a variable with
// the specified name.
func (c Context) HasVariable(name string) bool {
	_, ok := c.values[varPrefix+name]
	return ok
}

// HasConstant checks whether a Context has a constant with
// the specified name.
func (c Context) HasConstant(name string) bool {
	_, ok := c.values[constPrefix+name]
	return ok
}

// Get returns the the value of a variable or constant
// with the specified name. If no such variable or constant
// exists, it returns an error.
func (c Context) Get(name string) (Expression, error) {
	if c.HasVariable(name) {
		return c.values[varPrefix+name], nil
	} else if c.HasConstant(name) {
		return c.values[constPrefix+name], nil
	}
	return nil, errors.New("No variable or constant named " + name + ".")
}
