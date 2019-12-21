package agent

import (
	"encoding/json"
	"errors"
)

const (
	senderSourceErr  = "missing source id"
	senderTriggerErr = "missing trigger"
	senderDataErr    = "missing data"
)

// PrepSender is used to prepare a Sender.
func PrepSender(sender, trigger string, data interface{}) (Sender, error) {

	if sender == "" {
		return nil, errors.New(senderSourceErr)
	}

	if trigger == "" {
		return nil, errors.New(senderTriggerErr)
	}

	bd, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	d := string(bd)

	if d == "" {
		return nil, errors.New(senderDataErr)
	}

	return transfer{
		source:  sender,
		trigger: trigger,
		data:    d,
	}, nil
}

type transfer struct {
	source  string
	trigger string
	data    string
}

func (t transfer) Source() string {
	return t.source
}

func (t transfer) Trigger() string {
	return t.trigger
}

func (t transfer) Data() string {
	return t.data
}
