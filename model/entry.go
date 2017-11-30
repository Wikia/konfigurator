package model

type VariableType string

const (
	STANDARD   VariableType = "standard"
	INLINE                  = "inline"
	REFERENCED              = "referenced"
)

type Variable struct {
	Name        string
	Destination VariableDestination
	Type        VariableType
	Value       interface{}
}
