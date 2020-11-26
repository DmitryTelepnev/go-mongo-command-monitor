package metrics

import "github.com/DmitryTelepnev/mongo-command-monitor/event"

// Metric collector interface
type Metrics interface {
	CollectStartedQuery(query *event.Query)
	CollectFinishedQuery(query *event.Query, status string)
}
