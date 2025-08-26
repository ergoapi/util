// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

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
