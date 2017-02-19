package irc

type Capability string

const (

	// The account-notify spec defines a way for clients to be notified when other clients login to accounts.
	// This spec defines the ACCOUNT message to enable this, use of the a WHOX token, as well as outlining the
	// general restriction of account names not being * (as this is used to indicate logging out of accounts).
	AccountNotify Capability = "account-notify"

	// The account-tag spec defines a way for clients to receive a message tag on messages specifying the current
	// account that other client is logged into (or that they aren’t logged into one at all). This is especially
	// useful for letting bots make use of the network’s authentication and account mechanisms.
	AccountTag Capability = "account-tag"

	// The away-notify extension provides a way for clients to instantly know when other clients go away or come
	// back. This improves responsiveness and the display of channels for IRC clients that display this information.
	AwayNotify Capability = "away-notify"

	// The batch extension provides a way for servers to mark certain messages as related. This can simplify the
	// display of this information in clients as well as allow better post-processing on them.
	Batch Capability = "batch"

	// The cap-notify spec allows clients to be sent notifications when caps are added to or removed from the
	// server. This is useful in cases like SASL when the authentication layer disconnects (and thus, SASL
	// authentication is no longer possible). This extension is automatically enabled if clients request v3.2
	// capability negotiation.
	CapNotify Capability = "cap-notify"

	Chghost Capability = "chghost"

	EchoMessage Capability = "echo-message"

	// The extended-join spec defines a way to request that extra client information (including that client’s
	// account) is sent when clients join a given channel. This allows better tracking of accounts, particularly
	// when used with account-notify.
	ExtendedJoin Capability = "extended-join"

	InivteNotify Capability = "invite-notify"

	Metadata Capability = "metadata"

	Monitor Capability = "monitor"

	MultiPrefix Capability = "multi-prefix"

	SASL Capability = "sasl"

	// The server-time extension allows clients to see the exact time that messages were sent and received.
	// This allows bouncers to replay information with more accurate time tracking.
	ServerTime Capability = "server-time"

	TLS Capability = "tls"

	// The userhost-in-names extension allows clients to more easily see the user/hostnames of other clients when
	// joining channels. This allows clients to better track info and automate client features more easily.
	UserhostInNames Capability = "userhost-in-names"
)

func (c Capability) String() string {
	return string(c)
}
