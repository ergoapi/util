package main

import "github.com/ergoapi/util/exkube"

func main() {
	kubeClient, err := exkube.New(&exkube.ClientConfig{})
	if err != nil {
		panic(err)
	}
	v, err := kubeClient.Discovery().ServerVersion()
	if err != nil {
		panic(err)
	}
	println(v.String())
}
