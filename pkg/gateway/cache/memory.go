package cache

import (
	"errors"
	"fmt"
	"time"
	"tls-dns-proxy/pkg/domain/proxy"

	"golang.org/x/net/dns/dnsmessage"
)

func main() {
	fmt.Println("vim-go")
}

type memCache struct {
	ttl time.Time
}

func NewMemoryCache(ttl time.Time) proxy.Cache {
	return &memCache{ttl}
}

func (mc *memCache) Get(dnsm *dnsmessage.Message) (proxy.SolvedMsg, error) {
	return nil, errors.New("Not implemented")
}

func (mc *memCache) Store(dnsm *dnsmessage.Message, sm proxy.SolvedMsg) error {
	return errors.New("Not implemented")
}
