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
