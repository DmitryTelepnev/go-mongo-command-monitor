package mongodb_command_monitor

import (
	"context"
	"sync"
	"testing"

	"github.com/DmitryTelepnev/mongo-command-monitor/metrics"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
)

func BenchmarkGetCommandMonitor(b *testing.B) {
	metric := metrics.NewNilMetrics()
	commandMonitor := GetCommandMonitor(metric)

	database := "test"
	collection := "test"
	commandName := "find"

	bsonRaw := generateBson(commandName, collection)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(50)
	var requestId int64 = 0
	var m sync.Mutex
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Lock()
			requestId++
			m.Unlock()

			commandMonitor.Started(ctx, &event.CommandStartedEvent{
				Command:      bsonRaw,
				DatabaseName: database,
				CommandName:  commandName,
				RequestID:    requestId,
				ConnectionID: "1",
			})

			commandMonitor.Succeeded(ctx, &event.CommandSucceededEvent{
				CommandFinishedEvent: event.CommandFinishedEvent{
					DurationNanos: 0,
					CommandName:   commandName,
					RequestID:     requestId,
					ConnectionID:  "1",
				},
				Reply: nil,
			})
		}
	})
}

func generateBson(operation string, collection string) bson.Raw {
	raw, _ := bson.Marshal(map[string]interface{}{
		"asd":     "asd",
		"asd2":    "asd",
		"asd3":    "asd",
		"asd4":    "asd",
		"asd5":    "asd",
		"asd6":    "asd",
		"asd7":    "asd",
		"asd8":    "asd",
		"asd9":    "asd",
		operation: collection,
	})
	return raw
}
