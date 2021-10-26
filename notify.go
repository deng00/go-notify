package notify

import (
	"github.com/deng00/go-notify/pagerduty"
	"github.com/deng00/go-notify/pushover"
	"github.com/deng00/go-notify/slack"
)

type Platform string

const (
	PlatformSlack     Platform = "slack"
	PlatformPushover           = "pushover"
	PlatformDingDing           = "dingding"
	Platformpagerduty          = "pagerduty"
)

type Notify struct {
	config *Config
}

type Config struct {
	Platform Platform
	Token    string
	Channel  string
}

func NewNotify(config *Config) *Notify {
	return &Notify{
		config: config,
	}
}

func (n *Notify) Send(msg string) error {
	switch n.config.Platform {
	case PlatformPushover:
		return n.sendPushOverNotify(msg)
	case PlatformSlack:
		return n.sendSlackNotify(msg)
	case Platformpagerduty:
		return n.sendPagerdutyNotify(msg)
	default:
		panic("not supported notify platform")
	}
	return nil
}

func (n *Notify) sendPushOverNotify(msg string) error {
	app := pushover.New(pushover.Options{
		Token: n.config.Token,
		User:  n.config.Channel,
	})
	err := app.Send(msg)
	return err
}

func (n *Notify) sendSlackNotify(msg string) error {
	app := slack.New(slack.Options{
		Token:   n.config.Token,
		Channel: n.config.Channel,
	})
	err := app.Send(msg)
	return err
}

func (n *Notify) sendPagerdutyNotify(msg string) error {
	app := pagerduty.New(pagerduty.Options{
		Token: n.config.Token,
	})
	err := app.Send(msg)
	return err
}
