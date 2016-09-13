package policy

import (
	"fmt"
	"github.com/wdxxs2z/router-service-flow/models"
	"sync"
)

var (
	pos = 0
	mutex sync.Mutex
)

type RoundRobin struct{
	PolicyType
}

func NewRoundRobin(typeName string, nodes []models.Node) *RoundRobin{
	return &RoundRobin{PolicyType{typeName, nodes}}
}

func (roundRobin *RoundRobin) WinUrl() string {

	nodeLength := len(roundRobin.PolicyType.Nodes)

	if nodeLength <= 0 {
		fmt.Printf("Fetch node 0.\n")
		return ""
	}

	mutex.Lock()
	if pos >= nodeLength {
		pos = 0
	}
	node := roundRobin.Nodes[pos]
	winUrl := node.Url
	pos++
	mutex.Unlock()

	return winUrl
}
