package node

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sentinel-official/sentinelhub/v12/types"
)

type RequestGetNode struct {
	URI struct {
		Addr string `uri:"addr"`
	}
}

func NewRequestGetNode(ctx *gin.Context) (req *RequestGetNode, err error) {
	req = &RequestGetNode{}
	if err = ctx.ShouldBindUri(&req.URI); err != nil {
		return nil, fmt.Errorf("failed to bind uri: %w", err)
	}

	return req, nil
}

type RequestGetNodes struct {
}

func NewRequestGetNodes(_ *gin.Context) (req *RequestGetNodes, err error) {
	req = &RequestGetNodes{}

	return req, nil
}

type RequestInspectNode struct {
	Body struct {
		Addr string `json:"addr"`
	}
}

func NewRequestInspectNode(ctx *gin.Context) (req *RequestInspectNode, err error) {
	req = &RequestInspectNode{}
	if err = ctx.ShouldBindJSON(&req.Body); err != nil {
		return nil, fmt.Errorf("failed to bind json: %w", err)
	}

	if _, err := types.NodeAddressFromBech32(req.Body.Addr); err != nil {
		return nil, fmt.Errorf("failed to parse addr: %w", err)
	}

	return req, nil
}
