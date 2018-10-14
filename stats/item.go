package stats

import "time"

type Item struct {
	MsgCount   float64
	MsgSize    float64
	CacheHit   int64
	CacheMiss  int64
	Errors     int64
	Latency    time.Duration
	LastUpdate int64
}
