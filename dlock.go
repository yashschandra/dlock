package dlock

import (
	"context"
	"errors"
	"time"
)

const (
	lockNotSet     = "lock not set"
	lockAlreadySet = "lock already set"
)

type Lock struct {
	ctx    context.Context
	key    string
	uid    string
	locked bool
	expiry time.Time
	store  Store
}

func New(ctx context.Context, store Store) *Lock {
	if store == nil {
		store = getDefaultClient()
	}
	return &Lock{ctx: ctx, store: store, uid: newUUID()}
}

func (l *Lock) set(key string, expiry time.Time) {
	l.locked = true
	l.expiry = expiry
	l.key = key
}

func (l *Lock) unset() {
	l.locked = false
}

func (l *Lock) isSet() bool {
	if !l.locked {
		return false
	}
	if l.store.Get(l.ctx, l.key) != l.uid {
		l.unset()
		return false
	}
	if getCurrentTime().After(l.expiry) {
		l.unset()
		return false
	}
	return true
}

func (l *Lock) take(key string, ttl int64) error {
	if l.isSet() {
		return errors.New(lockAlreadySet)
	}
	expiry := getCurrentTime().Add(time.Duration(ttl) * time.Second)
	err := l.store.Set(l.ctx, key, l.uid, expiry)
	if err != nil {
		return err
	}
	l.set(key, expiry)
	return nil
}

func (l *Lock) release() error {
	if !l.isSet() {
		return errors.New(lockNotSet)
	}
	err := l.store.Delete(l.ctx, l.key)
	if err != nil {
		return err
	}
	l.unset()
	return nil
}

func (l *Lock) Take(key string, ttl int64) error {
	return l.take(key, ttl)
}

func (l *Lock) Release() error {
	return l.release()
}
