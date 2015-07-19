package irc

import (
	"fmt"
	"reflect"
	"testing"
)

func TestUnmarshalTooLong(t *testing.T) {
	const s = "command " +
		// this is 520 letter "a"s
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	var m Message

	if err := m.UnmarshalText([]byte(s)); err != ErrMsgTooLong {

		t.Fatalf("We didn't complain that "+
			"a message that was too long "+
			"was too long.%s", err)
	}
}

func TestUnmarshalTooShort(t *testing.T) {
	const s = ""

	var m Message

	if err := m.UnmarshalText([]byte(s)); err != ErrMsgTooShort {

		fmt.Println("32: ", s, m, err)

		t.Fatalf("We didn't complain that "+
			"a message that was too short "+
			"was too short. %s %+v", err, m)
	}
}

func TestUnmarshalParamsTooLong(t *testing.T) {
	const s = "command 1 2 3 4 5 6 7 8 9 10 11 " +
		"12 13 14 15 16" // one more parameter

	var m Message

	if err := m.UnmarshalText([]byte(s)); err != ErrParamsTooLong {

		t.Fatalf("We didn't complain that "+
			"a message that had too many params "+
			"had too many. %+q %+v", err, m)
	}
}

func TestLongMessage(t *testing.T) {
	const s = ":prefix command 1 2 3 4 5 6 7 8 9 10 11 " +
		"12 13 14 :hi this is a really long message lol" // one more parameter

	var m Message

	if err := m.UnmarshalText([]byte(s)); err != nil {

		t.Fatalf("%s %+v", err, m)
	}
}

func TestMarshalUnmarshalMessage(t *testing.T) {
	m := Message{
		Prefix:  "prefix",
		Command: "command",
		Params:  []string{"a", "b", "c"},
	}

	var b []byte
	var err error
	if b, err = m.MarshalText(); err != nil {
		panic(err)
	}

	t.Logf("Message string: %s", b)

	var m2 Message
	if err = m2.UnmarshalText(b); err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(m, m2) {
		t.Fatalf(
			"%+v != %+v",
			m,
			m2,
		)
	}
}
