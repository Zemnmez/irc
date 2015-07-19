package irc

import (
	"bufio"
	"bytes"
	"encoding"
	"errors"
	"github.com/Zemnmez/irc/internal/util"
	"io"
	"strings"
)

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = &Message{}

type MessageDecoder struct {
	r *bufio.Reader
}

var crlfTokenizer = util.NewMultiToken("\r\n")

func (d MessageDecoder) Decode() (m Message, err error) {
	rns, err := crlfTokenizer.ScanToken(d.r)
	if err != nil {
		return
	}

	if err = m.UnmarshalText([]byte(string(rns))); err != nil {
		return
	}

	return
}

func NewMessageDecoder(r io.Reader) (m MessageDecoder) {
	m.r = bufio.NewReader(r)

	return
}

//Message represents a single IRC message,
//such as this:
//
//	:prefix command parameter1 param2 param3
//
type Message struct {
	Prefix  string
	Command string
	Params  []string //up to 15
}

var ErrMsgTooShort = errors.New("irc: message too short")
var ErrMsgTooLong = errors.New("irc: message longer than 510 characters")
var ErrParamsTooLong = errors.New("irc: more than 15 parameters")

func (m Message) Encode(w io.Writer) (err error) {
	var b []byte
	if b, err = m.MarshalText(); err != nil {
		return
	}

	if _, err = w.Write(b); err != nil {
		return
	}

	if _, err = io.WriteString(w, "\r\n"); err != nil {
		return
	}

	return
}

//func (m *Message) Decode(r io.Reader) (err error) {
//
//}

//UnmarshalText parses the IRC message into the Message
//struct.
func (m *Message) UnmarshalText(b []byte) (err error) {
	m.Params = make([]string, 0, 15)

	s := string(b)
	if len(s) > 510 {
		err = ErrMsgTooLong
		return
	}

	if len(b) < 1 {
		err = ErrMsgTooShort
		return
	}

	// prefix
	if b[0] == ':' {
		i := bytes.IndexByte(b, ' ')
		b, m.Prefix = b[i+1:], string(b[1:i])
	}

	//suffix
	var suffix *string
	if i := bytes.IndexByte(b, ':'); i != -1 {
		var sp string
		b, sp = b[:i-1], string(b[i+1:])
		suffix = &sp
	}

	//command
	if i := bytes.IndexByte(b, ' '); i != -1 {
		b, m.Command = b[i+1:], string(b[:i])
	}

	var c int
	for {
		i := bytes.IndexByte(b, ' ')

		//no more spaces
		if i == -1 {
			m.Params = append(m.Params, string(b))
			c++
			break
		}

		b, m.Params = b[i+1:], append(m.Params, string(b[:i]))
		c++

	}

	if c > 15 {
		err = ErrParamsTooLong
		return
	}

	if suffix != nil {
		m.Params = append(m.Params, *suffix)
	}

	if m.Command == "" {
		err = ErrMsgTooShort
		return
	}

	return
}

func (m Message) MarshalText() (b []byte, err error) {

	if m.Prefix != "" {
		b = []byte(
			strings.Join(
				append(
					[]string{
						":" + m.Prefix,
						m.Command,
					}, m.Params...,
				),
				" ",
			),
		)
		return
	}

	b = []byte(
		strings.Join(
			append(
				[]string{
					m.Command,
				}, m.Params...,
			),
			" ",
		),
	)

	return
}
