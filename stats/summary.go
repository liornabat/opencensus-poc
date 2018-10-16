package stats

import (
	"fmt"
	"time"
)

type Summary struct {
	Node                string  `json:"node"`
	TotalMsgCount       float64 `json:"total_msg_count"`
	TotalMsgSize        float64 `json:"total_msg_size"`
	AvgMsgSize          float64 `json:"avg_msg_size"`
	TotalCacheHits      int64   `json:"total_cache_hits"`
	TotalCacheMiss      int64   `json:"total_cache_miss"`
	CacheHitsRatio      float64 `json:"cache_hits_ratio"`
	TotalErrors         int64   `json:"total_errors"`
	TotalActiveChannels int64   `json:"total_active_channels"`
	TotalActiveClients  int64   `json:"total_active_clients"`
	SuccessRate         float64 `json:"success_rate"`
	ErrorRate           float64 `json:"error_rate"`
}

func (s Summary) AddSummary(cs *ChannelSummary) Summary {
	if s.Node == "" {
		s.Node = cs.Node
	}
	s.TotalMsgCount += cs.TotalMsgCount
	s.TotalMsgSize += cs.TotalMsgSize
	s.TotalErrors += cs.TotalErrors
	if s.TotalMsgCount > 0 {
		s.AvgMsgSize = s.TotalMsgSize / s.TotalMsgCount
		s.ErrorRate = float64(s.TotalErrors) / s.TotalMsgCount * 100

	}
	s.SuccessRate = 100 - s.ErrorRate

	s.TotalCacheHits += cs.TotalCacheHits
	s.TotalCacheMiss += cs.TotalCacheMiss
	if s.TotalCacheHits+s.TotalCacheMiss > 0 {
		s.CacheHitsRatio = float64(s.TotalCacheHits) / float64(s.TotalCacheHits+s.TotalCacheMiss)
	}
	s.TotalActiveChannels++
	s.TotalActiveClients++

	return s

}

type ChannelSummary struct {
	Node            string    `json:"node"`
	Channel         string    `json:"channel"`
	Group           string    `json:"group"`
	ClientID        string    `json:"client_id"`
	Kind            string    `json:"kind"`
	TotalMsgCount   float64   `json:"total_msg_count"`
	TotalMsgSize    float64   `json:"total_msg_size"`
	AvgMsgSize      float64   `json:"avg_msg_size"`
	TotalCacheHits  int64     `json:"total_cache_hits"`
	TotalCacheMiss  int64     `json:"total_cache_miss"`
	CacheHitsRatio  float64   `json:"cache_hits_ratio"`
	TotalErrors     int64     `json:"total_errors"`
	AvgLatency      float64   `json:"avg_latency"`
	SuccessRate     float64   `json:"success_rate"`
	ErrorRate       float64   `json:"error_rate"`
	LastUpdatedUnix int64     `json:"last_updated_unix"`
	LastUpdateTime  time.Time `json:"last_update_time"`
}

func NewChannelSummary(key Key) *ChannelSummary {
	subKind := key.SubKind()
	if subKind != "" {
		subKind = "_" + subKind
	}
	return &ChannelSummary{
		Node:            key.Node(),
		Channel:         key.Channel(),
		Group:           key.Group(),
		ClientID:        key.ClientID(),
		Kind:            fmt.Sprintf("%s%s", key.Kind(), subKind),
		TotalMsgCount:   0,
		TotalMsgSize:    0,
		AvgMsgSize:      0,
		TotalCacheHits:  0,
		TotalCacheMiss:  0,
		CacheHitsRatio:  0,
		TotalErrors:     0,
		AvgLatency:      0,
		SuccessRate:     0,
		ErrorRate:       0,
		LastUpdatedUnix: 0,
		LastUpdateTime:  time.Time{},
	}
}
