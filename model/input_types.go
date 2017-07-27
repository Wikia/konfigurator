package model

type InputType string

const (
	SIMPLE         InputType = "simple"
	VAULT                    = "vault"
	CONSUL                   = "consul"
	LAYERED_CONSUL           = "layered_consul"
)
