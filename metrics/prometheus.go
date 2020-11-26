package metrics

import (
	"math"

	"github.com/DmitryTelepnev/mongo-command-monitor/event"
	prometheusClient "github.com/prometheus/client_golang/prometheus"
)

type prometheus struct {
	appName              string
	StartedQueryCounter  *prometheusClient.CounterVec
	FinishedQueryCounter *prometheusClient.CounterVec
	QueryDurationHisto   *prometheusClient.HistogramVec
}

//NewPrometheus is the prometheus metric collector constructor
func NewPrometheus(appName string) *prometheus {
	startedQueryCounter := prometheusClient.NewCounterVec(prometheusClient.CounterOpts{
		Name: "mongo_started_query_counter",
		Help: "The total number of started queries",
	}, []string{"app", "database", "collection", "command"})

	finishedQueryCounter := prometheusClient.NewCounterVec(prometheusClient.CounterOpts{
		Name: "mongo_finished_query_counter",
		Help: "The total number of finished queries",
	}, []string{"app", "database", "collection", "command", "status"})

	queryDurationHisto := prometheusClient.NewHistogramVec(prometheusClient.HistogramOpts{
		Name:    "mongo_query_duration",
		Help:    "The mongo query duration",
		Buckets: []float64{.1, .25, .5, 1, 2.5, 5, 10, 20, 30, 50, math.Inf(1)},
	}, []string{"app", "database", "collection", "command", "status"})

	prometheusClient.MustRegister(
		startedQueryCounter,
		finishedQueryCounter,
		queryDurationHisto,
	)

	return &prometheus{
		appName:              appName,
		StartedQueryCounter:  startedQueryCounter,
		FinishedQueryCounter: finishedQueryCounter,
		QueryDurationHisto:   queryDurationHisto,
	}
}

func (p *prometheus) CollectStartedQuery(query *event.Query) {
	p.StartedQueryCounter.WithLabelValues(p.appName, query.Database, query.Collection, query.Command).Inc()
}
func (p *prometheus) CollectFinishedQuery(query *event.Query, status string) {
	p.FinishedQueryCounter.WithLabelValues(p.appName, query.Database, query.Collection, query.Command, status).Inc()
	p.QueryDurationHisto.WithLabelValues(p.appName, query.Database, query.Collection, query.Command, status).Observe(query.Duration.Seconds())
}
