package loader

/*
This is a not type safe data loaders for go, inspired by https://github.com/facebook/dataloader.

The intended use is in graphql servers, to reduce the number of queries being sent to e.g. a database.

*/

import (
	"context"
	"time"

	cache "github.com/OrlovEvgeny/go-mcache"
)

const (
	// DefaultTTL is the default TTL used if nothing else is specified
	DefaultTTL = time.Minute * 10
)

type (
	// FetchFunc abstracts the process of loading a resource
	FetchFunc func(context.Context, string) (interface{}, error)

	// Loader holds cached resources. The cache is a simple in-memory cache with TTL.
	Loader struct {
		fetch        FetchFunc
		c            *cache.CacheDriver
		expiresAfter time.Duration
	}
)

// New initializes the loader
func New(f FetchFunc, ttl time.Duration) *Loader {
	return &Loader{
		fetch:        f,
		c:            cache.New(),
		expiresAfter: ttl,
	}
}

// Load returns either a cached instance or calls the fetch function to retriece the requested instance
func (l *Loader) Load(ctx context.Context, key string) (interface{}, error) {

	if data, ok := l.c.Get(key); ok {
		return data, nil
	}

	data, err := l.fetch(ctx, key)

	if err != nil {
		return nil, err
	}
	if data != nil {
		if err := l.c.Set(key, data, l.expiresAfter); err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, nil
}
