package model

type VariableType string

type Variable struct {
	Name        string
	Destination VariableDestination
	Source      InputType
	Value       interface{}
}
