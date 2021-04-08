package slack

import (
	"encoding/json"
	"errors"
	"github.com/imroc/req"
)

const (
	ApiURL = "https://slack.com/api/chat.postMessage"
)

// Options allows full configuration of the message sent to the Pushover API
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

func (c *client) Send(message string) error {
	if c.opt.Token == "" {
		return errors.New("missing token")
	}
	if c.opt.Channel == "" {
		return errors.New("missing user")
	}
	if message == "" {
		return errors.New("missing message")
	}
	c.opt.Text = message
	inrec, _ := json.Marshal(c.opt)
	params := &req.Param{}
	json.Unmarshal(inrec, params)
	resp, err := req.Post(ApiURL, *params)
	if err != nil {
		return nil
	}
	r := &Resp{}
	err = resp.ToJSON(r)
	if err != nil {
		return nil
	}
	if !r.Ok {
		return errors.New(r.Error)
	}
	return nil
}
