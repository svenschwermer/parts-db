package session

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

type Manager struct {
	sync.Mutex
	sessions        map[string]time.Time
	expirationTimer *time.Timer
}

func NewManager() *Manager {
	m := &Manager{sessions: make(map[string]time.Time)}
	m.expirationTimer = time.AfterFunc(time.Hour, func() {
		m.Lock()
		defer m.Unlock()

		for id, expires := range m.sessions {
			if expires.Before(time.Now()) {
				delete(m.sessions, id)
			}
		}
	})
	return m
}

func (m *Manager) Valid(id string) bool {
	m.Lock()
	defer m.Unlock()

	expires, ok := m.sessions[id]
	return ok && expires.After(time.Now())
}

func (m *Manager) Delete(id string) {
	m.Lock()
	defer m.Unlock()

	delete(m.sessions, id)
	m.expire()
}

func (m *Manager) New(valid time.Duration) string {
	rawID := make([]byte, 64)
	_, err := rand.Read(rawID)
	if err != nil {
		panic(err)
	}
	id := base64.URLEncoding.EncodeToString(rawID)

	m.Lock()
	defer m.Unlock()

	m.sessions[id] = time.Now().Add(valid)
	m.expire()
	return id
}

func (m *Manager) expire() {
	var next time.Time
	for _, expires := range m.sessions {
		if (next == time.Time{}) || expires.Before(next) {
			next = expires
		}
	}

	if (next == time.Time{}) {
		m.expirationTimer.Stop()
	} else {
		m.expirationTimer.Reset(time.Until(next))
	}
}
