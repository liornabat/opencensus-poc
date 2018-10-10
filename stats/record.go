package stats

import (
	"context"

	ocstats "go.opencensus.io/stats"
)

func Record(item *Item, keys ...Key) error {

	var ms []ocstats.Measurement
	if item.MsgCount > 0 {
		ms = append(ms, typeIntMeasures[typeMsgCount].M(item.MsgCount))
	}
	if item.MsgSize > 0 {
		ms = append(ms, typeFloatMeasures[typeMsgSize].M(item.MsgSize))
	}
	if item.Errors > 0 {
		ms = append(ms, typeIntMeasures[typeErrors].M(item.Errors))
	}
	if item.CacheHit > 0 {
		ms = append(ms, typeIntMeasures[typeCacheHits].M(item.CacheHit))
	}
	if item.CacheMiss > 0 {
		ms = append(ms, typeIntMeasures[typeCacheMiss].M(item.CacheMiss))
	}
	if item.Latency > 0 {
		ms = append(ms, typeFloatMeasures[typeLatency].M(float64(item.Latency)/1e6))
	}
	if item.LastUpdate > 0 {
		ms = append(ms, typeIntMeasures[typeLastUpdate].M(item.LastUpdate))
	}
	for i := 0; i < len(keys); i++ {
		ctx, err := keys[i].context(context.Background())
		if err != nil {
			return err
		}
		ocstats.Record(ctx, ms...)
	}

	return nil
}
