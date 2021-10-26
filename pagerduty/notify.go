package pagerduty

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	ApiURL = "https://events.pagerduty.com/v2/enqueue"
)

type Options struct {
	Token string `json:"token"`
	Text  string `json:"text"`
}

type pagerduty struct {
	Payload     payload `json:"payload"`
	RoutingKey  string  `json:"routing_key"`
	EventAction string  `json:"event_action"`
}

type payload struct {
	Summary  string `json:"summary"`
	Source   string `json:"source"`
	Severity string `json:"severity"`
}

type pagerdutyRes struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	DedupKey string `json:"dedup_key"`
}

type client struct {
	opt Options
}

type Resp struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

func New(opt Options) *client {
	return &client{opt: opt}
}

func (c *client) Send(message string) error {
	err := c.check(message)
	if err != nil {
		return err
	}

	pdOpt := &pagerduty{
		Payload: payload{
			Summary:  message,
			Source:   "monitoringtool:cloudvendor:central-region-dc-01:852559987:cluster/api-stats-prod-003",
			Severity: "info",
		},
		RoutingKey:  c.opt.Token,
		EventAction: "trigger",
	}

	inrec, _ := json.Marshal(pdOpt)
	resp, err := http.Post(ApiURL, "application/json", bytes.NewBuffer(inrec))
	if err != nil {
		return fmt.Errorf("pagerduty error: %s", err)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)
	res := &pagerdutyRes{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return fmt.Errorf("pagerduty server error: %s", string(body))
	}

	if res.Status != "success" {
		return fmt.Errorf("send notify failed: %s", string(body))
	}
	return nil
}

func (c *client) check(msg string) error {
	if c.opt.Token == "" {
		return errors.New("missing auth token")
	}

	if msg == "" {
		return errors.New("missing message")
	}
	return nil
}
