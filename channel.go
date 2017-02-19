package irc

import (
	"fmt"
	"regexp"
)

// channelNameRegex is a regular expression that is being used to identify valid IRC
// channel names.
var channelNameRegex = regexp.MustCompile(`[#&+!][a-zA-Z0-9#&+!]{1,49}`)

// Validates a given IRC channel name and returns either true, if
// the given name is a valid name for an IRC channel, or false
// if it not correct according to the RfC(s).
func isValidChannelName(name string) bool {
	return channelNameRegex.MatchString(name)
}

// Channel is basically a group of gathered users.
type Channel interface {
	Name() string
	Topic() string
	Equal(ch Channel) bool
	fmt.Stringer
}

type channel struct {
	name  string
	topic string
}

// NewChannel creates a new channel with the given name.
func NewChannel(name string) (ch Channel, err error) {
	if isValidChannelName(name) {
		ch = &channel{
			name: name,
		}
	} else {
		err = fmt.Errorf("invalid channel name: '%s'", name)
	}
	return
}

func (ch *channel) Name() string {
	return ch.name
}

func (ch *channel) Topic() string {
	return ch.topic
}

// Equal compares two channel definitions for equality.
// Channel names in IRC are case in-sensitive, so therefore we'll
// have to keep that in mind when comparing using string comparisons.
func (ch *channel) Equal(that Channel) bool {
	return toLowercase(ch.name) == toLowercase(that.Name())
}

// String returns a string representation of the channel, usually it's name.
func (ch *channel) String() (str string) {
	if ch == nil {
		str = "??NIL??"
	} else {
		str = ch.name
	}
	return
}
