package main

import (
	"bytes"
	"fmt"
	"github.com/zemnmez/irc"
	"io"
	"net"
	"os"
)

type rwc struct {
	io.ReadWriteCloser           // from the conn
	teeReader          io.Reader // everything we get from the server
	multiWriter        io.Writer // capture everything we send
}

var passBytes = []byte(pass)
var stars = []byte("***")

func (r rwc) Write(p []byte) (n int, err error) {
	if _, err = r.ReadWriteCloser.Write(p); err != nil {
		return
	}

	if _, err = os.Stdout.Write(bytes.Replace(p, passBytes, stars, -1)); err != nil {
		return
	}

	return
}
func (r rwc) Read(p []byte) (n int, err error) { return r.teeReader.Read(p) }
func newRwc(_rwc io.ReadWriteCloser) (r rwc) {
	r.ReadWriteCloser = _rwc
	r.teeReader = io.TeeReader(_rwc, os.Stdout)
	r.multiWriter = io.MultiWriter(os.Stdout, _rwc)

	return
}

func do() (err error) {
	conn, err := net.Dial("tcp", "irc.twitch.tv:6667")
	if err != nil {
		return
	}

	c := irc.NewClient(newRwc(conn))
	defer c.Close()

	finished := make(chan bool)
	//Print the messages we get
	go func(finished <-chan bool) {

		d := irc.NewMessageDecoder(c)

		var m irc.Message
		for {
			m, err = d.Decode()
			if err != nil {
				panic(err)
			}

			_ = m
		}

	}(finished)

	if err = c.Pass(pass); err != nil {
		return
	}

	if err = c.Nick("pdupcja9srq3trmtl8nah2hmx"); err != nil {
		return
	}

	if err = c.Join("#zemnmez"); err != nil {
		return
	}

	<-finished

	return
}

func main() {
	if err := do(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
