package main

import (
	"github.com/ergoapi/util/exhttp"
	"github.com/ergoapi/util/version"
	_ "github.com/ergoapi/util/version/prometheus"
	"github.com/ergoapi/util/zos"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	exhttp.MetricsHandler(g)
	g.GET("/", func(c *gin.Context) {
		c.JSON(200, version.Get())
	})
	g.GET("/uname", func(c *gin.Context) {
		c.String(200, zos.UNAME())
	})
	g.Run()
}
