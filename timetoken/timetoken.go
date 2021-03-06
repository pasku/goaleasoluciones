package timetoken

import (
	"sync"
	"time"

	"github.com/aleasoluciones/goaleasoluciones/clock"
	"github.com/aleasoluciones/goaleasoluciones/scheduledtask"
)

type tokenManager struct {
	periode      time.Duration
	ttl          time.Duration
	periodicFunc func(id string)
	tokens       map[string]time.Time
	mutex        sync.Mutex
	clock        clock.Clock
}

func NewTokenManager(periode, ttl time.Duration, periodicFunc func(id string)) *tokenManager {
	tm := tokenManager{
		periode:      periode,
		ttl:          ttl,
		periodicFunc: periodicFunc,
		tokens:       make(map[string]time.Time),
		clock:        clock.NewClock(),
	}

	scheduledtask.NewScheduledTask(
		tm.executePeriodicFunc,
		tm.periode,
		tm.ttl)
	return &tm
}

func (tm tokenManager) executePeriodicFunc() {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	now := tm.clock.Now()

	for k, v := range tm.tokens {
		if v.Sub(now) >= 0 {
			go tm.periodicFunc(k)
		} else {
			delete(tm.tokens, k)
		}
	}
}

func (tm *tokenManager) Add(id string, ttl time.Duration) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	newExpirationTime := tm.clock.Now().Add(ttl)

	actualExpirationTime, found := tm.tokens[id]
	if !found {
		tm.tokens[id] = newExpirationTime
		return
	}

	if newExpirationTime.Sub(actualExpirationTime) > 0 {
		tm.tokens[id] = newExpirationTime
	}
}
