package qless

import (
	"context"
	"encoding/json"
	"fmt"

	jsonc "github.com/nwidger/jsoncolor"
)

type Job struct {
	JID            string                 `redis:"jid"`
	Klass          string                 `redis:"klass"`
	State          string                 `redis:"state"`
	Queue          string                 `redis:"queue"`
	Worker         string                 `redis:"worker"`
	Data           map[string]interface{} `redis:"data"`
	Tags           string                 `redis:"tags"`
	Failure        map[string]interface{} `redis:"failure"`
	SpawnedFromJid string                 `redis:"spawned_from_jid"`
	Time           string                 `redis:"time"`
	Dependents     []string
	Dependencies   []string
	Tracked        float64
	Priority       int `redis:"priority"`
	Expires        int `redis:"expires"`
	Retries        int `redis:"retries"`
	Remaining      int `redis:"remaining"`
	client         Client
}

func (j *Job) Load() *Job {
	fields := []string{
		"jid", "klass", "state", "queue", "worker", "tags", "spawned_from_jid", "time", "priority", "expires", "retries", "remaining",
	}

	j.client.redis.HMGet(context.Background(), "ql:j:"+j.JID, fields...).Scan(j)
	data, _ := j.client.redis.HGet(context.Background(), "ql:j:"+j.JID, "data").Bytes()
	json.Unmarshal(data, &j.Data)

	failure, _ := j.client.redis.HGet(context.Background(), "ql:j:"+j.JID, "failure").Bytes()
	json.Unmarshal(failure, &j.Failure)

	return j
}

func (j *Job) updateTracked() *Job {
	j.Tracked = j.client.redis.ZScore(context.Background(), "ql:tracked", j.JID).Val()
	return j
}

func (j *Job) updateDependents() *Job {
	j.Dependents = j.client.redis.SMembers(context.Background(), "ql:tracked"+j.JID+"-dependents").Val()
	return j
}

func (j *Job) updateDependencies() *Job {
	j.Dependencies = j.client.redis.SMembers(context.Background(), "ql:tracked"+j.JID+"-dependencies").Val()
	return j
}

func (j *Job) Pretty() []byte {
	b, err := jsonc.MarshalIndent(j, "", "    ")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return b
}
