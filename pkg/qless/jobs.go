package qless

import "context"

type Jobs struct {
	client Client
}

func (j *Jobs) GetFailedGroups() []string {
	return j.client.redis.SMembers(context.Background(), "ql:failures").Val()
}

func (j *Jobs) GetFailedCounts() map[string]uint64 {
	groups := make(map[string]uint64, 0)

	for _, group := range j.GetFailedGroups() {
		groups[group] = uint64(j.client.redis.LLen(context.Background(), "ql:f:"+group).Val())

	}

	return groups
}

func (j *Jobs) GetFailedJobs() map[string][]*Job {
	jobGroups := make(map[string][]*Job, 0)

	for _, jobGroup := range j.GetFailedGroups() {
		for _, jid := range j.client.redis.LRange(context.Background(), "ql:f:"+jobGroup, 0, -1).Val() {
			jobGroups[jobGroup] = append(jobGroups[jobGroup], &Job{JID: jid, client: j.client})
		}
	}
	return jobGroups
}
