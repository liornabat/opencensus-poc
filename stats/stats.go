package stats

import (
	"time"

	"go.opencensus.io/exporter/prometheus"

	ocstats "go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

type Stats struct {
	opts             statsOptions
	internalExporter *exporter
	promExporter     *prometheus.Exporter
}

func Init(opts ...StateOption) (*Stats, error) {
	s := &Stats{}
	so := statsOptions{
		exportInterval:         5 * time.Second,
		enableInternalExporter: false,
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.apply(&so)
		}
	}
	s.opts = so
	var err error
	if s.opts.enablePrometheus {
		s.promExporter, err = prometheus.NewExporter(prometheus.Options{
			Namespace: s.opts.namespace,
			OnError:   s.opts.errFunc,
		})
		if err != nil {
			return nil, err
		}
		view.RegisterExporter(s.promExporter)

	}
	if s.opts.enableInternalExporter {
		s.internalExporter = NewExporter()
		view.RegisterExporter(s.internalExporter)
	}
	view.SetReportingPeriod(s.opts.exportInterval)
	for _, v := range typeViews {
		if err := view.Register(v); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *Stats) GetMetricsMap() map[string]*Metric {
	return s.internalExporter.aggMap.GetMetricsMap()
}
func (s *Stats) GetPrometheusHandler() *prometheus.Exporter {
	return s.promExporter
}

type statType int

const (
	typeMsgCount statType = iota
	typeMsgSize
	typeCacheHits
	typeCacheMiss
	typeErrors
	typeLatency
	typeLastUpdate
)

func (t statType) String() string {
	return typeNames[t]
}

func (t statType) Stat() ocstats.Measure {
	return typeIntMeasures[t]
}

func (t statType) View() *view.View {
	return typeViews[t]
}

var typeNames = map[statType]string{
	typeMsgCount:   "total_messages",
	typeMsgSize:    "total_message_size",
	typeCacheHits:  "total_cache_hits",
	typeCacheMiss:  "total_cache_miss",
	typeErrors:     "total_errors",
	typeLatency:    "total_latency",
	typeLastUpdate: "LastUpdate",
}

var typeIntMeasures = map[statType]*ocstats.Int64Measure{
	typeCacheHits:  ocstats.Int64("total_cache_hits", "count the number of requests with cache hits", "1"),
	typeCacheMiss:  ocstats.Int64("total_cache_miss", "count the number of requests with cache miss", "1"),
	typeErrors:     ocstats.Int64("total_errors", "count the number of errors", "1"),
	typeLastUpdate: ocstats.Int64("LastUpdate", "unix time of current update", "ns"),
}

var typeFloatMeasures = map[statType]*ocstats.Float64Measure{
	typeMsgSize:  ocstats.Float64("total_message_size", "sum the size of messages", "by"),
	typeLatency:  ocstats.Float64("total_latency", "distribution of requests latency", "ms"),
	typeMsgCount: ocstats.Float64("total_messages", "count the number of messages", "1"),
}

var (
	KeyNode, _     = tag.NewKey("node")
	KeyClientID, _ = tag.NewKey("client_id")
	KeyChannel, _  = tag.NewKey("channel")
	KeyGroup, _    = tag.NewKey("group")
	KeyKind, _     = tag.NewKey("kind")
	KeySubKind, _  = tag.NewKey("sub_kind")
)
var (
	Keys = []tag.Key{KeyNode, KeyClientID, KeyChannel, KeyGroup, KeyKind, KeySubKind}
)

var typeViews = map[statType]*view.View{
	typeMsgCount: &view.View{
		TagKeys:     Keys,
		Measure:     typeFloatMeasures[typeMsgCount],
		Aggregation: view.Sum(),
	},
	typeMsgSize: &view.View{
		TagKeys:     Keys,
		Measure:     typeFloatMeasures[typeMsgSize],
		Aggregation: view.Sum(),
	},
	typeCacheHits: &view.View{
		TagKeys:     Keys,
		Measure:     typeIntMeasures[typeCacheHits],
		Aggregation: view.Count(),
	},
	typeCacheMiss: &view.View{
		TagKeys:     Keys,
		Measure:     typeIntMeasures[typeCacheMiss],
		Aggregation: view.Count(),
	},
	typeErrors: &view.View{
		TagKeys:     Keys,
		Measure:     typeIntMeasures[typeErrors],
		Aggregation: view.Count(),
	},
	typeLatency: &view.View{
		TagKeys:     Keys,
		Measure:     typeFloatMeasures[typeLatency],
		Aggregation: view.Distribution(0, 25, 50, 75, 100, 200, 400, 600, 800, 1000, 2000, 4000, 6000),
	},
	typeLastUpdate: &view.View{
		TagKeys:     Keys,
		Measure:     typeIntMeasures[typeLastUpdate],
		Aggregation: view.LastValue(),
	},
}
