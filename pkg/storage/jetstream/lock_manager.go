package jetstream

import (
	"context"
	"log/slog"

	"github.com/aity-cloud/monty/pkg/storage"
	"github.com/aity-cloud/monty/pkg/storage/lock"
	"github.com/nats-io/nats.go"
)

// Requires jetstream 2.9+
type LockManager struct {
	ctx context.Context
	js  nats.JetStreamContext

	lg *slog.Logger

	prefix string
}

func NewLockManager(ctx context.Context, js nats.JetStreamContext, prefix string, lg *slog.Logger) *LockManager {
	prefix = sanitizePrefix(prefix)
	return &LockManager{
		ctx:    ctx,
		js:     js,
		lg:     lg,
		prefix: prefix,
	}
}

var _ storage.LockManager = (*LockManager)(nil)

func (l *LockManager) NewLock(key string, opts ...lock.LockOption) storage.Lock {
	options := lock.DefaultLockOptions()
	options.Apply(opts...)
	return NewLock(l.js, l.prefix, key, l.lg, options)
}
