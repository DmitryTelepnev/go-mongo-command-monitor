package event

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Query struct {
	Command      string
	Database     string
	Collection   string
	RequestID    int64
	ConnectionID string
	Duration     time.Duration
	StartAt      time.Time
	CompletedAt  time.Time
}

func (q *Query) Key() string {
	var str strings.Builder
	str.WriteString(q.Command + "-" + q.ConnectionID + "-" + strconv.FormatInt(q.RequestID, 10))
	return str.String()
}

type QueryEventInmemoryStorage struct {
	sync.RWMutex

	queries map[string]*Query
}

const queryDefaultSize = 1000

func NewQueryEventInmemoryStorage() *QueryEventInmemoryStorage {
	return &QueryEventInmemoryStorage{
		queries: make(map[string]*Query, queryDefaultSize),
	}
}

func (s *QueryEventInmemoryStorage) Upsert(query *Query) {
	s.Lock()
	defer s.Unlock()

	key := query.Key()
	if _, exists := s.queries[key]; exists {
		s.queries[key].CompletedAt = query.CompletedAt
		s.queries[key].Duration = query.Duration
	} else {
		s.queries[key] = query
	}
}

func (s *QueryEventInmemoryStorage) Extract(key string) (*Query, error) {
	s.Lock()
	defer s.Unlock()

	query, exists := s.queries[key]
	if exists {
		delete(s.queries, key)
		return query, nil
	}
	return query, fmt.Errorf("query with key #%s not exists", key)
}
