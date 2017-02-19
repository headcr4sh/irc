package irc

import (
	"reflect"
	"testing"
)

var exampleServerReplies = map[string]*message{
	":irc.example.net 001 testuser :Welcome testuser!~john.doe@172.17.0.1": {
		prefix:     &prefix{ptype: PrefixHostname, hostname: "irc.example.net"},
		command:    WelcomeReply,
		parameters: []string{"testuser", "Welcome testuser!~john.doe@172.17.0.1"},
	},
	":irc.some.host.org 002 testuser :Your host is irc.some.host.org, running version ngircd-23 (x86_64/alpine/linux-musl)": {
		prefix:     &prefix{ptype: PrefixHostname, hostname: "irc.some.host.org"},
		command:    YourHostReply,
		parameters: []string{"testuser", "Your host is irc.some.host.org, running version ngircd-23 (x86_64/alpine/linux-musl)"},
	},
	":irc.my-special-host.co.uk 003 testuser :This server has been started Mon Mar 06 2017 at 22:01:14 (UTC)": {
		prefix:     &prefix{ptype: PrefixHostname, hostname: "irc.my-special-host.co.uk"},
		command:    CreatedReply,
		parameters: []string{"testuser", "This server has been started Mon Mar 06 2017 at 22:01:14 (UTC)"},
	},
	":localhost 004 johndoe localhost ngircd-23 abBcCFiIoqrRswx abehiIklmMnoOPqQrRstvVz": {
		prefix:     &prefix{ptype: PrefixHostname, hostname: "localhost"},
		command:    MyInfoReply,
		parameters: []string{"johndoe", "localhost", "ngircd-23", "abBcCFiIoqrRswx", "abehiIklmMnoOPqQrRstvVz"},
	},
	":irc.example.com 005 testuser CHANNELLEN=50 NICKLEN=9 TOPICLEN=490 AWAYLEN=127 KICKLEN=400 MODES=5 MAXLIST=beI:50 EXCEPTS=e INVEX=I PENALTY :are supported on this server": {
		prefix:     &prefix{ptype: PrefixHostname, hostname: "irc.example.com"},
		command:    BounceReply,
		parameters: []string{"testuser", "CHANNELLEN=50", "NICKLEN=9", "TOPICLEN=490", "AWAYLEN=127", "KICKLEN=400", "MODES=5", "MAXLIST=beI:50", "EXCEPTS=e", "INVEX=I", "PENALTY", "are supported on this server"},
	},
}

func TestNewPrefixFromString(t *testing.T) {
	var testdata = []struct {
		str string
		pfx *prefix
	}{
		{
			"irc.example.com",
			&prefix{
				ptype:    PrefixHostname,
				hostname: "irc.example.com",
			},
		},
		{
			"johndoe!jdoe@client.example.com",
			&prefix{
				ptype:    PrefixNicknameUserHost,
				nickname: "johndoe",
				user:     "jdoe",
				host:     "client.example.com",
			},
		},
		{
			"janedoe@anotherclient.example.com",
			&prefix{
				ptype:    PrefixNicknameHost,
				nickname: "janedoe",
				host:     "anotherclient.example.com",
			},
		},
	}
	for _, tt := range testdata {
		pfx := NewPrefixFromString(tt.str)
		if pfx.Type() != tt.pfx.Type() || pfx.Hostname() != tt.pfx.Hostname() || pfx.Nickname() != tt.pfx.Nickname() || pfx.User() != tt.pfx.User() || pfx.Host() != tt.pfx.Host() {
			t.Errorf("NewPrefixFromString(%s) -> %v (type: %v) , expected: %v (type: %v)", tt.str, pfx, pfx.Type(), tt.pfx, tt.pfx.Type())
		}
		if str := pfx.String(); tt.str != str {
			t.Errorf("<prefix>.String() -> %s, expected: %s", str, tt.str)
		}
	}
}

func TestNewMessage(t *testing.T) {
	prefix := EmptyPrefix
	command := NickCommand
	parameters := []string{"john_doe"}
	msg := NewMessage(prefix, command, parameters...)
	if msg.Prefix() != prefix {
		t.Errorf("unexpected Prefix, expected %s, got %s", prefix, msg.Prefix())
	}
	if msg.Command() != command {
		t.Errorf("unexpected Command, expected %s, got %s", command, msg.Command())
	}
	if msg.Parameters()[0] != parameters[0] {
		t.Errorf("unexpected Parameters, expected %s, got %s", parameters[0], msg.Parameters()[0])
	}
}

func TestMessage_IsValid(t *testing.T) {
	m := NewMessageWithoutPrefix(UserCommand, "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16")
	valid, errs := m.IsValid()
	if valid {
		t.Error("message should not be valid")
	} else if errs == nil || len(errs) == 0 {
		t.Error("errors should be reported when validating invalid message")
	}
}

func TestMessage_String(t *testing.T) {
	for k, v := range exampleServerReplies {
		msg := NewMessage(v.prefix, v.command, v.parameters...)
		str := msg.String()
		if str != k {
			t.Errorf("expected '%s', got '%s'", k, str)
		}
	}
}

func TestParseMessage(t *testing.T) {
	for k, v := range exampleServerReplies {
		msg, err := NewMessageFromString(k)
		if err != nil {
			t.Errorf("could not parse valid server reply: \"%s\"", k)
		} else if !reflect.DeepEqual(v, msg) {
			t.Errorf("unexpected result while parsing message: \"%s\"", k)
		}
	}
}

func TestQuitMessage_Reason(t *testing.T) {
	r1 := "Bye, folks!"
	m := NewQuitMessage(EmptyPrefix, r1)
	if r2 := m.Reason(); r2 != r1 {
		t.Errorf(`QuitMessage.Reaon() -> "%s", expected: "%s"`, r2, r1)
	}
}
