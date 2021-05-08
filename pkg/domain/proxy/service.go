package proxy

import (
	"errors"
	"fmt"
	"io"
	"net"

	"golang.org/x/net/dns/dnsmessage"
)

type Service interface {
	Solve(um UnsolvedMsg, msjFormat string) (SolvedMsg, error)
	Direct(conn *net.Conn) error
}

type service struct {
	resolver Resolver
	mparser  MsgParser
	cache    Cache
	logger   Logger
}

func NewDNSProxy(r Resolver, mp MsgParser, c Cache, l Logger) Service {
	return &service{r, mp, c, l}
}

func (s *service) Direct(conn *net.Conn) error {
	resolverConn, err := s.resolver.GetTLSConnection()
	if err != nil {
		s.logger.Err("Could not get TLS Resolver connection: %v", err)
		return err
	}
	defer resolverConn.Close()
	go io.Copy(resolverConn, *conn) // Holdea
	io.Copy(*conn, resolverConn)
	return nil
}

func (s *service) Solve(um UnsolvedMsg, msgFormat string) (SolvedMsg, error) {
	// Parse the UnsolvedMsg
	var parseErr error
	var dnsm *dnsmessage.Message
	if msgFormat == "udp" {
		dnsm, parseErr = s.mparser.ParseUDPMsg(um)
		um, parseErr = s.mparser.PackMsg(dnsm, "tcp") // Parse UDP message in order to send it to DoT Resolver (TCP)
	} else if msgFormat == "tcp" {
		dnsm, parseErr = s.mparser.ParseTCPMsg(um)
	} else {
		return nil, errors.New(fmt.Sprintf("Invalid msg format: %s \n", msgFormat))
	}
	if parseErr != nil {
		s.logger.Err("Error parsing UnsolvedMsg: %v \n", parseErr)
		return nil, parseErr
	}
	for _, q := range dnsm.Questions {
		s.logger.Info("DNS %s: %s ", msgFormat, q.Name.String())
	}

	// Check if the response is cached
	cm, cacheErr := s.cache.Get(*dnsm)
	if cacheErr != nil {
		s.logger.Err("Cache error: %v", cacheErr)
	}
	if cm != nil {
		s.logger.Debug("Message found in cache")
		cm.Header.ID = dnsm.Header.ID
		sm, parseErr := s.mparser.PackMsg(cm, msgFormat)
		if parseErr != nil {
			return nil, errors.New("Wrong value stored in cache")
		}
		return sm, nil
	}

	sm, resolutionErr := s.resolver.Solve(um)
	if resolutionErr != nil {
		s.logger.Err("Resolution Error: %v \n", resolutionErr)
		return nil, resolutionErr
	}
	var storeErr error
	var solvedDNSM *dnsmessage.Message
	solvedDNSM, storeErr = s.mparser.ParseTCPMsg(sm) // Resolver response is always TCP
	storeErr = s.cache.Store(*solvedDNSM)
	if storeErr != nil {
		s.logger.Err("Cache error: %v", cacheErr)
	}
	sm, _ = s.mparser.PackMsg(solvedDNSM, msgFormat)
	return sm, nil
}
