package irc

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
)

// ConnectionState is a bit mask that determines the current connection state.
type ConnectionState int8

func (cs ConnectionState) String() string {
	var str string
	switch cs {
	case ConnectionStateClosed:
		str = "CLOSED"
	case ConnectionStateClosing:
		str = "CLOSING"
	case ConnectionStateOpen:
		str = "OPEN"
	case ConnectionStateReady:
		str = "READY"
	default:
		str = fmt.Sprintf("UNKNOWN (%d)", cs)
	}
	return str
}

// ConnectionAlreadyEstablished is an error that is being used to indicate that
// the connection to a server has already been
var ConnectionAlreadyEstablished = fmt.Errorf("connection has already been established")

// connectionMsgBufSize defines the buffer size to be used for incoming and outgoing
// messages that are send and received by a connection.
const connectionMsgBufSize = 64

const (
	// ConnectionStateClosed indicates, that the connection has been closed.
	// Pseudo state where all the bits have been flipped of.
	ConnectionStateClosed ConnectionState = 0x0
	// ConnectionStateClosing indicates, that the connection is about to be closed.
	ConnectionStateClosing ConnectionState = 0x1
	// ConnectionStateOpen indicates, that the TCP connection has been established
	// and messages can be send to the server and received from there.
	ConnectionStateOpen ConnectionState = 0x2
	// ConnectionStateReady indicates, that a connection can fully be used.
	// Usually, this is the case, after the NICK, USER and possible the PASS commands have
	// been send and successfully acknowledged by the server.
	ConnectionStateReady ConnectionState = 0x4
)

// ClientConnection encapsulates details about a client<>server connection.
type ClientConnection interface {
	State() <-chan ConnectionState
	Hostname() string
	Port() int
	Capabilities() []Capability
	HasCapability(Capability) bool
	In() <-chan Message
	Out() chan<- Message
	Err() <-chan error
	// Open establishes a connection to the configured IRC server.
	// If the connection cannot be established, "ConnectionAlreadyEstablished" will be returned
	// as error. If a "real" error happens while trying to establish the connection, this error
	// will be returned instead.
	Open() (err error)
	Wait()
	io.Closer
}

type clientConnection struct {
	state        chan ConnectionState
	hostname     string
	port         int
	capabilities map[Capability]bool
	tcpConn      net.Conn // Underlying TCP connection.
	in           chan Message
	out          chan Message
	err          chan error
	wg           sync.WaitGroup
}

type ConnectionHandler func(
	conn *ClientConnection,
	state <-chan ConnectionState,
	in <-chan *Message,
	out chan<- *Message,
	err <-chan error,
)

// NewClientConnection prepares a new connection that can be used to connect to the
// given IRC server.
func NewClientConnection(hostname string, port int) ClientConnection {
	return &clientConnection{
		state:        make(chan ConnectionState, 4),
		hostname:     hostname,
		capabilities: make(map[Capability]bool),
		port:         port,
		in:           make(chan Message, connectionMsgBufSize), // from server
		out:          make(chan Message, connectionMsgBufSize), // to server
		err:          make(chan error),                         // message-related errors
		wg:           sync.WaitGroup{},
	}
}

func (conn *clientConnection) Open() (err error) {
	if conn.tcpConn != nil {
		err = ConnectionAlreadyEstablished
		return
	}

	addr := fmt.Sprintf("%s:%d", conn.hostname, conn.port)
	if conn.tcpConn, err = net.Dial("tcp", addr); err != nil {
		err = fmt.Errorf("connection to IRC server %s failed: %v", addr, err)
		conn.err <- err
	} else {
		conn.wg.Add(1)

		// INPUT and CONNECTION CHECKS
		go func() {
			defer conn.wg.Done()
			defer func() { conn.state <- ConnectionStateClosed }()
			conn.state <- ConnectionStateOpen
			reader := bufio.NewReader(conn.tcpConn)
			scanner := bufio.NewScanner(reader)
			for conn.tcpConn != nil && scanner.Scan() {
				str := string(scanner.Text())
				var msg Message
				if msg, err = NewMessageFromString(str); err != nil {
					conn.err <- err
					continue
				}

				// Some messages will be used to trigger certain actions at this point.
				// They will still be relayed to the "in" channel, though.
				switch msg.Command() {
				case WelcomeReply:
					// Once we receive RPL_WELCOME, we can rest assured that the connection to
					// the server has been established successfully and the USER and NICK commands
					// have been acknowledged.
					conn.state <- ConnectionStateReady
				case PingCommand:
					// PING messages will be handled directly at this point, thus a PONG reply is
					// going to be send immediately.
					conn.out <- NewPongMessage(EmptyPrefix, msg.Parameters()[0])
				default:
					break
				}

				conn.in <- msg
			}
		}()

		// OUTPUT
		go func() {
			for conn.tcpConn != nil {
				select {
				case msg := <-conn.out:
					conn.send(msg)
				case state := <-conn.state:
					if state != ConnectionStateReady {
						break
					}
				}
			}
		}()

	}

	return
}

func (conn *clientConnection) Hostname() string {
	return conn.hostname
}

func (conn *clientConnection) Port() int {
	return conn.port
}

func (conn *clientConnection) State() <-chan ConnectionState {
	return conn.state
}

func (conn *clientConnection) In() <-chan Message {
	return conn.in
}

func (conn *clientConnection) Out() chan<- Message {
	return conn.out
}

func (conn *clientConnection) Err() <-chan error {
	return conn.err
}

// Wait will block (thus: wait) until the connection has been closed down.
func (conn *clientConnection) Wait() {
	conn.wg.Wait()
}

// send transmits an IRC message
// There is usually no need to call this method directly, because the provided
// channel, which can be accessed by the "Out()" method offers a much better
// way to dispatch messages.
func (conn *clientConnection) send(msg Message) (err error) {
	str := fmt.Sprintf("%s\r\n", msg.String())
	if _, err = fmt.Fprint(conn.tcpConn, str); err != nil {
		err = fmt.Errorf("could not send message: %v", err)
	}
	return
}

// Close closes a client connection.
// If the connection has already been closed, no attempt to
// close the connection (again) will be made.
// If the underlying TCP socket cannot be closed, the error
// that has been returned from attempting to close the socket
// will be returned and should be handled by the caller.
func (conn *clientConnection) Close() (err error) {
	conn.state <- ConnectionStateClosing
	if conn.tcpConn != nil {
		err = conn.tcpConn.Close()
		conn.state <- ConnectionStateClosed
	}
	conn.tcpConn = nil
	return
}

// Capabilities lists all the capabilities that the
// server reports as being supported.
func (conn *clientConnection) Capabilities() []Capability {
	caps := make([]Capability, len(conn.capabilities))
	i := 0
	act := 0
	for k := range conn.capabilities {
		if conn.capabilities[k] {
			caps[act] = k
			act++
		}
		i++
	}
	return caps[:act]
}

// HasCapability checks if the server that the client is connected to supports a given capability.
func (conn *clientConnection) HasCapability(capability Capability) bool {
	v, e := conn.capabilities[capability]
	return e && v
}
