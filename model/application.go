package model

type Application struct {
	Name        string
	Namespace   string
	Definitions []VariableDef
}
