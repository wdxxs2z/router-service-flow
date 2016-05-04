package models

type Node struct {
	Index     int
	Url       string
	Weight    int
}

func NewNode(index int, url string, weight int)  *Node {
	return &Node{
		Index:   index,
		Url:     url,
		Weight : weight,
	}
}

