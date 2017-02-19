package irc

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Maximum message length for any IRC message according to RfC-2812
const MsgMaxLen = 512

// Regular expression used to validate nicknames.
var nickNameRegexp = regexp.MustCompile("\\A[a-z_\\-\\[\\]\\\\^{}|`][a-z0-9_\\-\\[\\]\\\\^{}|`]*\\z")

// MessageDelimiter is the message delimiter that is sent after each
// message (Carriage-return + line-feed)
const messageDelimiter string = "\r\n"

// MessagePartSeparator is used in raw message to separate the Prefix,
// the Command and the Parameters from each other. (Space, thus ASCII: 0x20)
const messagePartSeparator string = " "

// messagePrefixPresenceIndicator is used in the beginning of raw messages
// to indicate the presence of a message Prefix. (Colon, thus ASCII: 0x3b)
var messagePrefixPresenceIndicator = ":"

// TrailingMessageParameterPresenceIndicator is used in message to indicate
// that the last message parameter may contain whitespace characters.
const trailingMessagePartPresenceIndicator = " :"

// MaxMessageParameterCount defines the maximum amount of Parameters that
// can be contained within a message.
const MaxMessageParameterCount = 15

// EmptyPrefix is used to mark messages that do not contain a prefix.
var EmptyPrefix Prefix = &prefix{
	ptype: PrefixEmpty,
}

type PrefixType int

const (
	PrefixEmpty            PrefixType = -1
	PrefixHostname         PrefixType = 0
	PrefixNicknameUserHost            = 1
	PrefixNicknameHost     PrefixType = 2
)

func (p PrefixType) String() string {
	switch p {
	case PrefixEmpty:
		return "PrefixEmpty"
	case PrefixHostname:
		return "PrefixHostname"
	case PrefixNicknameUserHost:
		return "PrefixNicknameUserHost"
	case PrefixNicknameHost:
		return "PrefixNicknameHost"
	default:
		panic(fmt.Errorf("unknown prefix type: %d", p))
	}
}

// Prefix encapsulates the contents of an IRC message prefix.
// According to the specification, a prefix can either contain the hostname of the server
// from which a message originated or the nickname and host (and optionally the user) of
// a message's origin. Before the contents of a prefix can be analyzed, it is therefore
// mandatory to check it's type, because that's how the contents can be determined.
type Prefix interface {
	Type() PrefixType
	Hostname() string
	Nickname() string
	User() string
	Host() string
	fmt.Stringer
}

type prefix struct {
	ptype    PrefixType
	hostname string
	nickname string
	user     string
	host     string
}

func (p *prefix) Type() PrefixType {
	return p.ptype
}

func (p *prefix) Hostname() string {
	return p.hostname
}

func (p *prefix) Nickname() string {
	return p.nickname
}

func (p *prefix) User() string {
	return p.user
}

func (p *prefix) Host() string {
	return p.host
}

func (p *prefix) String() string {
	if p == nil {
		return ""
	}
	switch p.ptype {
	case PrefixEmpty:
		return ""
	case PrefixHostname:
		return fmt.Sprintf("%s", p.hostname)
	case PrefixNicknameUserHost:
		return fmt.Sprintf("%s!%s@%s", p.nickname, p.user, p.host)
	case PrefixNicknameHost:
		return fmt.Sprintf("%s@%s", p.nickname, p.host)
	default:
		panic(fmt.Errorf("unexpected PrefixType: %d", p.ptype))
	}
}

func NewPrefixFromString(str string) Prefix {
	pfx := &prefix{}
	if strings.IndexRune(str, '@') != -1 {
		restAndHost := strings.Split(str, "@")
		pfx.host = restAndHost[1]
		if strings.IndexRune(restAndHost[0], '!') != -1 {
			pfx.ptype = PrefixNicknameUserHost
			nicknameAndUser := strings.Split(restAndHost[0], "!")
			pfx.nickname = nicknameAndUser[0]
			pfx.user = nicknameAndUser[1]

		} else {
			pfx.ptype = PrefixNicknameHost
			pfx.nickname = restAndHost[0]
		}

	} else {
		pfx.ptype = PrefixHostname
		pfx.hostname = str
	}
	return pfx
}

type Message interface {
	Prefix() Prefix
	Command() Command
	Parameters() []string
	IsValid() (valid bool, errs []error)
	fmt.Stringer
}

/**
 * Message is a structure that contains the parts of a raw string message after
 * they parsing of the message has been performed.
 */
type message struct {
	prefix     Prefix
	command    Command
	parameters []string
}

func (msg *message) Prefix() Prefix {
	return msg.prefix
}

func (msg *message) Command() Command {
	return msg.command
}

func (msg *message) Parameters() []string {
	return msg.parameters
}

func NewMessage(prefix Prefix, command Command, parameters ...string) Message {
	return &message{
		prefix:     prefix,
		command:    command,
		parameters: parameters,
	}
}

func NewMessageWithoutPrefix(command Command, parameters ...string) Message {
	return NewMessage(EmptyPrefix, command, parameters...)
}

// NewMessageFromString create a new message by parsing a raw CR-LF-terminated
// raw string as received from a connection.
func NewMessageFromString(rawStr string) (msg Message, err error) {

	var prefix Prefix
	var command Command
	var parameters []string

	// Let's first cut of the CR-LF message separator
	rawStr = strings.TrimRight(rawStr, messageDelimiter)

	// Checks, if the message contains a Prefix and processes it if it is present.
	if strings.HasPrefix(rawStr, messagePrefixPresenceIndicator) {
		rawStr = strings.TrimLeft(rawStr, messagePrefixPresenceIndicator)
		pfxAndRest := strings.SplitN(rawStr, messagePartSeparator, 2)
		prefix = NewPrefixFromString(pfxAndRest[0])
		rawStr = pfxAndRest[1]
	}

	// Extracts the Command / reply code from the message and processes it.
	cmdAndParams := strings.SplitN(rawStr, messagePartSeparator, 2)
	command = Command(cmdAndParams[0])
	rawStr = cmdAndParams[1]

	// Now let's check the Parameters
	var trailingParam *string = nil
	if strings.Contains(rawStr, trailingMessagePartPresenceIndicator) {
		paramsMiddleAndTrailing := strings.SplitN(rawStr, trailingMessagePartPresenceIndicator, 2)
		trailingParam = &paramsMiddleAndTrailing[1]
		rawStr = paramsMiddleAndTrailing[0]
	}

	// And now all the non-trailing (middle) Parameters...
	middleParams := strings.Split(rawStr, messagePartSeparator)
	pCount := len(middleParams)
	if trailingParam != nil {
		pCount++
	}
	parameters = make([]string, pCount)
	for i := 0; i < len(middleParams); i++ {
		parameters[i] = middleParams[i]
	}
	if trailingParam != nil {
		parameters[pCount-1] = *trailingParam
	}

	msg = NewMessage(prefix, command, parameters...)
	return
}

// IsValid checks a message for validity. Several checks will be performed to
// check whether the contents of the message complies with the various RfCs and
// specifications. The returned value valid will be set to false, if any check
// detects issues with the message in which case the returned slice of errors err
// will contain at least one issue with the message.
func (msg *message) IsValid() (valid bool, errs []error) {
	pc := len(msg.parameters)
	msgLen := len(msg.String()) + len(messageDelimiter)
	if msgLen > MsgMaxLen {
		errs = append(errs, fmt.Errorf("message length is %d bytes, which exceeds the allowed maximum of %d bytes", msgLen, MsgMaxLen))
	}
	if pc > MaxMessageParameterCount {
		errs = append(errs, fmt.Errorf("message has %d parameters, which exceeds the allowed maximum of %d parameters", pc, MaxMessageParameterCount))
	}
	valid = len(errs) == 0
	return
}

// String converts a message to a string representation.
// This operation achieves the exact opposite of the NewMessageFromString
// function.
func (msg *message) String() (str string) {
	str = fmt.Sprintf(":%v %s", msg.prefix, msg.command)
	if msg.parameters != nil {
		for _, p := range msg.parameters {
			if strings.ContainsRune(p, ' ') {
				str += fmt.Sprintf(" :%s", p)
			} else {
				str += fmt.Sprintf(" %s", p)
			}
		}
	}
	return
}

type PassMessage interface {
	Message
	Password() string
}

type passMessage struct {
	message
}

func (msg *passMessage) Password() string {
	return msg.Parameters()[0]
}

func NewPassMessage(prefix Prefix, password string) (msg PassMessage) {
	return &passMessage{
		message{
			prefix,
			PassCommand,
			[]string{password},
		},
	}
}

type PongMessage interface {
	Message
	Server1() string
}

type pongMessage struct {
	message
}

func (msg *pongMessage) Server1() string {
	return msg.parameters[0]
}

func NewPongMessage(prefix Prefix, server1 string) PongMessage {
	return &pongMessage{
		message{
			prefix,
			PongCommand,
			[]string{server1},
		},
	}
}

func NickMessage(nickname string) (msg Message) {
	msg = NewMessageWithoutPrefix(NickCommand, nickname)
	if !isValidNickname(nickname) {
		panic(fmt.Errorf("invalid nickname: %s", nickname))
	}
	return
}

type UserMessage interface {
	Message
	Username() string
	Realname() string
	Mode() UserModes
}

type userMessage struct {
	message
}

func (msg *userMessage) Username() string {
	return msg.parameters[0]
}

func (msg *userMessage) Mode() UserModes {
	if mode, err := strconv.Atoi(msg.parameters[1]); err != nil {
		panic(err)
	} else {
		return UserModesFromBitmask(mode)
	}
}

func (msg *userMessage) Realname() string {
	return msg.parameters[3]
}

func NewUserMessage(prefix Prefix, username string, realname string, mode ...UserMode) (msg Message) {
	// TODO(headcr4sh): Validate username
	return &userMessage{
		message{
			prefix,
			UserCommand,
			[]string{username, strconv.Itoa(UserModes(mode).Bitmask()), "*", realname},
		},
	}
}

func NewJoinMessage(channel string) (msg Message) {
	// TODO(headcr4sh): Validate channel name
	msg = NewMessageWithoutPrefix(JoinCommand, channel)
	return
}

type QuitMessage interface {
	Message
	Reason() string
}

type quitMessage struct {
	message
}

func (msg *quitMessage) Reason() string {
	return msg.Parameters()[0]
}

func NewQuitMessage(prefix Prefix, reason string) (msg QuitMessage) {
	return &quitMessage{
		message{
			prefix,
			QuitCommand,
			[]string{reason},
		},
	}
}
