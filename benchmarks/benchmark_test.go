package benchmarks

import (
	"context"
	"fmt"

	commandMonitor "github.com/DmitryTelepnev/mongo-command-monitor"
	"github.com/DmitryTelepnev/mongo-command-monitor/metrics"
	"go.mongodb.org/mongo-driver/bson"
	driver "go.mongodb.org/mongo-driver/mongo"
	driverOpts "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"math/rand"
	"testing"
	"time"
)

func mongoConnect(options *driverOpts.ClientOptions) (*driver.Client, error) {
	connectCtx, connectCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer connectCancel()
	client, err := driver.Connect(connectCtx, options)
	if err != nil {
		return nil, err
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer pingCancel()
	err = client.Ping(pingCtx, options.ReadPreference)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	database   = "test"
	collection = "test"
)

type testDoc struct {
	I int64  `bson:"i"`
	S string `bson:"s"`
}

var mongoConnection *driver.Client

func getMongoConnect() *driver.Client {
	if mongoConnection == nil {
		options := driverOpts.Client()
		options.SetHosts([]string{"benchmarks_mongo-db-with-monitor_1:27017"})
		options.SetReadPreference(readpref.Primary())
		options.SetConnectTimeout(5 * time.Second)
		options.SetMonitor(commandMonitor.GetCommandMonitor(metrics.NewPrometheus("bench-app")))
		client, err := mongoConnect(options)
		panicOnErr(err)

		_ = client.Database(database).CreateCollection(context.Background(), collection)

		_, _ = client.Database(database).Collection(collection).Indexes().CreateOne(context.Background(), driver.IndexModel{
			Keys: bson.M{
				"i": 1,
				"s": 1,
			},
		})

		mongoConnection = client
	}

	return mongoConnection
}

var mongoConnectionWithoutMonitor *driver.Client

func getMongoConnectWithoutMonitor() *driver.Client {
	if mongoConnectionWithoutMonitor == nil {
		options := driverOpts.Client()
		options.SetHosts([]string{"benchmarks_mongo-db-without-monitor_1:27017"})
		options.SetReadPreference(readpref.Primary())
		options.SetConnectTimeout(5 * time.Second)
		client, err := mongoConnect(options)
		panicOnErr(err)

		_ = client.Database(database).CreateCollection(context.Background(), collection)

		_, _ = client.Database(database).Collection(collection).Indexes().CreateOne(context.Background(), driver.IndexModel{
			Keys: bson.M{
				"i": 1,
				"s": 1,
			},
		})

		mongoConnectionWithoutMonitor = client
	}

	return mongoConnectionWithoutMonitor
}

func BenchmarkMongoInsertQueries_WithCommandMonitor(b *testing.B) {
	client := getMongoConnect()

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			randNum := rand.Int63()
			_, _ = client.Database(database).Collection(collection).InsertOne(ctx, testDoc{
				I: randNum,
				S: fmt.Sprintf("test-%d", randNum),
			})
		}
	})
}

func BenchmarkMongoInsertQueries_WithoutCommandMonitor(b *testing.B) {
	client := getMongoConnectWithoutMonitor()

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			randNum := rand.Int63()
			_, _ = client.Database(database).Collection(collection).InsertOne(ctx, testDoc{
				I: randNum,
				S: fmt.Sprintf("test-%d", randNum),
			})
		}
	})
}

func BenchmarkMongoUpdateQueries_WithCommandMonitor(b *testing.B) {
	client := getMongoConnect()

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			updateOpts := driverOpts.Update()
			updateOpts.SetUpsert(true)

			randNum := rand.Int63()
			_, _ = client.Database(database).Collection(collection).UpdateOne(ctx,
				bson.D{{Key: "i", Value: randNum}},
				bson.D{{Key: "$set", Value: bson.D{{Key: "s", Value: fmt.Sprintf("%d-test", randNum)}}}},
				updateOpts,
			)
		}
	})
}

func BenchmarkMongoUpdateQueries_WithoutCommandMonitor(b *testing.B) {
	client := getMongoConnectWithoutMonitor()

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			updateOpts := driverOpts.Update()
			updateOpts.SetUpsert(true)

			randNum := rand.Int63()
			_, _ = client.Database(database).Collection(collection).UpdateOne(ctx,
				bson.D{{Key: "i", Value: randNum}},
				bson.D{{Key: "$set", Value: bson.D{{Key: "s", Value: fmt.Sprintf("%d-test", randNum)}}}},
				updateOpts,
			)
		}
	})
}

func BenchmarkMongoFindQueries_WithCommandMonitor(b *testing.B) {
	client := getMongoConnect()

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			updateOpts := driverOpts.Update()
			updateOpts.SetUpsert(true)

			randNum := rand.Int63()
			_, _ = client.Database(database).Collection(collection).Find(ctx,
				bson.D{{Key: "i", Value: randNum}},
			)
		}
	})
}

func BenchmarkMongoFindQueries_WithoutCommandMonitor(b *testing.B) {
	client := getMongoConnectWithoutMonitor()

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			updateOpts := driverOpts.Update()
			updateOpts.SetUpsert(true)

			randNum := rand.Int63()
			_, _ = client.Database(database).Collection(collection).Find(ctx,
				bson.D{{Key: "i", Value: randNum}},
			)
		}
	})
}
