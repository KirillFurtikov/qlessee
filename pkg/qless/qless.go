package qless

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// Options stores properties
type Options struct {
	Name string
	URL  string
}

// Client Qless instance
type Client struct {
	Name   string
	URL    string
	redis  *redis.Client
	Queues []Queue
	Config *Config
}

var client *Client

func (c *Client) Redis() *redis.Client {
	return c.redis
}

// NewClient initialize Qless instance by redis
func NewClient(options *Options) *Client {
	rdo, err := redis.ParseURL(options.URL)
	if err != nil {
		panic(err)
	}

	client = &Client{
		Name:  options.Name,
		URL:   options.URL,
		redis: redis.NewClient(rdo),
	}
	client.Config = &Config{client: *client}

	return client
}

// LoadQueues - load data from redis and write it in Qless properties
func (c *Client) LoadQueues() *[]Queue {
	// ZRANGE ql:queues 0 -1
	response, err := c.redis.ZRange(context.Background(), "ql:queues", 0, -1).Result()
	if err != nil {
		panic(err)
	}

	for _, name := range response {
		c.Queues = append(c.Queues, Queue{Name: name, client: *c})
	}

	return &c.Queues
}

// GetQueue - get Queue by name
func (c *Client) GetQueue(name string) *Queue {
	queue := new(Queue)
	for _, element := range c.Queues {
		if element.Name == name {
			*queue = element
		}
	}
	return queue
}

func (c *Client) GetJob(jid string) *Job {
	j := &Job{JID: jid, client: *c}
	j.Load()
	return j
}

func (c *Client) Jobs() *Jobs {
	return &Jobs{client: *c}
}

// Close - close connection with redis
func (c *Client) Close() {
	c.redis.Close()
}
