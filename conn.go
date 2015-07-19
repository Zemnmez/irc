package irc

import (
	"io"
)

type Client struct {
	io.ReadWriteCloser
}

func (c Client) SendMessage(m Message) (err error) { return m.Encode(c) }
func (c Client) Command(command string, params ...string) error {
	return c.SendMessage(Message{Command: command, Params: params})
}

func (c Client) Nick(name string) error        { return c.Command("NICK", name) }
func (c Client) Pass(password string) error    { return c.Command("PASS", password) }
func (c Client) Join(channels ...string) error { return c.Command("JOIN", channels...) }
func (c Client) JoinLockedChannel(channel, password string) error {
	return c.Command("JOIN", channel, password)
}

func NewClient(r io.ReadWriteCloser) (c Client) {
	c.ReadWriteCloser = r

	return
}
