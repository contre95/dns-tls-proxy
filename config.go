package main

import (
	"os"
	"strconv"
)

const TCP_PORT = "PROXY_CONFIG_TCP_PORT"
const UDP_PORT = "PROXY_CONFIG_UDP_PORT"
const PROXY_METHOD_PORT = "PROXY_CONFIG_METHOD"
const TCP_MAX_CONN_POOL = "PROXY_CONFIG_TCP_MAX_CONN_POOL"
const CACHE_TTL = "PROXY_CONFIG_CACHE_TTL"
const RESOLVER_READ_TO = "PROXY_RESOLVER_READ_TO"

type config struct {
	TCP_PORT          int
	UDP_PORT          int
	TCP_DIRECT        bool
	TCP_MAX_CONN_POOL int
	CACHE_TLL         int
	RESOLVER_READ_TO  uint
}

func GetConfig() config {

	cacheTTL, err := strconv.Atoi(os.Getenv(CACHE_TTL))
	if err != nil {
		panic("Config Err: Could not parse Cache TTL")
	}
	tcpPort, err := strconv.Atoi(os.Getenv(TCP_PORT))
	if err != nil {
		panic("Config Err: Could not parse TCP port")
	}
	udpPort, err := strconv.Atoi(os.Getenv(UDP_PORT))
	if err != nil {
		panic("Config Err: Could not parse UDP port")
	}
	maxConnPool, err := strconv.Atoi(os.Getenv(TCP_MAX_CONN_POOL))
	if err != nil {
		panic("Config Err: Could not parse TCP_MAX_CONN_POOL.")
	}
	resolverReadTO, err := strconv.Atoi(os.Getenv(RESOLVER_READ_TO))
	if err != nil {
		panic("Config Err: Could not parse RESOLVER_READ_TO.")
	}
	tcpProxyMehod := os.Getenv(PROXY_METHOD_PORT)
	directTCPPRoxy := false
	if tcpProxyMehod == "direct" {
		directTCPPRoxy = true
	}

	return config{
		TCP_PORT:          tcpPort,
		UDP_PORT:          udpPort,
		TCP_DIRECT:        directTCPPRoxy,
		TCP_MAX_CONN_POOL: maxConnPool,
		CACHE_TLL:         cacheTTL,
		RESOLVER_READ_TO:  uint(resolverReadTO),
	}
}
