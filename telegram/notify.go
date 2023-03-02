package telegram

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ApiURL = "https://api.telegram.org/bot"
)

// https://core.telegram.org/bots/api#sendmessage
// https://github.com/go-telegram-bot-api/telegram-bot-api

// Options allows full configuration of the message sent to the Pushover API
type Options struct {
	Token    string `json:"token"`
	Channel  int64  `json:"channel"`
	ChatName string `json:"chat_name"`
	// User may be either a user key or a group key.
}

type client struct {
	opt Options
	bot *tgbotapi.BotAPI
}

func New(opt Options) *client {
	api, err := tgbotapi.NewBotAPI(opt.Token)
	if err != nil {
		return nil
	}

	return &client{opt: opt, bot: api}
}

type Resp struct {
	Status int      `json:"status"`
	Errors []string `json:"errors"`
}

func (c *client) Send(message string) error {
	if c.opt.Token == "" {
		return errors.New("missing token")
	}

	if message == "" {
		return errors.New("missing message")
	}
	var msg tgbotapi.MessageConfig
	if c.opt.Channel != 0 {
		msg = tgbotapi.NewMessage(c.opt.Channel, message)
	} else {
		msg = tgbotapi.NewMessageToChannel(c.opt.ChatName, message)
	}

	sendRes, err := c.bot.Send(msg)
	if err != nil {
		return err
	}

	fmt.Println("sendRes:", sendRes)
	return nil
}
