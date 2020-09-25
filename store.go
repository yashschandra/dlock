package dlock

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Store interface {
	Set(ctx context.Context, key, uid string, expiry time.Time) error
	Get(ctx context.Context, key string) string
	Delete(ctx context.Context, key string) error
}

const (
	keyAlreadyPresent = "key already present"
)

type value struct {
	uid    string
	expiry time.Time
}

type defaultClientImpl map[string]value

var (
	defaultClient defaultClientImpl
	mu            *sync.Mutex
)

func init() {
	mu = &sync.Mutex{}
}

func getDefaultClient() Store {
	defaultClient = make(defaultClientImpl)
	return defaultClient
}

func (c defaultClientImpl) Set(ctx context.Context, key, uid string, expiry time.Time) error {
	mu.Lock()
	defer mu.Unlock()
	v, ok := c[key]
	if ok && getCurrentTime().Before(v.expiry) {
		return errors.New(keyAlreadyPresent)
	}
	c[key] = value{uid: uid, expiry: expiry}
	return nil
}

func (c defaultClientImpl) Get(ctx context.Context, key string) string {
	v, ok := c[key]
	if !ok {
		return ""
	}
	if getCurrentTime().After(v.expiry) {
		c.Delete(ctx, key)
		return ""
	}
	return v.uid
}

func (c defaultClientImpl) Delete(ctx context.Context, key string) error {
	delete(c, key)
	return nil
}
