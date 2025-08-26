// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package tg

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	tb "gopkg.in/telebot.v3"
)

type TGHook struct {
	Config TGConfig
	bot    *tb.Bot
}

type TGConfig struct {
	Level  logrus.Level
	API    string
	Token  string
	ChatID int64
}

func NewTGHook(cfg TGConfig) (*TGHook, error) {
	tbcfg := tb.Settings{
		Token:  cfg.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}
	if len(cfg.API) > 0 && strings.HasPrefix(cfg.API, "http") {
		tbcfg.URL = cfg.API
	}
	bot, err := tb.NewBot(tbcfg)
	if err != nil {
		return nil, err
	}
	return &TGHook{
		Config: cfg,
		bot:    bot,
	}, nil
}

func (hook *TGHook) Fire(entry *logrus.Entry) error {
	var notifyErr string
	if err, ok := entry.Data["error"].(error); ok {
		notifyErr = err.Error()
	} else {
		notifyErr = entry.Message
	}
	message := fmt.Sprintf("%s: %s", strings.ToUpper(entry.Level.String()), notifyErr)
	return hook.sendMessage(message, hook.Config.ChatID)
}

func (hook *TGHook) Levels() []logrus.Level {
	return logrus.AllLevels[:hook.Config.Level+1]
}

func (hook *TGHook) sendMessage(msg string, to int64) error {
	opt := &tb.SendOptions{}
	_, err := hook.bot.Send(tb.ChatID(to), msg, opt)
	return err
}
