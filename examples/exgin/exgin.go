package main

import (
	"fmt"

	"github.com/ergoapi/util/exgin"
	"github.com/gin-gonic/gin"
)

func main() {
	g := exgin.Init(&exgin.Config{
		Debug: true,
	})
	g.Use(exgin.ExHackHeader())
	g.GET("/ping", func(ctx *gin.Context) {
		exgin.GinsData(ctx, "pong", nil)
	})
	g.GET("/error", func(ctx *gin.Context) {
		exgin.GinsData(ctx, "pong", fmt.Errorf("error"))
	})
	g.GET("/errormap", func(ctx *gin.Context) {
		exgin.GinsCodeData(ctx, 351, "pong", fmt.Errorf("451error"))
	})
	g.Run()
}
