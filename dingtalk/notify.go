package dingtalk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/imroc/req"
	"net/url"
	"time"
)

type Options struct {
	WebhookUrl string `json:"webhook_url"`
	Secret     string `json:"secret"`
}

type client struct {
	opt Options
}

func New(opt Options) *client {
	return &client{opt: opt}
}

type Resp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func (c *client) Send(message string) error {
	if "" == c.opt.WebhookUrl {
		return errors.New("missing webhook url")
	}

	if "" == c.opt.Secret {
		return errors.New("missing secret")
	}

	if "" == message {
		return errors.New("missing message")
	}

	sign, timestamp := c.getSign()

	header := req.Header{
		"Content-Type": "application/json",
	}

	dingUrl := fmt.Sprintf("%s&timestamp=%d&sign=%s", c.opt.WebhookUrl, timestamp, sign)
	json := fmt.Sprintf("{\"msgtype\": \"text\",\"text\": {\"content\":\"%s\"}}", message)

	resp, _ := req.Post(dingUrl, json, header)
	r := &Resp{}
	err := resp.ToJSON(&r)
	if err != nil {
		return err
	}
	if r.Errcode != 0 {
		return errors.New(r.Errmsg)
	}
	return nil
}

func (c *client) getSign() (string, int64) {
	timestamp := time.Now().UnixMilli()
	secret := c.opt.Secret

	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(stringToSign))
	signData := hash.Sum(nil)
	sign := url.QueryEscape(base64.StdEncoding.EncodeToString(signData))
	return sign, timestamp
}
