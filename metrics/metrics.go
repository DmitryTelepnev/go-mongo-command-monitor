package metrics

import "github.com/DmitryTelepnev/mongo-command-monitor/event"

type Metrics interface {
	CollectStartedQuery(query *event.Query)
	CollectFinishedQuery(query *event.Query, status string)
}
