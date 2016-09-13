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

type RoundRobinWeght struct{
	PolicyType
}

func NewRoundRobinWeight(typeName string, nodes []models.Node) *RoundRobinWeght{
	return &RoundRobinWeght{PolicyType{typeName, nodes}}
}

func (roundRobinWeght *RoundRobinWeght) WinUrl() string {
	nodeLength := len(roundRobinWeght.PolicyType.Nodes)

	if nodeLength <= 0 {
		fmt.Printf("Fetch node 0.\n")
		return ""
	}

	var serverList []string
	listLen := 0

	for _,node := range RoundRobinWeght.Nodes {
		nodeWeight := node.Weight
		for i:=0;i<nodeWeight;i++ {
			serverList[listLen] = node.Url
			listLen ++
		}
	}

	mutex.Lock()
	if pos >= nodeLength {
		pos = 0
	}
	winUrl := serverList[pos]
	pos++
	mutex.Unlock()
	return winUrl
}
