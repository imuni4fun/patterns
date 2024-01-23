package synchronizer

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// channels are exposed for testing
func GetSynchronizer(
	ctx context.Context,
	timeout time.Duration,
	sessionChan chan Session,
	completedChan chan string,
) ISychronizer {
	if cap(sessionChan) < 50 {
		panic("session channel must have capacity to buffer events, minimum of 50 is required")
	}
	if cap(completedChan) < 5 {
		panic("completion channel must have capacity to buffer events, minimum of 5 is required")
	}
	obj := synchronizer{
		sessionMap:      make(map[string]Session),
		mutex:           sync.Mutex{},
		timeout:         timeout,
		sessions:        sessionChan,
		completedEvents: completedChan,
		ctx:             ctx,
	}
	return &obj
}

func (s *synchronizer) StartPull(image, pullPolicy string) (Session, error) {
	fmt.Printf("start pull: asked to pull image %s with policy %s\n", image, pullPolicy)
	s.mutex.Lock() // lock mutex
	defer s.mutex.Unlock()
	ses, ok := s.sessionMap[image] // try get session
	if !ok {                       // if no session, create session
		ses = Session{
			image:      image,
			pullPolicy: pullPolicy,
			timeout:    s.timeout,
			done:       make(chan interface{}),
			isComplete: false,
			isTimedOut: false,
			err:        nil,
		}
		select {
		case s.sessions <- ses: // start session, this could deadlock so must think this through
			fmt.Printf("start pull: new session created for %s with policy %s and timeout %v\n", ses.image, ses.pullPolicy, ses.timeout)
		default: // catch deadlock
			panic("start pull: announce deadlock")
		}
		s.sessionMap[image] = ses // add session to map
	} else {
		fmt.Printf("start pull: found open session for %s\n", ses.image)
	}
	// return session and unlock
	return ses, nil
}

func (s *synchronizer) WaitForPull(session Session, callerTimeout time.Duration) error {
	fmt.Printf("wait for pull: starting to wait for image %s\n", session.image)
	defer fmt.Printf("wait for pull: exiting wait for image %s\n", session.image)
	select {
	case <-session.done: // success or error (including session timeout)
		fmt.Printf("wait for pull: pull completed for %s, isError: %t\n",
			session.image, session.err != nil)
		return session.err
	case <-time.After(callerTimeout):
		return fmt.Errorf("wait for pull: this wait for image %s has timed out after %v",
			session.image, callerTimeout)
	case <-s.ctx.Done(): //TODO: might wait for puller to do this instead
		return fmt.Errorf("wait for pull: synchronizer is shutting down") // must return error since not success
	}
}

func (s *synchronizer) RunCompletionsChecker() {
	go func() {
		for {
			select {
			case <-s.ctx.Done(): // shut down loop
				close(s.sessions) // the writer is supposed to close channels
				s.mutex.Lock()
				for image := range s.sessionMap {
					delete(s.sessionMap, image) // no-op if already deleted due to race
				}
				s.mutex.Unlock()
			case image := <-s.completedEvents: // remove session (no longer active)
				s.mutex.Lock()
				delete(s.sessionMap, image) // no-op if already deleted due to race
				s.mutex.Unlock()
			}
		}
	}()
}
