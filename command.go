package irc

import "regexp"

type Command string

// replyCommandRegexp is a regular expression that can be used to check if a
// given command is a (numeric) IRC message reply code.
var replyCommandRegexp = regexp.MustCompile("\\d{3}")

const (
	JoinCommand Command = "JOIN"
	NickCommand Command = "NICK"
	OperCommand Command = "OPER"
	PassCommand Command = "PASS"
	PingCommand Command = "PING"
	PongCommand Command = "PONG"
	UserCommand Command = "USER"
	QuitCommand Command = "QUIT"
)

// Numerics in the range from 001 to 099 are used for client-server
// connections only and should never travel between servers.  Replies
// generated in the response to commands are found in the range from 200
// to 399.
//
// The server sends Replies 001 to 004 to a user upon
// successful registration.
const (

	// "Welcome to the Internet Relay Network
	// <nick>!<user>@<host>"
	WelcomeReply Command = "001"

	// "Your host is <servername>, running version <ver>"
	YourHostReply Command = "002"

	// "This server was created <date>"
	CreatedReply Command = "003"

	// "<servername> <version> <available user modes>
	MyInfoReply Command = "004"

	// Sent by the server to a user to suggest an alternative
	// server.  This is often used when the connection is
	// refused because the server is already full.
	//
	// "Try server <server name>, port <port number>"
	BounceReply Command = "005"

	// Reply format used by USERHOST to list replies to
	// the query list.  The reply string is composed as
	// follows:
	//
	// reply = nickname [ "*" ] "=" ( "+" / "-" ) hostname
	//
	// The '*' indicates whether the client has registered
	// as an Operator.  The '-' or '+' characters represent
	// whether the client has set an AWAY message or not
	// respectively.
	//
	// ":*1<reply> *( " " <reply> )"
	UserHostReply Command = "302"
)

// String returns a string-representation of the command.
// Mainly used to satisfy the "runtime.stringer" interface and to allow for "%s"-format strings.
func (c Command) String() string {
	return string(c)
}

// IsNumericReply checks if the command is a (numeric, 3-digit) reply-code.
func (c Command) IsNumericReply() bool {
	return replyCommandRegexp.MatchString(c.String())
}
