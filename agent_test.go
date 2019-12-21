package agent

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// ALL - extend | actor/#/
var key = []string{
	"ua_o66yvMhUS_7JchL_64WMXRwWHbm9C",
	"6wXo66yvMhUS_7JchL_64WMXRwWHbm9C",
	"uxzo66yvMhUS_7JchL_64WMXRwWHbm9C",
}

const (
	server = "127.0.0.1"
	port   = 8084
	link   = "actor"
)

type testSend struct {
	DL  bool   `json:"dl"`
	URL string `json:"url"`
}

func TestActions(t *testing.T) {

	l := logrus.New()

	em1 := UseEmitter(server, port, key[0], link, l, false)
	em2 := UseEmitter(server, port, key[1], link, l, false)
	em3 := UseEmitter(server, port, key[2], link, l, true)

	// Local

	local, err := New(em1)
	assert.NoError(t, err)

	assert.NoError(t, local.Start())

	local.Do("ping", sayThanks)

	local.Do("dl", dl)

	// Local 2

	local2, err := New(em2)
	assert.NoError(t, err)

	assert.NoError(t, local2.Start())

	local2.Do("ping", sayThanks)

	local2.Do("dl", dl)

	// Remote

	remote, err := New(em3)
	assert.NoError(t, err)

	assert.NoError(t, remote.Start())

	remote.Do("msg", parseMessage)

	remote.SendRaw("ping", "remote ping ball")

	sendDT := testSend{true, "https://google.com"}

	sendDL, err := PrepSender("remote", "dl", sendDT)

	remote.Send(sendDL)

	local.SendRaw("msg", "local1 msg red")

	time.Sleep(1 * time.Second)

}

func parseMessage(_ Client, content string) {
	fmt.Printf("remote recieved: %s\n", content)
}

func sayThanks(fr Client, content string) {
	fr.SendRaw("msg", "local2 msg thanks")
}

func dl(fr Client, content string) {

	var resp testSend
	err := json.Unmarshal([]byte(content), &resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("local received %+v\n", resp)

}
