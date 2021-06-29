package qless

import (
	"context"
	"strings"

	"github.com/go-redis/redis/v8"
)

/*

https://github.com/seomoz/qless-core/blob/94076f350b45bfe9d423ac82262c5ba4b7cc4ac9/queue.lua#L965

stalled - locks.lenght(now)
waiting - work.lenght
running - (locks.length() -  stalled)
scheduled - scheduled.lenght()
depends - depends.length()
recurring - recurring.length()
paused - Queue.Paused()

*/
var postfixForStatus = map[string]string{
	"stalled":   "-locks",
	"work":      "-work",
	"depends":   "-depends",
	"scheduled": "-scheduled",
	"recurring": "-recur",
}
var _ string

const qlPausedQueues = "ql:paused_queues"

// Queue - store properties
type Queue struct {
	client Client
	Name   string
	Stats  map[string]string
}

// Pause - pause queue
func (q *Queue) Pause() *redis.IntCmd {
	result := q.client.redis.SAdd(context.Background(), qlPausedQueues, q.Name)

	return result
}

// Continue - unpause queue
func (q *Queue) Continue() *redis.IntCmd {
	return q.client.redis.SRem(context.Background(), qlPausedQueues, q.Name)
}

// Paused - check queue pause status
func (q *Queue) Paused() bool {
	result, err := q.client.redis.SIsMember(context.Background(), qlPausedQueues, q.Name).Result()

	if err != nil {
		panic(err)
	}

	return result
}

// Counts - count jobs by states in redis
func (q *Queue) Counts() map[string]int64 {
	var result = map[string]int64{}

	for k := range postfixForStatus {
		result[k] = q.client.redis.ZCard(context.Background(), q.getKeyForStatus(k)).Val()
	}

	return result
}

// Faild return failed jobs groups
func (q *Queue) Failed() []string {
	return q.client.redis.SMembers(context.Background(), "ql:failures").Val()

}

// FaildCount return failed count jobs
func (q *Queue) FailedCount() int64 {
	count := int64(0)
	for _, group := range q.Failed() {
		count = count + q.client.redis.LLen(context.Background(), "ql:f:"+group).Val()
	}
	return count
}

func (q *Queue) getKeyForStatus(status string) string {
	str := []string{"ql:q:", q.Name, postfixForStatus[status]}

	return strings.Join(str, "")
}
