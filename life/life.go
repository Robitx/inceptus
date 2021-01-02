// Package life provides application wide context with cancel function
// and some helper functions to register syscalls for smooth app termination
// and reloading (conf, logger and such).
//
//
package life

import (
	"context"
	"fmt"
	"os"
	os_signal "os/signal"
	"runtime"
	"syscall"
	"time"
)

// AppContext represents application wide context with cancel function
// and some helper functions to register syscalls for smooth app termination
// and reloading (conf, logger and such)
type AppContext struct {
	Context context.Context
	StopApp context.CancelFunc
}

// ByeTime is a helper that returns true if it's time for application to stop
func (a *AppContext) ByeTime() bool {
	select {
	case <-a.Context.Done():
		return true
	default:
		return false
	}
}

// RunForever keeps main function running
// until there is a panic or kill signal
// NOTE: Call this directly from main, without go!
func (a *AppContext) RunForever() {
	// Memory ballast hack to reduce GC times
	// https://blog.twitch.tv/en/2019/04/10/go-memory-ballast-how-i-learnt-to-stop-worrying-and-love-the-heap-26c2462549a2/
	ballast := make([]byte, 10<<30)
	_ = len(ballast)

	runtime.Goexit()
}

// Reloader interface ..
type Reloader interface {
	Reload()
}

// RegisterReloaders allows you to specify stuff you wan't to "reload" based on some signal
// (like reloading config, or rotating log file based on SIGHUP)
func (a *AppContext) RegisterReloaders(
	signal syscall.Signal, reloaders ...Reloader) {
	listener := make(chan os.Signal)
	os_signal.Notify(listener, signal)

	go func() {
		for {
			<-listener
			for _, r := range reloaders {
				r.Reload()
			}
		}
	}()
}

// RegisterStopSignals handles signals that will terminate application
// After signal is received it cancels app wide context
// and waits specified interval before exiting
func (a *AppContext) RegisterStopSignals(dieTimeout time.Duration,
	signals ...syscall.Signal) {
	listener := make(chan os.Signal)
	for _, signal := range signals {
		os_signal.Notify(listener, signal)
	}
	go func() {
		<-listener
		// _signal := <- listener

		os_signal.Stop(listener)
		a.StopApp()

		fmt.Fprintf(os.Stderr, "\nWaiting %v before exiting.\n", dieTimeout)
		time.Sleep(dieTimeout)
		time.Sleep(200 * time.Millisecond)
		os.Exit(0)
	}()
}

// New prepares application wide context with cancel function
// and some helper functions to register syscalls for smooth app termination
// and reloading (conf, logger and such)
func New() AppContext {
	ctx, stopApp := context.WithCancel(context.Background())
	return AppContext{ctx, stopApp}
}
