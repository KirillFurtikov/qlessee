package qless

import "context"

// Config store qless configuration data
type Config struct {
	client Client
}

var config *Config

// Get value for key
func (c *Config) Get(key string) string {
	return c.client.redis.HGet(context.Background(), "ql:config", key).Val()
}

// GetAll return all values in ql:config
func (c *Config) GetAll() map[string]string {
	return c.client.redis.HGetAll(context.Background(), "ql:config").Val()
}

// Set value
func (c *Config) Set(option, value string) int64 {
	result := c.client.redis.HSet(context.Background(), "ql:config", option, value).Val()
	return result
}

// Unset value
func (c *Config) Unset(option string) int64 {
	result := c.client.redis.HDel(context.Background(), "ql:config", option).Val()
	return result
}
