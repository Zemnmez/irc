//Package util provides some helpful functions.
package util

import (
	"bufio"
)

//MultiToken represents a run of runes that can be used to tokenize a
//stream of runes.
type MultiToken []rune

//ScanToken returns the run of runes terminated by m, not including the
//terminator.
func (m MultiToken) ScanToken(rd *bufio.Reader) (rns []rune, err error) {
	var pos int // position in m

	for r := rune(0); ; {
		r, _, err = rd.ReadRune()

		rns = append(rns, r)
		switch {
		case err != nil:
			return
		case r == m[pos]:
			pos++
			if pos == len(m) {
				rns = rns[:len(rns)-len(m)]
				return
			}
		default:
			pos = 0
		}

	}
}

func NewMultiToken(s string) MultiToken {
	return MultiToken([]rune(s))
}
