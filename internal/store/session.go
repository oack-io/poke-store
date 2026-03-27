package store

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"sync"

	"github.com/oack-io/poke-store/internal/model"
)

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]model.User
	secret   []byte
}

func NewSessionStore(secret string) *SessionStore {
	return &SessionStore{
		sessions: make(map[string]model.User),
		secret:   []byte(secret),
	}
}

func (s *SessionStore) Create(user model.User) string {
	token := s.generateToken()
	s.mu.Lock()
	s.sessions[token] = user
	s.mu.Unlock()
	return token
}

func (s *SessionStore) Get(token string) (model.User, bool) {
	s.mu.RLock()
	user, ok := s.sessions[token]
	s.mu.RUnlock()
	return user, ok
}

func (s *SessionStore) Delete(token string) {
	s.mu.Lock()
	delete(s.sessions, token)
	s.mu.Unlock()
}

func (s *SessionStore) generateToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	mac := hmac.New(sha256.New, s.secret)
	mac.Write(b)
	return hex.EncodeToString(mac.Sum(nil))
}
