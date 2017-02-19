package irc

import (
	"fmt"
	"testing"
)

var validUrlStrs = [...]string{
	"irc://irc.example.com",
	"irc://irc.example.com:6667",
	"irc://irc.example.com/#channel",
	"irc://127.0.0.1/",
	"irc://127.0.0.1:6667",
	"irc://127.0.0.1/#channel",
	"ircs://irc.example.com:6697/#channel",
}

var invalidUrlStrs = [...]string{
	"http://www.example.com/",
}

func TestNewUrl(t *testing.T) {
	for _, str := range validUrlStrs {
		url, err := NewURL(str)
		if err != nil || !url.IsValid() {
			t.Errorf("IRC URL should be valid but couldn't be parsed: %s", url)
		}
	}
	for _, str := range invalidUrlStrs {
		url, err := NewURL(str)
		if err == nil || url.IsValid() {
			t.Errorf("string should be invalid IRC URL, but could somehow be parsed: %s", str)
		}

	}
}

func TestURL_String(t *testing.T) {

	str := "irc://irc.example.com:6667/"
	var url fmt.Stringer
	url, _ = NewURL(str)
	if str != url.String() {
		t.Error("Implementation of fmt.Stringer interface produces unexpected output")
	}

	for _, str := range validUrlStrs {
		url, _ := NewURL(str)
		if str != url.String() && str != url.String() {
			t.Error("invocation of ToString yields unexpected result")
		}
	}
}
