package kvutil

import (
	"context"
	"path"
	"strings"
	"sync"

	"github.com/aity-cloud/monty/pkg/storage"
)

type KeyValueStoreLocker[T any] interface {
	storage.KeyValueStoreT[T]
	sync.Locker
}

type kvStoreLockerImpl[T any] struct {
	storage.KeyValueStoreT[T]
	sync.Mutex
}

func NewKeyValueStoreLocker[T any](s storage.KeyValueStoreT[T]) KeyValueStoreLocker[T] {
	return &kvStoreLockerImpl[T]{
		KeyValueStoreT: s,
	}
}

type ValueStoreLocker[T any] interface {
	storage.ValueStoreT[T]
	sync.Locker
}

type valueStoreLockerImpl[T any] struct {
	storage.ValueStoreT[T]
	sync.Locker
}

func NewValueStoreLocker[T any](s storage.ValueStoreT[T], mutex ...sync.Locker) ValueStoreLocker[T] {
	var locker sync.Locker
	if len(mutex) == 0 {
		locker = &sync.Mutex{}
	} else {
		locker = mutex[0]
	}
	return &valueStoreLockerImpl[T]{
		ValueStoreT: s,
		Locker:      locker,
	}
}

type kvStorePrefixImpl[T any] struct {
	base   storage.KeyValueStoreT[T]
	prefix string
}

func (s *kvStorePrefixImpl[T]) Put(ctx context.Context, key string, value T, opts ...storage.PutOpt) error {
	return s.base.Put(ctx, s.prefix+key, value, opts...)
}

func (s *kvStorePrefixImpl[T]) Get(ctx context.Context, key string, opts ...storage.GetOpt) (T, error) {
	return s.base.Get(ctx, s.prefix+key, opts...)
}

func (s *kvStorePrefixImpl[T]) Watch(ctx context.Context, key string, opts ...storage.WatchOpt) (<-chan storage.WatchEvent[storage.KeyRevision[T]], error) {
	c, err := s.base.Watch(ctx, s.prefix+key, opts...)
	if err != nil {
		return nil, err
	}
	out := make(chan storage.WatchEvent[storage.KeyRevision[T]], 1)
	go func() {
		defer close(out)
		for e := range c {
			// TODO: what the / doin
			if e.Current != nil {
				e.Current.SetKey(strings.TrimPrefix("/"+e.Current.Key(), s.prefix))
			}
			if e.Previous != nil {
				e.Previous.SetKey(strings.TrimPrefix("/"+e.Previous.Key(), s.prefix))
			}
			out <- e
		}
	}()
	return out, nil
}

func (s *kvStorePrefixImpl[T]) Delete(ctx context.Context, key string, opts ...storage.DeleteOpt) error {
	return s.base.Delete(ctx, s.prefix+key, opts...)
}

func (s *kvStorePrefixImpl[T]) ListKeys(ctx context.Context, prefix string, opts ...storage.ListOpt) ([]string, error) {
	return s.base.ListKeys(ctx, s.prefix+prefix, opts...)
}

func (s *kvStorePrefixImpl[T]) History(ctx context.Context, key string, opts ...storage.HistoryOpt) ([]storage.KeyRevision[T], error) {
	return s.base.History(ctx, s.prefix+key, opts...)
}

func WithPrefix[T any](base storage.KeyValueStoreT[T], prefix string) storage.KeyValueStoreT[T] {
	return &kvStorePrefixImpl[T]{
		base:   base,
		prefix: prefix,
	}
}

type singleValueStoreImpl[T any] struct {
	base storage.KeyValueStoreT[T]
	key  string
}

func (s *singleValueStoreImpl[T]) Put(ctx context.Context, value T, opts ...storage.PutOpt) error {
	return s.base.Put(ctx, s.key, value, opts...)
}

func (s *singleValueStoreImpl[T]) Get(ctx context.Context, opts ...storage.GetOpt) (T, error) {
	return s.base.Get(ctx, s.key, opts...)
}

func (s *singleValueStoreImpl[T]) Watch(ctx context.Context, opts ...storage.WatchOpt) (<-chan storage.WatchEvent[storage.KeyRevision[T]], error) {
	c, err := s.base.Watch(ctx, s.key, opts...)
	if err != nil {
		return nil, err
	}
	out := make(chan storage.WatchEvent[storage.KeyRevision[T]], 1)
	go func() {
		defer close(out)
		for e := range c {
			if e.Current != nil {
				e.Current.SetKey(path.Base(s.key))
			}
			if e.Previous != nil {
				e.Previous.SetKey(path.Base(s.key))
			}
			out <- e
		}
	}()
	return out, nil
}

func (s *singleValueStoreImpl[T]) Delete(ctx context.Context, opts ...storage.DeleteOpt) error {
	return s.base.Delete(ctx, s.key, opts...)
}

func (s *singleValueStoreImpl[T]) History(ctx context.Context, opts ...storage.HistoryOpt) ([]storage.KeyRevision[T], error) {
	return s.base.History(ctx, s.key, opts...)
}

func WithKey[T any](base storage.KeyValueStoreT[T], key string) storage.ValueStoreT[T] {
	return &singleValueStoreImpl[T]{
		base: base,
		key:  key,
	}
}

type ValueStoreAdapter[T any] struct {
	PutFunc     func(ctx context.Context, value T, opts ...storage.PutOpt) error
	GetFunc     func(ctx context.Context, opts ...storage.GetOpt) (T, error)
	WatchFunc   func(ctx context.Context, opts ...storage.WatchOpt) (<-chan storage.WatchEvent[storage.KeyRevision[T]], error)
	DeleteFunc  func(ctx context.Context, opts ...storage.DeleteOpt) error
	HistoryFunc func(ctx context.Context, opts ...storage.HistoryOpt) ([]storage.KeyRevision[T], error)
}

func (s ValueStoreAdapter[T]) Put(ctx context.Context, value T, opts ...storage.PutOpt) error {
	return s.PutFunc(ctx, value, opts...)
}

func (s ValueStoreAdapter[T]) Get(ctx context.Context, opts ...storage.GetOpt) (T, error) {
	return s.GetFunc(ctx, opts...)
}

func (s ValueStoreAdapter[T]) Watch(ctx context.Context, opts ...storage.WatchOpt) (<-chan storage.WatchEvent[storage.KeyRevision[T]], error) {
	return s.WatchFunc(ctx, opts...)
}

func (s ValueStoreAdapter[T]) Delete(ctx context.Context, opts ...storage.DeleteOpt) error {
	return s.DeleteFunc(ctx, opts...)
}

func (s ValueStoreAdapter[T]) History(ctx context.Context, opts ...storage.HistoryOpt) ([]storage.KeyRevision[T], error) {
	return s.HistoryFunc(ctx, opts...)
}
