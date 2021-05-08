package cache

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
	"tls-dns-proxy/pkg/domain/proxy"

	"golang.org/x/net/dns/dnsmessage"
)

type cacheValue struct {
	msg        *dnsmessage.Message
	expiration time.Time
}

type memCache struct {
	ttl     time.Duration
	mx      sync.RWMutex
	entries map[string]cacheValue
	logger  proxy.Logger
}

func (mc *memCache) AutoPurge() {
	for now := range time.Tick(time.Second) {
		for key, cValue := range mc.entries {
			if cValue.expiration.Before(now) {
				mc.mx.Lock()
				mc.logger.Info("Clearing entry: %v \n", key)
				delete(mc.entries, key)
				mc.mx.Unlock()
			}
		}
	}
}
func NewMemoryCache(ttl time.Duration, logger proxy.Logger) proxy.Cache {
	c := memCache{
		ttl:     ttl,
		entries: map[string]cacheValue{},
		logger:  logger,
	}
	return &c
}

func (mc *memCache) hashKey(dnsm dnsmessage.Message) string {

	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", dnsm.Questions)))
	return fmt.Sprintf("%x", h.Sum(nil))

}

func (mc *memCache) Get(dnsm dnsmessage.Message) (*dnsmessage.Message, error) {
	mc.mx.RLock()
	defer mc.mx.RUnlock()
	if cValue, ok := mc.entries[mc.hashKey(dnsm)]; ok {
		return cValue.msg, nil
	}
	return nil, nil
}

func (mc *memCache) Store(dnsm dnsmessage.Message) error {
	mc.mx.Lock()
	mc.entries[mc.hashKey(dnsm)] = cacheValue{&dnsm, time.Now().Add(mc.ttl)}
	mc.mx.Unlock()
	return nil
}
