package synchronizer

import (
	"context"
	"sync"
	"time"
)

type Session struct {
	image      string
	pullPolicy string
	timeout    time.Duration    // this is the session timeout, not the caller timeout
	done       chan interface{} // chan will block until result
	isComplete bool
	isTimedOut bool
	err        error
}

type synchronizer struct {
	sessionMap      map[string]Session
	mutex           sync.Mutex
	timeout         time.Duration
	sessions        chan Session
	completedEvents chan string
	ctx             context.Context
}

// allows mocking/dependency injection
type ISychronizer interface {
	// returns session that is ready to wait on, or error
	StartPull(image, pullPolicy string) (Session, error)
	// waits for session to time out or succeed
	WaitForPull(session Session, timeout time.Duration) error
}
