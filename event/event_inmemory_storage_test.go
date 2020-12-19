package event

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInmemoryQueryMetricStorage_Upsert(t *testing.T) {
	storage := NewQueryEventInmemoryStorage()

	startAt, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	queryStart := &Query{
		Command:      "find",
		Database:     "test",
		RequestID:    1,
		ConnectionID: "1",
		Duration:     0,
		StartAt:      startAt,
		CompletedAt:  startAt,
	}

	storage.Upsert(queryStart)

	queryEnd := &Query{
		Command:      "find",
		Database:     "",
		RequestID:    1,
		ConnectionID: "1",
		Duration:     5 * time.Second,
		StartAt:      time.Now(),
		CompletedAt:  startAt.Add(5 * time.Second),
	}

	storage.Upsert(queryEnd)

	query, _ := storage.Extract(queryEnd.Key())

	assert.Equal(t, queryStart.Command, query.Command)
	assert.Equal(t, queryStart.Database, query.Database)
	assert.Equal(t, queryStart.ConnectionID, query.ConnectionID)
	assert.Equal(t, queryEnd.Duration, query.Duration)
	assert.Equal(t, queryStart.StartAt, query.StartAt)
	assert.Equal(t, queryEnd.CompletedAt, query.CompletedAt)
}

func TestInmemoryQueryMetricStorage_Extract(t *testing.T) {
	storage := NewQueryEventInmemoryStorage()

	startAt, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	queryStart := &Query{
		Command:      "find",
		Database:     "test",
		RequestID:    1,
		ConnectionID: "1",
		Duration:     0,
		StartAt:      startAt,
		CompletedAt:  startAt,
	}

	storage.Upsert(queryStart)
	query, err := storage.Extract(queryStart.Key())
	assert.NoError(t, err)
	assert.NotNil(t, query)

	query, err = storage.Extract(queryStart.Key())
	assert.Errorf(t, err, fmt.Sprintf("query by key #%s not exists", queryStart.Key()))
	assert.Nil(t, query)
}

func BenchmarkInmemoryQueryMetricStorage_Upsert(b *testing.B) {
	storage := NewQueryEventInmemoryStorage()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			storage.Upsert(&Query{
				Command:      "find",
				Database:     "test",
				RequestID:    int64(rand.Intn(1000000000) + 1),
				ConnectionID: "1",
				Duration:     0,
				StartAt:      time.Now(),
				CompletedAt:  time.Now(),
			})
		}
	})
}

func BenchmarkInmemoryQueryMetricStorage_Usecase(b *testing.B) {
	storage := NewQueryEventInmemoryStorage()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			query := &Query{
				Command:      "find",
				Database:     "test",
				RequestID:    int64(rand.Intn(1000000000) + 1),
				ConnectionID: "1",
				Duration:     0,
				StartAt:      time.Now(),
				CompletedAt:  time.Now(),
			}

			if rand.Int()%2 == 0 {
				storage.Upsert(query)
			} else {
				_, _ = storage.Extract(query.Key())
			}
		}
	})
}
