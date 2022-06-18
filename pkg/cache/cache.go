package cache

import (
	"fmt"
)

type Cache interface {
	Set(key string, value string) error
	Get(key string) (string, error)
}

type memoryCache struct {
	underlyingData map[string]string
}

func (p *memoryCache) Set(key string, value string) error {
	p.underlyingData[key] = value
	return nil
}

func (p *memoryCache) Get(key string) (string, error) {
	value, ok := p.underlyingData[key]
	if !ok {
		return "", fmt.Errorf("Key %s not found", key)
	}

	return value, nil
}

func NewMemoryCache() Cache {
	pCache := &memoryCache{}
	pCache.underlyingData = make(map[string]string)

	return pCache
}
