package irc

import (
	"fmt"
)

// ServerNameMaxLen specifies the maximum possible length of any given server name
// as specified by RfC-2812
const ServerNameMaxLen = 63

type ServerInfo struct {
	hostname     string
	port         int
	capabilities map[Capability]bool
}

func NewServerInfo(hostname string, port int) *ServerInfo {
	return &ServerInfo{
		hostname: hostname,
		port:     port,
	}
}

// NewServerInfoFromUrl creates a new server descriptor from the
// given URL url. If the given URL is invalid, the returned error err
// will be non-nil.
func NewServerInfoFromUrl(url URL) (srv *ServerInfo, err error) {
	if !url.IsValid() {
		err = fmt.Errorf("invalid URL: %s", url)
		return
	}
	srv = NewServerInfo(url.Hostname(), url.Port())
	return

}

func (server ServerInfo) String() string {
	return fmt.Sprintf("%s:%d", server.hostname, server.port)
}

type ServerInfoSlice []ServerInfo

func (p ServerInfoSlice) Len() int           { return len(p) }
func (p ServerInfoSlice) Less(i, j int) bool { return p[i].hostname < p[j].hostname }
func (p ServerInfoSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
