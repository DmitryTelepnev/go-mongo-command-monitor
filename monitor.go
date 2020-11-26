package monitor

import (
	"context"
	"time"

	"github.com/DmitryTelepnev/mongo-command-monitor/event"
	"github.com/DmitryTelepnev/mongo-command-monitor/metrics"
	driverEvent "go.mongodb.org/mongo-driver/event"
)

// GetCommandMonitor is mongo command monitor constructor
func GetCommandMonitor(metrics metrics.Metrics) *driverEvent.CommandMonitor {
	queryEventStorage := event.NewQueryEventInmemoryStorage()

	return &driverEvent.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *driverEvent.CommandStartedEvent) {
			query := &event.Query{
				Command:      startedEvent.CommandName,
				Database:     startedEvent.DatabaseName,
				Collection:   determineCollection(startedEvent),
				RequestID:    startedEvent.RequestID,
				ConnectionID: startedEvent.ConnectionID,
				Duration:     0,
				StartAt:      time.Now(),
				CompletedAt:  time.Now(),
			}

			queryEventStorage.Upsert(query)
			metrics.CollectStartedQuery(query)
		},
		Succeeded: func(ctx context.Context, succeededEvent *driverEvent.CommandSucceededEvent) {
			duration := time.Duration(succeededEvent.DurationNanos) * time.Nanosecond
			query := &event.Query{
				Command:      succeededEvent.CommandName,
				Database:     "",
				Collection:   "",
				RequestID:    succeededEvent.RequestID,
				ConnectionID: succeededEvent.ConnectionID,
				Duration:     duration,
				StartAt:      time.Now(),
				CompletedAt:  time.Now(),
			}
			queryEventStorage.Upsert(query)
			enrichedQuery, _ := queryEventStorage.Extract(query.Key())
			metrics.CollectFinishedQuery(enrichedQuery, "success")
		},
		Failed: func(ctx context.Context, failedEvent *driverEvent.CommandFailedEvent) {
			duration := time.Duration(failedEvent.DurationNanos) * time.Nanosecond
			query := &event.Query{
				Command:      failedEvent.CommandName,
				Database:     "",
				Collection:   "",
				RequestID:    failedEvent.RequestID,
				ConnectionID: failedEvent.ConnectionID,
				Duration:     duration,
				StartAt:      time.Now(),
				CompletedAt:  time.Now(),
			}
			queryEventStorage.Upsert(query)
			enrichedQuery, _ := queryEventStorage.Extract(query.Key())
			metrics.CollectFinishedQuery(enrichedQuery, "error")
		},
	}
}

var operationsMap = map[string]struct{}{
	"insert":        struct{}{},
	"find":          struct{}{},
	"update":        struct{}{},
	"delete":        struct{}{},
	"findAndModify": struct{}{},
}

func determineCollection(startedEvent *driverEvent.CommandStartedEvent) string {
	elements, _ := startedEvent.Command.Elements()
	for i := range elements {
		_, ok := operationsMap[elements[i].Key()]
		if ok {
			return elements[i].Value().String()
		}
	}
	return ""
}
