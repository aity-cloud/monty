package waitctx

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/ttacon/chalk"
	"go.uber.org/atomic"
)

type waitCtxDataKeyType struct{}
type timeout struct{}

var waitCtxDataKey waitCtxDataKeyType
var waitTimeout timeout

type waitCtxData struct {
	wg      sync.WaitGroup
	waiting *atomic.Bool
	cond    *sync.Cond
}

func FromContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, waitCtxDataKey, &waitCtxData{
		waiting: atomic.NewBool(false),
		cond:    sync.NewCond(&sync.Mutex{}),
	})
}

func Background() context.Context {
	return FromContext(context.Background())
}

type RestrictiveContext = context.Context
type PermissiveContext = context.Context

type WaitContextInterface interface {
	AddOne(ctx context.Context)
	Done(ctx context.Context)
	Wait(ctx context.Context, notifyAfter ...time.Duration)

	// Go runs the given function in a background goroutine, and adds it to the
	// WaitContext. Shorthand for the following pattern:
	//  waitctx.AddOne(ctx)
	//  go func() {
	//    defer waitctx.Done(ctx)
	//    // do stuff
	//  }()
	Go(ctx context.Context, fn func())
}

type restrictive struct{}

func (restrictive) FromContext(ctx RestrictiveContext) RestrictiveContext {
	return context.WithValue(ctx, waitCtxDataKey, &waitCtxData{
		wg: sync.WaitGroup{},
	})
}

func (restrictive) AddOne(ctx RestrictiveContext) {
	data := ctx.Value(waitCtxDataKey)
	if data == nil {
		panic("context is not a WaitContext")
	}
	d := data.(*waitCtxData)
	d.cond.L.Lock()
	for d.waiting.Load() {
		d.cond.Wait()
	}
	d.wg.Add(1)
	d.cond.L.Unlock()
}

func (restrictive) Done(ctx RestrictiveContext) {
	data := ctx.Value(waitCtxDataKey)
	if data == nil {
		panic("context is not a WaitContext")
	}
	data.(*waitCtxData).wg.Done()
}

func (restrictive) Wait(ctx RestrictiveContext, notifyAfter ...time.Duration) {
	data := ctx.Value(waitCtxDataKey)
	if data == nil {
		panic("context is not a WaitContext")
	}
	done := make(chan struct{})
	d := data.(*waitCtxData)
	go func() {
		defer close(done)
		d.cond.L.Lock()
		d.waiting.Store(true)
		d.wg.Wait()
		d.waiting.Store(false)
		d.cond.Broadcast()
		d.cond.L.Unlock()
	}()
	if len(notifyAfter) > 0 {
		stack := string(debug.Stack())
		go func(duration time.Duration) {
			for {
				select {
				case <-done:
					return
				case <-time.After(duration):
					fmt.Fprint(os.Stderr, chalk.Yellow.Color("\n=== WARNING: waiting longer than expected for context to cancel ===\n"+stack+"\n"))
				}
			}
		}(notifyAfter[0])
	}
	<-done
}

// WaitWithTimeout follows the same pattern as wait, but with a timeout.  If the timeout expires
// WaitWithTimeout will panic.
func (restrictive) WaitWithTimeout(ctx RestrictiveContext, timeout time.Duration, notifyAfter ...time.Duration) {
	data := ctx.Value(waitCtxDataKey)
	if data == nil {
		panic("context is not a WaitContext")
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		d := data.(*waitCtxData)
		d.cond.L.Lock()
		d.waiting.Store(true)
		d.wg.Wait()
		d.waiting.Store(false)
		d.cond.Broadcast()
		d.cond.L.Unlock()
	}()
	if len(notifyAfter) > 0 {
		stack := string(debug.Stack())
		notifyDone := make(chan struct{})
		defer close(notifyDone)
		go func() {
			for {
				select {
				case <-notifyDone:
					return
				case <-time.After(notifyAfter[0]):
					fmt.Fprint(os.Stderr, chalk.Yellow.Color("\n=== WARNING: waiting longer than expected for context to cancel ===\n"+stack+"\n"))
				}
			}
		}()
	}
	select {
	case <-done:
		return
	case <-time.After(timeout):
		panic(waitTimeout)
	}
}

func (w restrictive) Go(ctx RestrictiveContext, fn func()) {
	w.AddOne(ctx)
	go func() {
		defer w.Done(ctx)
		fn()
	}()
}

type permissive struct{}

func (permissive) AddOne(ctx PermissiveContext) {
	if data := ctx.Value(waitCtxDataKey); data != nil {
		d := data.(*waitCtxData)
		d.cond.L.Lock()
		for d.waiting.Load() {
			d.cond.Wait()
		}
		d.wg.Add(1)
		d.cond.L.Unlock()
	}
}

func (permissive) Done(ctx PermissiveContext) {
	if data := ctx.Value(waitCtxDataKey); data != nil {
		data.(*waitCtxData).wg.Done()
	}
}

func (permissive) Wait(ctx PermissiveContext, notifyAfter ...time.Duration) {
	data := ctx.Value(waitCtxDataKey)
	if data == nil {
		return
	}
	done := make(chan struct{})
	d := data.(*waitCtxData)
	go func() {
		defer close(done)
		d.cond.L.Lock()
		d.waiting.Store(true)
		d.wg.Wait()
		d.waiting.Store(false)
		d.cond.Broadcast()
		d.cond.L.Unlock()
	}()
	if len(notifyAfter) > 0 {
		stack := string(debug.Stack())
		go func(d time.Duration) {
			for {
				select {
				case <-done:
					return
				case <-time.After(d):
					fmt.Fprint(os.Stderr, chalk.Yellow.Color("\n=== WARNING: waiting longer than expected for context to cancel ===\n"+stack+"\n"))
				}
			}
		}(notifyAfter[0])
	}
	<-done
}

func (w permissive) Go(ctx PermissiveContext, fn func()) {
	w.AddOne(ctx)
	go func() {
		defer w.Done(ctx)
		fn()
	}()
}

func RecoverTimeout() {
	r := recover()
	if r == nil {
		return
	}
	switch r.(type) {
	case timeout:
		fmt.Println("timeout expired, exiting process")
		os.Exit(0)
	default:
		fmt.Println(r)
		fmt.Println(string(debug.Stack()))
		os.Exit(1)
	}
}

var (
	Restrictive = restrictive{}
	Permissive  = permissive{}

	AddOne          = Restrictive.AddOne
	Done            = Restrictive.Done
	Wait            = Restrictive.Wait
	WaitWithTimeout = Restrictive.WaitWithTimeout
	Go              = Restrictive.Go
)
