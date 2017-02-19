package irc

import (
	"fmt"
	"regexp"
	"strconv"
)

// Default port to be used for IRC server connections
const DefaultServerPort int = 6667

// Default port to be used for IRC server connections when using TLS.
// See RfC-7194 for further details.
const DefaultServerPortTls int = 6697

// urlRegex can be used to match valid IRC URLs.
var urlRegexp = regexp.MustCompile(`(?P<Protocol>ircs?)(?:://)(?P<Host>[a-z0-9\\.\\-]*)(?::(?P<Port>\d+))?(?:/(?P<Channel>[#&+!][a-zA-Z0-9#&+!]{1,50}))?`)

type URL interface {
	String() string
	IsValid() bool
	Protocol() string
	Hostname() string
	Port() int
}

type url struct {
	valid    bool
	str      string
	port     int
	hostname string
	protocol string
}

// NewURL creates a new IRC URL.
// The given string str is parsed and the normalized portions of the URL
// are stored to allow for fast access. If parsing of the URL string fails,
// the returned error err will be non-nil.
func NewURL(str string) (u URL, err error) {
	uStruct := &url{
		str:   str,
		valid: urlRegexp.MatchString(str),
	}
	u = uStruct
	if !u.IsValid() {
		err = fmt.Errorf("invalid URL: %s", u)
		return
	}

	match := urlRegexp.FindStringSubmatch(u.String())
	portSet := false
	for i, name := range urlRegexp.SubexpNames() {
		value := match[i]
		switch name {
		case "Protocol":
			uStruct.protocol = value
		case "Host":
			uStruct.hostname = value
		case "Port":
			uStruct.port, _ = strconv.Atoi(value)
			portSet = true
		}

	}

	// If the port has not (yet) been set, we'll use the default port for the given protocol.
	if !portSet {
		switch uStruct.protocol {
		case "irc":
			uStruct.port = DefaultServerPort
		case "ircs":
			uStruct.port = DefaultServerPortTls
		default:
			// We should NOT be able to reach this point, because the regular
			// expression doesn't allow it. But... hey... can we be really sure? ;-)
			panic(fmt.Errorf("unknown protocol %s", uStruct.protocol))
		}
	}

	return
}

func (u *url) Hostname() string {
	return u.hostname
}

func (u *url) Protocol() string {
	return u.protocol
}

func (u *url) Port() int {
	return u.port
}

// String returns a string representation of the IRC URL.
func (u *url) String() string {
	return u.str
}

// IsValid validates the correctness of a given irc://-URL.
func (u *url) IsValid() bool {
	return u.valid
}
