package policy

import (
	"fmt"
	"github.com/wdxxs2z/router-service-flow/models"
	"sync"
)

var (
	round_weight_pos = 0
	round_weight_mutex sync.Mutex
)

type RoundRobinWeight struct{
	PolicyType
}

func NewRoundRobinWeight(typeName string, nodes []models.Node) *RoundRobinWeight{
	return &RoundRobinWeight{PolicyType{typeName, nodes}}
}

func (roundRobinWeight *RoundRobinWeight) WinUrl() string {
	nodeLength := len(roundRobinWeight.PolicyType.Nodes)

	if nodeLength <= 0 {
		fmt.Printf("Fetch node 0.\n")
		return ""
	}

	var serverList []string
	listLen := 0

	for _,node := range roundRobinWeight.Nodes {
		nodeWeight := node.Weight
		for i := int64(0); i < nodeWeight; i++ {
			serverList[listLen] = node.Url
			listLen ++
		}
	}

	round_weight_mutex.Lock()
	if round_weight_pos >= nodeLength {
		round_weight_pos = 0
	}
	winUrl := serverList[round_weight_pos]
	round_weight_pos ++
	round_weight_mutex.Unlock()
	return winUrl
}
