package config

type VariableType uint

const (
	_                      = iota
	CONFIGMAP VariableType = 1
	SECRET
)
