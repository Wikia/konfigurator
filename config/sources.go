package config

type VariableSource uint

const (
	_                     = iota
	SIMPLE VariableSource = 1
	POD
	VAULT
	FAILE
	PODYAML
)
