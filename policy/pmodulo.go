package policy
import (
	"github.com/wdxxs2z/router-service-flow/models"
	"math/rand"
	"fmt"
	"time"
)

const DEFAULT_NUMBER  = 1000

type Modulo struct {
	PolicyType
}

func NewModulo(typeName string, nodes [] models.Node) *Modulo{
	return &Modulo{PolicyType{typeName,nodes}}
}

func (modulo *Modulo) WinUrl() string {
	sum := int64(0)
	for _,node := range modulo.Nodes {
		sum += node.Weight
	}
	randNumber := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(DEFAULT_NUMBER)
	mod := randNumber % sum
	fmt.Println(mod)
	for index,node := range modulo.Nodes {
		weight := node.Weight
		url    := node.Url
		// case 1 b > a
		if weight > (sum - weight) {
			if mod >= weight && mod == 0{
				return url
			} else {
				return modulo.Nodes[index +1].Url
			}
		}
		// case 2 b < a
		if weight < (sum - weight) {
			if mod < (sum - weight) && mod != 0{
				return url
			} else {
				return modulo.Nodes[index+1].Url
			}
		}
	}
	return ""
}