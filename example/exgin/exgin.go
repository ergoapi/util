package main

import "github.com/ergoapi/util/exgin"

func main() {
	g := exgin.Init(&exgin.Config{
		Debug: true,
	})
	g.Use(exgin.ExHackHeader())
	g.Run()
}
