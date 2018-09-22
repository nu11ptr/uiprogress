// Package strutil provides various utilities for manipulating strings
package strutil

import (
	"bytes"
	"time"
)

type strState int

const (
	notInEscape strState = iota
	inEscape
	// startingEscape means we have seen the ESC char but are not yet in a full escape sequence
	startingEscape

	esc      = '\x1b'
	lBracket = '['
	execEsc  = 'm'
)

// Len computes the length of a string, but unlike the builtin len, it ignores ANSI escape codes
func Len(s string) (count int) {
	state := notInEscape

	for _, c := range s {
		switch state {
		case notInEscape:
			if c == esc {
				state = startingEscape
			} else {
				count++
			}
		case inEscape:
			if c == execEsc {
				state = notInEscape
			}
		case startingEscape:
			if c == lBracket {
				state = inEscape
			} else {
				state = notInEscape
				// We increment count because this escape was a false positive and wasn't counted earlier
				// Additionally, this 2nd char (that wasn't lBracket) was also not counted and should be
				count += 2
			}
		}
	}
	return
}

// PadRight returns a new string of a specified length in which the end of the current string is padded with spaces or with a specified Unicode character.
func PadRight(str string, length int, pad byte) string {
	if Len(str) >= length {
		return str
	}
	buf := bytes.NewBufferString(str)
	for i := 0; i < length-Len(str); i++ {
		buf.WriteByte(pad)
	}
	return buf.String()
}

// PadLeft returns a new string of a specified length in which the beginning of the current string is padded with spaces or with a specified Unicode character.
func PadLeft(str string, length int, pad byte) string {
	if Len(str) >= length {
		return str
	}
	var buf bytes.Buffer
	for i := 0; i < length-Len(str); i++ {
		buf.WriteByte(pad)
	}
	buf.WriteString(str)
	return buf.String()
}

// Resize resizes the string with the given length. It ellipses with '...' when the string's length exceeds
// the desired length or pads spaces to the right of the string when length is smaller than desired
func Resize(s string, length uint) string {
	n := int(length)
	if Len(s) == n {
		return s
	}
	// Pads only when length of the string smaller than len needed
	s = PadRight(s, n, ' ')
	// FIXME: This won't work if there is an escape sequence at the end of the string
	if Len(s) > n {
		b := []byte(s)
		var buf bytes.Buffer
		for i := 0; i < n-3; i++ {
			buf.WriteByte(b[i])
		}
		buf.WriteString("...")
		s = buf.String()
	}
	return s
}

// PrettyTime returns the string representation of the duration. It rounds the time duration to a second and returns a "---" when duration is 0
func PrettyTime(t time.Duration) string {
	if t == 0 {
		return "---"
	}
	return (t - (t % time.Second)).String()
}
