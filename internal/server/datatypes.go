package server

import (
	"time"
)

// RedisValue represents different Redis data types
type RedisValue struct {
	Type      string                 // "string", "list", "set", "hash", "zset"
	String    string                 // For string values
	List      []string               // For list values
	Set       map[string]bool        // For set values (using map for O(1) lookup)
	Hash      map[string]string      // For hash values
	ZSet      map[string]float64     // For sorted set values (member -> score)
	ExpiresAt *time.Time             // For TTL support
}

// NewStringValue creates a new string value
func NewStringValue(s string) *RedisValue {
	return &RedisValue{
		Type:   "string",
		String: s,
	}
}

// NewListValue creates a new list value
func NewListValue() *RedisValue {
	return &RedisValue{
		Type: "list",
		List: make([]string, 0),
	}
}

// NewSetValue creates a new set value
func NewSetValue() *RedisValue {
	return &RedisValue{
		Type: "set",
		Set:  make(map[string]bool),
	}
}

// NewHashValue creates a new hash value
func NewHashValue() *RedisValue {
	return &RedisValue{
		Type: "hash",
		Hash: make(map[string]string),
	}
}

// IsExpired checks if the value has expired
func (rv *RedisValue) IsExpired() bool {
	if rv.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*rv.ExpiresAt)
}

// SetExpiration sets an expiration time for the value
func (rv *RedisValue) SetExpiration(ttl time.Duration) {
	expireTime := time.Now().Add(ttl)
	rv.ExpiresAt = &expireTime
}

// List operations
func (rv *RedisValue) ListPush(value string, left bool) int {
	if rv.Type != "list" {
		return -1
	}
	if left {
		rv.List = append([]string{value}, rv.List...)
	} else {
		rv.List = append(rv.List, value)
	}
	return len(rv.List)
}

func (rv *RedisValue) ListPop(left bool) (string, bool) {
	if rv.Type != "list" || len(rv.List) == 0 {
		return "", false
	}
	
	var value string
	if left {
		value = rv.List[0]
		rv.List = rv.List[1:]
	} else {
		value = rv.List[len(rv.List)-1]
		rv.List = rv.List[:len(rv.List)-1]
	}
	return value, true
}

func (rv *RedisValue) ListLength() int {
	if rv.Type != "list" {
		return 0
	}
	return len(rv.List)
}

// Set operations
func (rv *RedisValue) SetAdd(member string) bool {
	if rv.Type != "set" {
		return false
	}
	_, exists := rv.Set[member]
	rv.Set[member] = true
	return !exists // Return true if it's a new member
}

func (rv *RedisValue) SetRemove(member string) bool {
	if rv.Type != "set" {
		return false
	}
	_, exists := rv.Set[member]
	delete(rv.Set, member)
	return exists
}

func (rv *RedisValue) SetContains(member string) bool {
	if rv.Type != "set" {
		return false
	}
	_, exists := rv.Set[member]
	return exists
}

func (rv *RedisValue) SetMembers() []string {
	if rv.Type != "set" {
		return nil
	}
	members := make([]string, 0, len(rv.Set))
	for member := range rv.Set {
		members = append(members, member)
	}
	return members
}

// Hash operations
func (rv *RedisValue) HashSet(field, value string) bool {
	if rv.Type != "hash" {
		return false
	}
	_, exists := rv.Hash[field]
	rv.Hash[field] = value
	return !exists // Return true if it's a new field
}

func (rv *RedisValue) HashGet(field string) (string, bool) {
	if rv.Type != "hash" {
		return "", false
	}
	value, exists := rv.Hash[field]
	return value, exists
}

func (rv *RedisValue) HashDelete(field string) bool {
	if rv.Type != "hash" {
		return false
	}
	_, exists := rv.Hash[field]
	delete(rv.Hash, field)
	return exists
}

func (rv *RedisValue) HashGetAll() map[string]string {
	if rv.Type != "hash" {
		return nil
	}
	result := make(map[string]string)
	for k, v := range rv.Hash {
		result[k] = v
	}
	return result
}
