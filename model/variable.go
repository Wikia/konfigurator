package model

type VariableType string

const (
	CONFIGMAP VariableType = "config"
	SECRET                 = "secret"
)
