package models

import (
	"github.com/sentinel-official/sentinel-go-sdk/libs/geoip"
)

type Node struct {
	Addr      string          `json:"addr,omitempty" bson:"addr"`
	Error     string          `json:"error,omitempty" bson:"error"`
	Location  *geoip.Location `json:"location,omitempty" bson:"location"`
	Timestamp int64           `json:"timestamp,omitempty" bson:"timestamp"`
}

func (n *Node) IsZero() bool {
	return n.Addr == ""
}
