package models

import (
	"github.com/sentinel-official/sentinel-go-sdk/libs/geoip"
)

type Node struct {
	Addr      string          `bson:"addr"      json:"addr,omitempty"`
	Error     string          `bson:"error"     json:"error,omitempty"`
	Location  *geoip.Location `bson:"location"  json:"location,omitempty"`
	Timestamp int64           `bson:"timestamp" json:"timestamp,omitempty"`
}

func (n *Node) IsZero() bool {
	return n.Addr == ""
}
