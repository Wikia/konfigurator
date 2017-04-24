package model

type VariableType string

const (
	CONFIGMAP VariableType = "config"
	SECRET                 = "secret"
)

type VariableDef struct {
	Name   string
	Source InputType
	Type   VariableType
	Value  interface{}
}
