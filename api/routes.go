package api

import (
	"github.com/gin-gonic/gin"

	"github.com/sentinel-official/sentinel-dvpnhealthx/api/node"
	"github.com/sentinel-official/sentinel-dvpnhealthx/core"
)

func RegisterRoutes(r gin.IRouter, c *core.Context) {
	node.RegisterRoutes(r, c)
}
