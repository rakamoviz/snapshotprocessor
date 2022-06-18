package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryCache(t *testing.T) {
	key, value := "clave", "valor"

	pCache := NewMemoryCache()
	pCache.Set(key, value)

	expectedValue := value
	actualValue, err := pCache.Get(key)

	assert.Nil(t, err, "The Get from cache must succeed")
	assert.Equal(t, expectedValue, actualValue, "The value from Get not the same as the value in Set")
}
