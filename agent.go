package agent

// Perform is the function that occurs.
type Perform = func(from Client, data string)

// Client has the connection that receives instruction.
type Client interface {
	Start() error
	Do(trigger string, action Perform)
	SendRaw(trigger string, data string)
	Send(Sender)
}

// Config gives options.
type Config interface {
	Prepare() (Client, error)
}

// Logger allows for information to be given.
type Logger interface {
	Debugf(string, ...interface{})
	Errorf(string, ...interface{})
}

// Sender is used to send data to a remote agent.
type Sender interface {
	Source() string
	Trigger() string
	Data() string
}

// New returns an actor.
func New(c Config) (Client, error) {
	return c.Prepare()
}
