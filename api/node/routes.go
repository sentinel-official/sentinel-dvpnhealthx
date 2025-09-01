package node

import (
	"github.com/gin-gonic/gin"

	"github.com/sentinel-official/sentinel-dvpnhealthx/core"
)

func RegisterRoutes(r gin.IRouter, c *core.Context) {
	r.GET("/nodes/:addr", handlerGetNode(c))
	r.GET("/nodes", handlerGetNodes(c))
	r.POST("/nodes", handlerInspectNode(c))
}
