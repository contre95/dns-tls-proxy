package proxy

import (
	"errors"
	"fmt"
	"io"
	"log"
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
}

func NewDNSProxy(r Resolver, mp MsgParser, c Cache) Service {
	return &service{r, mp, c}
}

func (s *service) Direct(conn *net.Conn) error {
	resolverConn, err := s.resolver.GetTLSConnection()
	if err != nil {
		log.Println("Could not get TLS Resolver connection")
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
		dnsm, um, parseErr = s.mparser.ParseUPDMsg(um)
	} else if msgFormat == "tcp" {
		dnsm, parseErr = s.mparser.ParseTCPMsg(um)
	} else {
		return nil, errors.New(fmt.Sprintf("Invalid msg format: %s \n", msgFormat))
	}
	if parseErr != nil {
		log.Printf("Error parsing UnsolvedMsg: %v \n", parseErr)
		return nil, parseErr
	}

	// Check if the response is cached
	cm, cacheErr := s.cache.Get(dnsm)
	if cacheErr != nil {
		log.Printf("Cache error: %v", cacheErr)
	}

	// If cache could resolve the query, then try with the resolver
	if cm == nil {
		sm, resolutionErr := s.resolver.Solve(um)
		if resolutionErr != nil {
			fmt.Printf("Resolution Error: %v \n", resolutionErr)
			return nil, resolutionErr
		}
		cacheErr := s.cache.Store(dnsm, sm)
		if cacheErr != nil {
			log.Printf("Error: Could not cache message: %v.", cacheErr)
		}
		return sm, nil
	}
	return cm, nil
}