package policy

import (
	"github.com/wdxxs2z/router-service-flow/models"
)

const (
	POLICY_MODULO        = "modulo"
	POLICY_ROUNDROBIN    = "roundrobin"
)

type PolicyType struct {
	TypeName string			`json:"typename"`
	Nodes    [] models.Node         `json:"nodes"`
}

type Policy interface {
	winUrl() string
}
