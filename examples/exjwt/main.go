// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/ergoapi/util/exjwt"
)

func main() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change-me-in-production"
	}
	// Create a token using env var JWT_SECRET
	_ = os.Setenv("JWT_SECRET", secret)
	token, err := exjwt.Auth("alice", "uuid-123")
	if err != nil {
		fmt.Println("auth error:", err)
		return
	}
	fmt.Println("token:", token)

	// Parse and inspect claims
	claims, err := exjwt.Parse(token)
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}
	fmt.Printf("username=%v uuid=%v sub=%v\n", claims["username"], claims["uuid"], claims["sub"])
}
