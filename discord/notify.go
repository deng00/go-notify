package discord

import (
	"encoding/json"
	"errors"
	"github.com/imroc/req"
)

var (
	ApiURL = "https://discord.com/api/webhooks/"
)

type Options struct {
	Token   string `json:"token"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type client struct {
	opt Options
}

func New(opt Options) *client {
	return &client{opt: opt}
}

type Resp struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type Webhook struct {
	Content string `json:"content"`
}

func (c *client) Send(message string) error {
	if "" == c.opt.Token {
		return errors.New("missing token")
	}

	if "" == c.opt.Channel {
		return errors.New("missing channel")
	}

	if "" == message {
		return errors.New("missing message")
	}
	c.opt.Text = message

	whMsg := &Webhook{
		Content: c.opt.Text,
	}

	inrec, _ := json.Marshal(whMsg)
	params := &req.Param{}
	json.Unmarshal(inrec, params)

	ApiURL = ApiURL + c.opt.Channel + "/" + c.opt.Token
	resp, err := req.Post(ApiURL, *params)
	if err != nil {
		return err
	}
	r := &Resp{}
	err = resp.ToJSON(r)

	if err != nil || err.Error() != "unexpected end of JSON input" {
		return err
	}
	return nil
}
