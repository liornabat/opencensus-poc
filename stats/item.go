package stats

type Item struct {
	MsgCount   int64
	MsgSize    float64
	CacheHit   int64
	CacheMiss  int64
	Errors     int64
	Latency    int64
	LastUpdate int64
}
