package token

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache use redis to store token and TtyParameter
type RedisCache struct {
	client *redis.Client
	prefix string
}

// NewRedisCache new redis token cache
func NewRedisCache(client *redis.Client, prefix string) *RedisCache {
	return &RedisCache{
		client: client,
		prefix: prefix,
	}
}

// Get token param from redis
func (r *RedisCache) Get(token string) *TtyParameter {
	b, err := r.client.Get(context.Background(), r.prefix+token).Bytes()
	if err != nil {
		log.Printf("get token %s from redis fatal:%s", token, err.Error())
		return nil
	}
	var param TtyParameter
	if err := json.Unmarshal(b, &param); err != nil {
		log.Printf("get token param %s unmarshal err:%s", string(b), err.Error())
		return nil
	}
	return &param
}

// Delete token from redis
func (r *RedisCache) Delete(token string) error {
	return r.client.Del(context.Background(), r.prefix+token).Err()
}

// Add token to redis
func (r *RedisCache) Add(token string, param *TtyParameter, d time.Duration) error {
	b, err := json.Marshal(param)
	if err != nil {
		return err
	}
	return r.client.Set(context.Background(), r.prefix+token, string(b), d).Err()
}

// GetSession retrieves session info from redis
func (r *RedisCache) GetSession(sessionToken string) *SessionInfo {
	b, err := r.client.Get(context.Background(), r.prefix+"session:"+sessionToken).Bytes()
	if err != nil {
		log.Printf("get session %s from redis error:%s", sessionToken, err.Error())
		return nil
	}
	var session SessionInfo
	if err := json.Unmarshal(b, &session); err != nil {
		log.Printf("get session info %s unmarshal error:%s", string(b), err.Error())
		return nil
	}
	return &session
}

// SetSession stores session info to redis
func (r *RedisCache) SetSession(sessionToken string, session *SessionInfo, d time.Duration) error {
	b, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return r.client.Set(context.Background(), r.prefix+"session:"+sessionToken, string(b), d).Err()
}

// UpdateSessionAccess updates session last access time
func (r *RedisCache) UpdateSessionAccess(sessionToken string) error {
	session := r.GetSession(sessionToken)
	if session == nil {
		return nil
	}
	session.LastAccess = time.Now().Unix()
	return r.SetSession(sessionToken, session, 0) // 使用原有的TTL
}

// DeleteSession removes session from redis
func (r *RedisCache) DeleteSession(sessionToken string) error {
	return r.client.Del(context.Background(), r.prefix+"session:"+sessionToken).Err()
}
