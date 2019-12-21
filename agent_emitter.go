package agent

import (
	"fmt"
	"path"
	"strings"

	emitter "github.com/emitter-io/go/v2"
)

// UseEmitter returns emitter config struct.
func UseEmitter(domain string, port int, key, link string, l Logger, debug bool) Config {
	return &useEmitter{
		log:    l,
		server: domain,
		port:   port,
		key:    key,
		link:   path.Join(link),
	}
}

type useEmitter struct {
	log    Logger
	server string
	port   int
	key    string
	link   string
}

func (u *useEmitter) Prepare() (Client, error) {

	s := fmt.Sprintf("tcp://%s:%v", u.server, u.port)

	client, err := emitter.Connect(s, func(c *emitter.Client, msg emitter.Message) {
		u.log.Errorf("unknown message: %s %s", msg.Topic(), string(msg.Payload()))
	}, emitter.WithAutoReconnect(true))

	if err != nil {
		return nil, err
	}

	client.OnError(func(c *emitter.Client, err emitter.Error) {

		if err.Error() == "the security key provided is not authorized to perform this operation" {
			u.log.Debugf("key: %s", u.key)
		}

		u.log.Errorf("%v\n", err)
	})

	return &ea{
		client:   client,
		log:      u.log,
		key:      u.key,
		link:     u.link,
		handlers: make(map[string]Perform),
	}, nil

}

// CLIENT
type ea struct {
	client   *emitter.Client
	handlers map[string]Perform
	link     string
	log      Logger
	key      string
}

func (e *ea) Start() error {

	link, err := e.client.CreatePrivateLink(e.key, e.link+"/", "s", e.actions)

	e.log.Debugf("subscribed to %s\n", link.Channel)

	return err

}

func (e *ea) actions(_ *emitter.Client, msg emitter.Message) {

	request := strings.SplitN(string(msg.Payload()), " ", 3)
	if len(request) != 3 {
		e.log.Errorf("received invalid message: %s\n", string(msg.Payload()))
		return
	}

	from := request[0]
	trigger := request[1]
	message := request[2]

	if fn, ok := e.handlers[trigger]; ok {

		e.log.Debugf("received the following from %s on %s channel: %s", from, trigger, message)

		fr := &ea{
			client: e.client,
			key:    e.key,
			log:    e.log,
		}

		fn(fr, message)

	}
}

func (e *ea) Do(trigger string, action Perform) {

	e.log.Debugf("adding %s to handler list\n", trigger)

	e.handlers[trigger] = action
}

func (e *ea) SendRaw(topic string, data string) {

	e.log.Debugf("sending %s\n", data)

	if err := e.client.PublishWithLink("s", data); err != nil {
		e.log.Errorf("unable to publish the following data to %s\n%s", topic, data)
	}

}

func (e *ea) Send(da Sender) {

	e.log.Debugf("sending to %s from %s\n", da.Trigger(), da.Source())

	if err := e.client.PublishWithLink("s", fmt.Sprintf("%s %s %s", da.Source(), da.Trigger(), da.Data())); err != nil {
		e.log.Errorf("unable to send data to %s\n", da.Trigger())
	}

}
