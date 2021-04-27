package proxy

import (
	"crypto/tls"

	"golang.org/x/net/dns/dnsmessage"
)

type Msg []byte
type UnsolvedMsg = Msg
type SolvedMsg = Msg

type Resolver interface {
	Solve(um UnsolvedMsg) (SolvedMsg, error)
	GetTLSConnection() (*tls.Conn, error)
}
type Cache interface {
	Get(dnsm dnsmessage.Message) (*dnsmessage.Message, error)
	Store(dnsm dnsmessage.Message) error
	AutoPurge()
}
type MsgParser interface {
	ParseUDPMsg(m Msg) (*dnsmessage.Message, error)
	ParseTCPMsg(m Msg) (*dnsmessage.Message, error)
	PackMsg(dnsm *dnsmessage.Message, msgFormat string) (Msg, error)
}
