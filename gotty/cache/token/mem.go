package token

import (
	"log"
	"time"

	"github.com/patrickmn/go-cache"
)

// MemCache use memory to store token and TtyParameter
type MemCache struct {
	cache *cache.Cache
}

// NewMemCache new MemCache
func NewMemCache() *MemCache {
	return &MemCache{
		cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

// Get token param from memory
func (r *MemCache) Get(token string) *TtyParameter {
	obj, exit := r.cache.Get(token)
	if !exit {
		return nil
	}
	param, ok := obj.(TtyParameter)
	if ok {
		return &param
	}
	log.Printf("get token %s from mem obj is not tty param", token)
	return nil
}

// Delete token from memory
func (r *MemCache) Delete(token string) error {
	r.cache.Delete(token)
	return nil
}

// Add token to memory
func (r *MemCache) Add(token string, param *TtyParameter, d time.Duration) error {
	return r.cache.Add(token, *param, d)
}

// GetSession retrieves session info from memory
func (r *MemCache) GetSession(sessionToken string) *SessionInfo {
	obj, exists := r.cache.Get("session:" + sessionToken)
	if !exists {
		return nil
	}
	session, ok := obj.(SessionInfo)
	if ok {
		return &session
	}
	log.Printf("get session %s from mem obj is not session info", sessionToken)
	return nil
}

// SetSession stores session info to memory
func (r *MemCache) SetSession(sessionToken string, session *SessionInfo, d time.Duration) error {
	return r.cache.Add("session:"+sessionToken, *session, d)
}

// UpdateSessionAccess updates session last access time
func (r *MemCache) UpdateSessionAccess(sessionToken string) error {
	session := r.GetSession(sessionToken)
	if session == nil {
		return nil
	}
	session.LastAccess = time.Now().Unix()
	r.cache.Set("session:"+sessionToken, *session, cache.DefaultExpiration)
	return nil
}

// DeleteSession removes session from memory
func (r *MemCache) DeleteSession(sessionToken string) error {
	r.cache.Delete("session:" + sessionToken)
	return nil
}
