package irc

import (
	"os"
	"strconv"
	"testing"
)

func TestConnection_Capabilities(t *testing.T) {
	conn := &clientConnection{}
	conn.hostname = "example.com"
	conn.port = DefaultServerPort
	conn.capabilities = make(map[Capability]bool)
	conn.capabilities[AccountNotify] = true
	conn.capabilities[CapNotify] = false
	caps := conn.Capabilities()
	if len(caps) != 1 {
		t.Errorf("expected a list that contains exactly 1 capability, but found: %d", len(caps))
	}
	if caps[0] != AccountNotify {
		t.Errorf(`expected a list that contains the capability "%s", but found "%s" instead`, AccountNotify, caps[0])
	}
}

func TestConnection__HasCapability(t *testing.T) {
	conn := &clientConnection{}
	conn.hostname = "example.com"
	conn.port = DefaultServerPort
	conn.capabilities = make(map[Capability]bool)
	conn.capabilities[AccountNotify] = true
	conn.capabilities[CapNotify] = false
	if !conn.HasCapability(AccountNotify) {
		t.Errorf(`expected server to support capability "%s"`, AccountNotify)
	}
	if conn.HasCapability(CapNotify) {
		t.Errorf(`didn't expect server to support capability "%s"`, CapNotify)
	}
	if conn.HasCapability(AccountTag) {
		t.Errorf(`didn't expect server to support capability "%s"`, AccountTag)
	}
}

func TestConnection__IRCD_Integration(t *testing.T) {
	server, serverSet := os.LookupEnv("IRC_IT_SERVER_HOSTNAME")
	portStr, portSet := os.LookupEnv("IRC_IT_SERVER_PORT")
	var port int
	if !serverSet {
		server = "localhost"
	}
	if portSet {
		port, _ = strconv.Atoi(portStr)
	} else {
		port = DefaultServerPort
	}
	if !testing.Short() {
		t.Logf("running integration tests against %s:%d", server, port)
		conn := NewClientConnection(server, DefaultServerPort)
		go func(t *testing.T) {
			for {
				select {
				case err := <-conn.Err():
					t.Logf("ERROR: %v\n", err)
					t.Error(err)
				case msg := <-conn.In():
					t.Logf("Message received: %v\n", msg)
					if msg.Command() == PingCommand {
						conn.Out() <- NewQuitMessage(EmptyPrefix, "Bye...")
					}
				case s := <-conn.State():
					t.Logf("ClientConnection state changed: %v\n", s)
					if s == ConnectionStateOpen {
						conn.Out() <- NickMessage("johndoe")
						conn.Out() <- NewUserMessage(EmptyPrefix, "j.doe", "John Doe")
					}
				}
			}
		}(t)
		if err := conn.Open(); err != nil {
			t.Error(err)
		}
		conn.Wait()
	}
}
