package config

type Config struct {
	KubeConfPath string
	Definitions  []VariableDef
}

type VariableDef struct {
	Name   string
	Source VariableSource
	Type   VariableType
}
