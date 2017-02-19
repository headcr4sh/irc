package irc

import (
	"testing"
)

func TestCommand_String(t *testing.T) {
	var testData = []struct {
		vrn string
		cmd Command
		str string
	}{
		{"JoinCommand", JoinCommand, "JOIN"},
		{"NickCommand", NickCommand, "NICK"},
		{"YourHostReply", YourHostReply, "002"},
	}
	for _, td := range testData {
		if str := td.cmd.String(); str != td.str {
			t.Errorf("%s.String() -> %s, expected: %s", td.vrn, str, td.str)
		}
	}
}

func TestCommand_IsNumericReply(t *testing.T) {
	var testData = []struct {
		cmd Command
		isr bool
	}{
		{JoinCommand, false},
		{OperCommand, false},
		{BounceReply, true},
		{WelcomeReply, true},
		{YourHostReply, true},
	}
	for _, td := range testData {
		if isr := td.cmd.IsNumericReply(); isr != td.isr {
			t.Errorf("%s.IsNumericReply() -> %b, expected: %b", td.cmd, isr, td.isr)
		}
	}
}
