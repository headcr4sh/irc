package irc

import (
	"strings"
	"testing"
)

// Tests the behavior of the validation function when facing
// correct channel names.
func TestIsValidChannelName(t *testing.T) {
	testdata := []struct {
		name  string
		valid bool
	}{
		{"#test", true},
		{"&seCrET", true},
		{"##golang", true},
		{"!haha", true},
		{"+plusChan", true},
		{"@what?", false},
		{"'yankeedoo'", false},
	}
	for _, tt := range testdata {
		if valid := isValidChannelName(tt.name); valid != tt.valid {
			t.Errorf(`IsValidChannelName("%s") => %v, expected: %v`, tt.name, valid, tt.valid)
		}
	}
}

func TestNewChannel(t *testing.T) {
	ch, err := NewChannel("#test")
	if err != nil {
		t.Fail()
	}
	if ch == nil {
		t.Fail()
	}
}

func TestChannel_String(t *testing.T) {
	const name = "#test"
	ch, _ := NewChannel(name)
	str := ch.String()
	if str == "" || !strings.Contains(str, name) {
		t.Fail()
	}
}

func TestChannel_Equal(t *testing.T) {
	testdata := []struct {
		name1 string
		name2 string
		equal bool
	}{
		{"#test", "#TEST", true},
		{"+chan", "+chan", true},
		{"#another", "+another", false},
	}
	for _, tt := range testdata {
		ch1, _ := NewChannel(tt.name1)
		ch2, _ := NewChannel(tt.name2)
		if equal := ch1.Equal(ch2); equal != tt.equal {
			t.Errorf(`Channel<"%s">.equal(Channel<"%s"> => %v, expected: %v`, ch1.Name(), ch2.Name(), equal, tt.equal)
		}
	}

}
