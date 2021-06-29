package qless_test

import (
	"context"

	"github.com/KirillFurtikov/qlessee/pkg/qless"
	"github.com/alicebob/miniredis"
	redis "github.com/go-redis/redis/v8"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var ctx = context.Background()
var qlessClient *qless.Client
var redisMock *redis.Client
var rawJob = []string{
	"jid", "testjid",
	"failure", "{\"foo\": \"bar\"}",
	"expires", "123",
	"klass", "JobTestClass",
	"priority", "1",
	"retries", "5",
	"data", "{\"hey\": \"bruh\"}",
	"worker", "myWorker",
	"state", "completed",
	"time", "55435345435",
	"remaining", "777",
	"queue", "system",
	"tags", "iwatch",
}

var expectedJob = &qless.Job{
	JID:       "testjid",
	Failure:   map[string]interface{}{"foo": "bar"},
	Expires:   123,
	Klass:     "JobTestClass",
	Priority:  1,
	Retries:   5,
	Data:      map[string]interface{}{"hey": "bruh"},
	Worker:    "myWorker",
	State:     "completed",
	Time:      "55435345435",
	Remaining: 777,
	Queue:     "system",
	Tags:      "iwatch",
}

var _ = Describe("Qless", func() {
	r, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	qlessClient = qless.NewClient(&qless.Options{
		Name: "api",
		URL:  "redis://" + r.Addr(),
	})

	redisMock = qlessClient.Redis()

	BeforeEach(func() {
		redisMock.FlushAll(ctx)

		for i, value := range []string{"system", "test"} {
			redisMock.ZAdd(ctx, "ql:queues", &redis.Z{float64(i), value})
		}

		redisMock.HSet(ctx, "ql:f:"+rawJob[1], rawJob)

	})
	Describe("Queues", func() {
		It("should load queues", func() {
			qlessClient.LoadQueues()

			Expect(len(qlessClient.Queues)).
				To(Equal(len(redisMock.ZRange(context.Background(), "ql:queues", 0, -1).Val())))
		})

		It("should get single queue by name", func() {
			qlessClient.LoadQueues()

			Expect(qlessClient.GetQueue("system").Name).To(Equal("system"))
		})

		It("should get counts", func() {
			qlessClient.LoadQueues()

			Expect(qlessClient.Queues[0].Counts()).
				To(Equal(map[string]int64{"depends": 0, "scheduled": 0, "recurring": 0, "stalled": 0, "work": 0}))
		})

		It("should get failed count", func() {
			qlessClient.LoadQueues()
			failedGroups := redisMock.SMembers(context.Background(), "ql:failures").Val()
			expectedCount := int64(0)

			for _, group := range failedGroups {
				expectedCount += redisMock.LLen(context.Background(), "ql:f:"+group).Val()
			}

			Expect(qlessClient.GetQueue("system").FailedCount()).
				To(Equal(expectedCount))
		})

		It("Should pause/continue queues", func() {
			qlessClient.LoadQueues()

			Expect(qlessClient.Queues[0].Paused()).
				To(Equal(false))
			Expect((qlessClient.Queues)[0].Pause().Val()).
				To(Equal(int64(1)))

			Expect(qlessClient.Queues[0].Paused()).
				To(Equal(true))
			Expect((qlessClient.Queues)[0].Continue().Val()).
				To(Equal(int64(1)))

			Expect(qlessClient.Queues[0].Paused()).
				To(Equal(false))
		})
	})

	Describe("Config", func() {
		It("should set/get value by key", func() {
			qlessClient.Config.Set("foo", "bar")
			Expect(qlessClient.Config.Get("foo")).To(Equal("bar"))
		})

		It("should unset key", func() {
			qlessClient.Config.Set("foo", "bar")
			Expect(qlessClient.Config.Unset("foo")).To(Equal(int64(1)))
		})
	})

	Describe("Job", func() {
		It("should get Failed job groups", func() {
			qlessClient.LoadQueues()

			Expect(qlessClient.Jobs().GetFailedGroups()).
				To(Equal(redisMock.SMembers(context.Background(), "ql:failures").Val()))
		})

		It("should get Failed job counts", func() {
			qlessClient.LoadQueues()
			groups := make(map[string]uint64, 0)
			for _, group := range redisMock.SMembers(context.Background(), "ql:failures").Val() {
				groups[group] = uint64(redisMock.LLen(context.Background(), "ql:f:"+group).Val())

			}
			Expect(qlessClient.Jobs().GetFailedCounts()).To(Equal(groups))
		})
	})
})
