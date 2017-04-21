package model

type InputType string

const (
	SIMPLE  InputType = "simple"
	POD               = "pod"
	VAULT             = "valut"
	FILE              = "file"
	PODYAML           = "pod_yaml"
)

type VariableDef struct {
	Name   string
	Source InputType
	Type   VariableType
	Value  interface{}
}
