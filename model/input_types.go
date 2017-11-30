package model

type InputType string

const (
	SIMPLE         InputType = "simple"
	REFERENCE                = "reference"
	VAULT                    = "vault"
	CONSUL                   = "consul"
	LAYERED_CONSUL           = "layered_consul"
)
