package models

import (
	"hash/crc32"
	"sync"
	"strconv"
	"sort"
)

//set default replicas 120
const DEFAULT_REPLICAS = 160

type HashRing []uint32

func (c HashRing) Len() int {
	return len(c)
}

func (c HashRing) Less(i, j int) bool {
	return c[i] < c[j]
}

func (c HashRing) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type Consistent struct  {
	Nodes    map[uint32]Node
	numReps   int
	Resources map[int]bool
	ring      HashRing
	sync.RWMutex
}

func NewConsistent() *Consistent {
	return &Consistent{
		Nodes:     make(map[uint32]Node),
		numReps:   DEFAULT_REPLICAS,
		Resources: make(map[int]bool),
		ring:      HashRing{},
	}
}

func (c *Consistent) Add(node *Node) bool {
	c.Lock()
	defer c.Unlock()

	if _,ok := c.Resources[node.Index];ok {
		return false
	}

	count := c.numReps * node.Weight
	//add the ring
	for i := 0;i < count;i++ {
		str := c.joinStr(i, node)
		c.Nodes[c.hashStr(str)] = *(node)
	}
	c.Resources[node.Index] = true
	c.sortHashRing()
	return true
}

func (c *Consistent)joinStr(countNum int, node *Node) string{
	return node.Url + "*" + strconv.Itoa(node.Weight) + "-" + strconv.Itoa(countNum) + "-" + strconv.Itoa(node.Index)
}

func (c *Consistent)hashStr(key string) uint32{
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *Consistent)sortHashRing() {
	c.ring = HashRing{}
	for key := range c.Nodes {
		c.ring = append(c.ring, key)
	}
	sort.Sort(c.ring)
}

func (c * Consistent)Get(key string) Node {
	c.Lock()
	defer c.Unlock()

	hash := c.hashStr(key)
	i := c.search(hash)
	return c.Nodes[c.ring[i]]
}

func (c *Consistent) search(hash uint32) int {
	i := sort.Search(len(c.ring),func(i int) bool {return c.ring[i] >= hash})
	if i < len(c.ring) {
		if i == len(c.ring)-1{
			return 0
		}else{
			return i
		}
	}else{
		return len(c.ring)-1
	}
}

func (c *Consistent) Remove(node *Node) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Resources[node.Index]; !ok {
		return
	}

	delete(c.Resources, node.Index)

	count := c.numReps * node.Weight
	for i := 0; i < count; i++ {
		str := c.joinStr(i, node)
		delete(c.Nodes, c.hashStr(str))
	}
	c.sortHashRing()
}