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
	Nodes      *list.List
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

func (n *Normal) AddNode(node *Node) bool {
	n.Lock()
	defer n.Unlock()

	if _,ok := n.Resources[node.Index];ok {
		return false
	}
	nodes := n.Nodes
	nodes.PushBack(node)
	n.Resources[node.Index] = true
	return true
}

func (n *Normal) GetWinUrl() string {
	n.Lock()
	defer n.Unlock()

	sum := int64(0)
	nodes := n.Nodes
	for e := nodes.Front(); e != nil; e = e.Next(){
		weight := reflect.ValueOf(e.Value).Elem().FieldByName("Weight").Int()
		sum += weight
	}
	mod := n.randNumber % sum
	log.Println("The mod is : %v", mod)
	for e := n.Nodes.Front(); e!= nil; e = e.Next() {
		weight := reflect.ValueOf(e.Value).Elem().FieldByName("Weight").Int()
		url    := reflect.ValueOf(e.Value).Elem().FieldByName("Url").String()
		// case 1 b > a
		if weight > (sum - weight) {
			if mod >= weight && mod == 0{
				return url
			} else {
				return reflect.ValueOf(e.Next().Value).Elem().FieldByName("Url").String()
			}
		}
		// case 2 b < a
		if weight < (sum - weight) {
			if mod < (sum - weight) && mod != 0{
				return url
			} else {
				return reflect.ValueOf(e.Next().Value).Elem().FieldByName("Url").String()
			}
		}
	}
	return ""
}