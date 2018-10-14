package stats

import "fmt"

type Metric struct {
	Node       string
	Channel    string
	Group      string
	ClientID   string
	Kind       string
	MsgCount   int64
	MsgSize    float64
	CacheHit   int64
	CacheMiss  int64
	Errors     int64
	Latency    float64
	LastUpdate int64
}

func NewMetric(key Key) *Metric {
	subKind := key.SubKind()
	if subKind != "" {
		subKind = "_" + subKind
	}
	return &Metric{
		Node:       key.Node(),
		Channel:    key.Channel(),
		Group:      key.Group(),
		ClientID:   key.ClientID(),
		Kind:       fmt.Sprintf("%s%s", key.Kind(), subKind),
		MsgCount:   0,
		MsgSize:    0,
		CacheHit:   0,
		CacheMiss:  0,
		Errors:     0,
		Latency:    0,
		LastUpdate: 0,
	}
}
