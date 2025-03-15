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
		exgin.SucessResponse(ctx, "pong")
	})
	g.GET("/error", func(ctx *gin.Context) {
		exgin.ErrorResponse(ctx, 400, fmt.Errorf("error"))
	})
	g.GET("/errormap", func(ctx *gin.Context) {
		exgin.ErrorResponse(ctx, 351, fmt.Errorf("451error"))
	})
	g.Run()
}
