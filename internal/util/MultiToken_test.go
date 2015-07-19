package util

import (
	"bufio"
	"fmt"
	"strings"
)

func ExampleMultiToken_ScanToken() {
	const input = "hi\r\nhow\r\nare\r\nyou"
	var tok = MultiToken([]rune{'\r', '\n'})

	r := bufio.NewReader(strings.NewReader(input))

	for {
		rs, err := tok.ScanToken(r)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Got input: %s\n", string(rs))
	}

	// Output:
	// Got input: hi
	// Got input: how
	// Got input: are
	// EOF
}
