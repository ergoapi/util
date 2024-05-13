package main

import (
	"fmt"
	"time"

	"github.com/ergoapi/util/exgin"
	ratelimit "github.com/ergoapi/util/feat/ginmid/ratelimit"
	"github.com/gin-gonic/gin"
)

func main() {
	g := exgin.Init(&exgin.Config{
		Debug: true,
	})
	g.Use(exgin.ExHackHeader())
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Minute,
		Limit: 3,
	})
	mw := ratelimit.RateLimiter(store, &ratelimit.Options{
		Mark: "local",
	})
	g.Use(mw)
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
