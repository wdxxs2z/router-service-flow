package policy

import (
	"github.com/wdxxs2z/router-service-flow/models"
	"fmt"
	"hash/fnv"
)

type SourceHash struct{
	PolicyType
	remoteHost string
}

func NewSourceHash(typeName string, nodes []models.Node, remoteIp string) *SourceHash{
	return &SourceHash{PolicyType{typeName, nodes},remoteIp}
}

func (sourceHash *SourceHash) WinUrl() string {
	nodeLength := len(sourceHash.PolicyType.Nodes)

	if nodeLength <= 0 {
		fmt.Printf("Fetch node 0.\n")
		return ""
	}

	hashcode := hash(sourceHash.remoteHost)
	pos      := hashcode % nodeLength

	serverNode := sourceHash.Nodes[pos]

	return serverNode.Url
}

func hash(ip string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(ip))
	return h.Sum32()
}
