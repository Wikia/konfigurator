package model

type VariableType string

const (
	CONFIGMAP VariableType = "config"
	SECRET                 = "secret"
	REFERENCE              = "reference"
)

type VariableDef struct {
	Name    string
	Source  InputType
	Type    VariableType
	Value   interface{}
	Context map[string]string
}

func NewVariableDef() VariableDef {
	return VariableDef{Context: map[string]string{}}
}
