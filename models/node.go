package models

type Node struct {
	Index     int64       	`json:"index"`
	Url       string    	`json:"url"`
	Weight    int64           `json:"weight"`
}

