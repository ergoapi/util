package main

import (
	"github.com/ergoapi/util/exhttp"
	"github.com/ergoapi/util/version"
	_ "github.com/ergoapi/util/version/prometheus"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	exhttp.MetricsHandler(g)
	g.GET("/", func(c *gin.Context) {
		c.JSON(200, version.Get())
	})
	g.Run()
}
