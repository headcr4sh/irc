package irc

import (
	"fmt"
	"sort"
	"testing"
)

func TestNewServerInfoFromUrl(t *testing.T) {
	url, err := NewURL("irc://irc.example.com:6667/#test")
	if err != nil {
		t.Fail()
	}
	srv, err := NewServerInfoFromUrl(url)
	if err != nil {
		t.Fail()
	}
	if srv.hostname != "irc.example.com" && srv.port != 6667 {
		fmt.Printf("%s:%d\n", srv.hostname, srv.port)
		t.Errorf("unexpected result while parsing URL: \"%s\"", url)
	}
}

func TestServerInfoSort(t *testing.T) {
	infos := ServerInfoSlice{
		*NewServerInfo("irc.hostc.com", 6667),
		*NewServerInfo("irc.hostb.com", 6667),
		*NewServerInfo("irc.hosta.com", 6667),
	}
	sort.Sort(infos)
	if len(infos) != 3 {
		t.Errorf("length of sorted slice must be 3, but was %d", len(infos))
	}
	if infos[0].hostname != "irc.hosta.com" || infos[1].hostname != "irc.hostb.com" || infos[2].hostname != "irc.hostc.com" {
		t.Error("sort order doesn't match expected result")
	}
}
