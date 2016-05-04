package models

import (
	"math/rand"
	"container/list"
	"time"
	"sync"
	"reflect"
	log "github.com/Sirupsen/logrus"
)

const DEFAULT_NUMBER  = 1000

type Normal struct  {
	Nodes      list.List
	randNumber int64
	Resources map[int]bool
	sync.RWMutex
}

func NewNormal() *Normal {
	return &Normal{
		Nodes:  list.New(),
		Resources: make(map[int]bool),
		randNumber: rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(DEFAULT_NUMBER),
	}
}

func (n *Normal) AddNode(node Node) bool {
	n.Lock()
	defer n.Unlock()

	if _,ok := n.Resources[node.Index];ok {
		return false
	}

	n.Nodes.PushBack(node)
	n.Resources[node.Index] = true
	return true
}

func (n *Normal) GetWinUrl() string {
	sum := int64(0)
	for e := n.Nodes.Front(); e != nil; e = e.Next(){
		weight := int64(reflect.ValueOf(&e.Value).Elem().FieldByName("Weight"))
		sum += weight
	}
	mod := n.randNumber % sum
	log.Println("The mod is : %v", mod)
	for e := n.Nodes.Front(); e!= nil; e = e.Next() {
		weight := int64(reflect.ValueOf(&e.Value).Elem().FieldByName("Weight"))
		url    := reflect.ValueOf(&e.Value).Elem().FieldByName("Url")
		if weight >= mod || mod == 0 {
			return url
		}
		if weight < mod {
			return url
		}
		if weight == 0 {
			continue
		}
	}
	return ""
}
