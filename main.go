package main

import (
	"log"
	"net/http"
	"time"

	"github.com/liornabat/opencensus-poc/stats"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
)

func main() {
	go http.ListenAndServe("localhost:8080", nil)

	// Create that Stackdriver stats exporter
	exporter, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		log.Fatalf("Failed to create the Stackdriver stats exporter: %v", err)
	}
	http.Handle("/metrics", exporter)

	// Register the stats exporter
	view.RegisterExporter(exporter)

	view.RegisterExporter(&stats.PrintExporter{})
	// Register the views
	view.SetReportingPeriod(10 * time.Second)
	err = stats.Init("some_node_name")
	if err != nil {
		log.Fatalf("Failed to init stats: %v", err)
	}
	key := stats.MakeKey("client_id", "channel_1", "group_1", "subscribe", "messages")
	key2 := stats.MakeKey("client_id", "channel_2", "group_1", "subscribe", "messages")
	keySet := stats.NewSet("some_set").Add(key, key2)

	for {
		select {
		case <-time.After(1 * time.Second):
			err := keySet.Record(&stats.Item{
				MsgCount:   1,
				MsgSize:    0,
				CacheHit:   0,
				CacheMiss:  0,
				Errors:     0,
				Latency:    0,
				LastUpdate: 0,
			})

			if err != nil {
				log.Fatalf("Failed to record to view: %v", err)
			}
		}
	}
}
