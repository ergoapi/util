// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package main

import (
	"github.com/ergoapi/util/log/hooks/tg"

	"github.com/sirupsen/logrus"
)

func main() {
	tghook, err := tg.NewTGHook(tg.TGConfig{
		// API:    "",
		Token:  "000000000:XXXXXXXX",
		ChatID: 0,
		Level:  logrus.WarnLevel,
	})
	if err != nil {
		panic(err)
	}
	logrus.AddHook(tghook)
	logrus.Error("error log!")
	logrus.Warn("warn log.")
}
