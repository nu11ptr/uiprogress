package strutil

import (
	"testing"
	"time"
)

func TestLen(t *testing.T) {
	tests := []struct {
		name, input string
		output      int
	}{
		{"Empty", "", len("")},
		{"Plain", "abc abc", len("abc abc")},
		{"FalsePositive", "abc\x1babc", len("abc abc")},
		{"SetAndResetColor", string([]rune{
			esc, lBracket, '3', '6', 'm', 'a', 'b', 'c', ' ', 'a', 'b', 'c', esc, lBracket, '0', 'm',
		}), len("abc abc")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := Len(test.input)
			if l != test.output {
				t.Error("want", test.output, "got", l)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name, input, output string
		l                   int
		trunc               bool
	}{
		{"TruncNotNeeded", "abc", "abc", 3, false},
		{"Basic", "abcabc", "abc", 3, true},
		{"HasEscapes", "abc\x1b[36mabc", "abc\x1b[36mab", 5, true},
		{"FalsePositive", "abc\x1babc", "abc\x1ba", 5, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s2, trunc := Truncate(test.input, test.l)
			if s2 != test.output {
				t.Error("want", test.output, "got", s2)
			}
			if trunc != test.trunc {
				t.Error("want", test.trunc, "got", trunc)
			}
		})
	}
}

func TestResize(t *testing.T) {
	s := "foo"
	got := Resize(s, 5)
	if len(got) != 5 {
		t.Fatal("want", 5, "got", len(got))
	}
	s = "foobar"
	got = Resize(s, 5)

	if got != "fo..." {
		t.Fatal("want", "fo...", "got", got)
	}
}

func TestPadRight(t *testing.T) {
	got := PadRight("foo", 5, '-')
	if got != "foo--" {
		t.Fatal("want", "foo--", "got", got)
	}
}

func TestPadLeft(t *testing.T) {
	got := PadLeft("foo", 5, '-')
	if got != "--foo" {
		t.Fatal("want", "--foo", "got", got)
	}
}

func TestPrettyTime(t *testing.T) {
	d, _ := time.ParseDuration("")
	got := PrettyTime(d)
	if got != "---" {
		t.Fatal("want", "---", "got", got)
	}
}
