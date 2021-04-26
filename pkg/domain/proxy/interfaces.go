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
	Get(dnsm *dnsmessage.Message) (SolvedMsg, error)
	Store(dnsm *dnsmessage.Message, sm SolvedMsg) error
}
type MsgParser interface {
	ParseUPDMsg(m Msg) (*dnsmessage.Message, UnsolvedMsg, error)
	ParseTCPMsg(m Msg) (*dnsmessage.Message, error)
}
