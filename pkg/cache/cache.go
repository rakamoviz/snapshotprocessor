package cache

import (
	"context"
	"fmt"
)

type Cache interface {
	Set(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (string, error)
}

type memoryCache struct {
	underlyingData map[string]string
}

func (p *memoryCache) Set(ctx context.Context, key string, value string) error {
	p.underlyingData[key] = value
	return nil
}

func (p *memoryCache) Get(ctx context.Context, key string) (string, error) {
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
