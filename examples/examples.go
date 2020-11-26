package main

import (
	"context"
	"fmt"
	"time"

	commandMonitor "github.com/DmitryTelepnev/mongo-command-monitor"
	"github.com/DmitryTelepnev/mongo-command-monitor/metrics"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	driver "go.mongodb.org/mongo-driver/mongo"
	driverOpts "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Options struct {
	AppName           string
	Hosts             []string
	ConnectionTimeout time.Duration
	ReadPreference    string
	MinPoolSize       uint64
	MaxPoolSize       uint64
}

func main() {
	options := Options{
		Hosts:             []string{"examples_mongo-db_1:27017"},
		ReadPreference:    "primary",
		ConnectionTimeout: 5 * time.Second,
		MinPoolSize:       2,
		MaxPoolSize:       4,
	}
	mongo, err := MongoConnect(options)
	panicOnErr(err)

	type A struct {
		A string `bson:"a"`
	}

	a := A{"asd"}

	fmt.Println("query insertOne")
	ctx := context.Background()
	_, _ = mongo.Database("test").Collection("test_col").InsertOne(ctx, a)

	fmt.Println("query insertMany")
	ctx = context.Background()
	_, _ = mongo.Database("test").Collection("test_col").InsertMany(ctx, []interface{}{a})

	fmt.Println("query findOne")
	ctx = context.Background()
	_ = mongo.Database("test").Collection("test_col").FindOne(ctx, bson.M{"a": "asd"})

	fmt.Println("query find")
	ctx = context.Background()
	_, _ = mongo.Database("test").Collection("test_col").Find(ctx, bson.M{"a": "asd"})
}

func getReadPref(readPref string) (*readpref.ReadPref, error) {
	mode, err := readpref.ModeFromString(readPref)
	if err != nil {
		return nil, fmt.Errorf("mongo readpref not exists %s", readPref)
	}
	rp, err := readpref.New(mode)
	if err != nil {
		return nil, err
	}

	return rp, nil
}

func getPoolMonitor() *event.PoolMonitor {
	return &event.PoolMonitor{Event: func(poolEvent *event.PoolEvent) {
		fmt.Printf("%d\n", poolEvent.ConnectionID)
	}}
}

func MongoConnect(options Options) (*driver.Client, error) {
	connectCtx, connectCancel := context.WithTimeout(context.Background(), options.ConnectionTimeout)
	opts := driverOpts.Client()
	opts.SetHosts(options.Hosts)
	opts.SetConnectTimeout(options.ConnectionTimeout)
	opts.SetAppName(options.AppName)
	opts.SetMinPoolSize(options.MinPoolSize)
	opts.SetMinPoolSize(options.MaxPoolSize)
	//opts.SetPoolMonitor(getPoolMonitor())
	opts.SetMonitor(commandMonitor.GetCommandMonitor(metrics.NewPrometheus("examples-app")))

	readPref, err := getReadPref(options.ReadPreference)
	if err != nil {
		return nil, err
	}

	opts.SetReadPreference(readPref)

	client, err := driver.Connect(connectCtx, opts)
	if err != nil {
		connectCancel()
		return nil, err
	}
	connectCancel()

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(pingCtx, readPref)
	if err != nil {
		pingCancel()
		return nil, err
	}
	pingCancel()

	return client, nil
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
