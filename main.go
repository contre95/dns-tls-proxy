package main

import (
	"time"
	"tls-dns-proxy/pkg/domain/proxy"
	"tls-dns-proxy/pkg/gateway/cache"
	"tls-dns-proxy/pkg/gateway/logger"
	"tls-dns-proxy/pkg/gateway/resolver"
	"tls-dns-proxy/pkg/helpers"
	"tls-dns-proxy/pkg/presenter/socket"
)

func main() {
	config := GetConfig()
	stdCacheLogger := logger.NewSTDLogger("CACHE")
	cache := cache.NewMemoryCache(time.Duration(config.CACHE_TLL)*time.Second, stdCacheLogger)
	go cache.AutoPurge()
	resolver := resolver.NewCloudFlareResolver("1.1.1.1", 853, config.RESOLVER_READ_TO)
	parser := helpers.NewMsgParser()
	stdProxyLogger := logger.NewSTDLogger("PROXY")
	proxy := proxy.NewDNSProxy(resolver, parser, cache, stdProxyLogger)

	go socket.StarUDPtServer(proxy, config.UDP_PORT, "0.0.0.0")
	socket.StartTCPServer(proxy, config.TCP_PORT, "0.0.0.0", config.TCP_DIRECT, config.TCP_MAX_CONN_POOL)
}
