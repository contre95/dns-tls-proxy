package helpers

import (
	"errors"
	"log"
	"tls-dns-proxy/pkg/domain/proxy"

	"golang.org/x/net/dns/dnsmessage"
)

type msgParser struct {
}

func NewMsgParser() proxy.MsgParser {
	return &msgParser{}
}

func (mp *msgParser) PackMsg(dnsm *dnsmessage.Message, msgFormat string) (proxy.Msg, error) {
	m, err := dnsm.Pack()
	if err != nil {
		return nil, err
	}
	if msgFormat == "tcp" {
		var tcpBytes []byte
		tcpBytes = make([]byte, 2)
		tcpBytes[0] = 0
		tcpBytes[1] = byte(len(m))
		m = append(tcpBytes, m...)
		return m, nil
	} else if msgFormat == "udp" {
		return m, nil
	} else {
		return nil, errors.New("Invalid msg format")
	}
}

func (mp *msgParser) ParseUDPMsg(m proxy.Msg) (*dnsmessage.Message, error) {
	var dnsm dnsmessage.Message
	err := dnsm.Unpack(m[:])
	if err != nil {
		log.Printf("Unable to parse UDP Message: %v \n", err)
		return nil, errors.New("Unable to unpack request, invalid message.")
	}

	return &dnsm, nil
}

func (mp *msgParser) ParseTCPMsg(m proxy.Msg) (*dnsmessage.Message, error) {
	var dnsm dnsmessage.Message
	err := dnsm.Unpack(m[2:])
	if err != nil {
		log.Printf("Unable to parse TCP Message: %v \n", err)
		return nil, errors.New("Unable to unpack request, invalid message.")
	}
	return &dnsm, nil
}
