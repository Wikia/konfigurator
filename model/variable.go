package model

type VariableDestination string

const (
	CONFIGMAP VariableDestination = "config"
	SECRET    VariableDestination = "secret"
	INLINE    VariableDestination = "inline"
)

type VariableDef struct {
	Name        string
	Source      InputType
	Destination VariableDestination
	Value       interface{}
	Context     map[string]string
}

func NewVariableDef() VariableDef {
	return VariableDef{Context: map[string]string{}}
}
