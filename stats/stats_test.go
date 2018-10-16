package stats

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func createResultMetric(key Key, base *ChannelSummary) *ChannelSummary {
	m := NewChannelSummary(key)
	m.TotalCacheHits = base.TotalCacheHits
	m.TotalMsgCount = base.TotalMsgCount
	m.TotalMsgSize = base.TotalMsgSize
	m.LastUpdatedUnix = base.LastUpdatedUnix
	m.AvgLatency = base.AvgLatency
	m.TotalErrors = base.TotalErrors
	m.TotalCacheMiss = base.TotalCacheMiss
	return m
}
func TestKey_SingleKeySingleItem(t *testing.T) {
	tests := []struct {
		name     string
		key      Key
		item     Item
		expected *ChannelSummary
	}{
		{
			name: "key_elements_1",
			key:  GetKey("node_1", "client_1", "", "", "", ""),
			item: Item{
				MsgCount:   2,
				MsgSize:    0,
				CacheHit:   0,
				CacheMiss:  0,
				Errors:     0,
				Latency:    0,
				LastUpdate: 0,
			},
			expected: &ChannelSummary{
				Node:            "node_1",
				Channel:         "",
				Group:           "",
				ClientID:        "client_1",
				Kind:            "",
				TotalMsgCount:   2,
				TotalMsgSize:    0,
				AvgMsgSize:      0,
				TotalCacheHits:  0,
				TotalCacheMiss:  0,
				CacheHitsRatio:  0,
				TotalErrors:     0,
				AvgLatency:      0,
				SuccessRate:     100,
				ErrorRate:       0,
				LastUpdatedUnix: 0,
				LastUpdateTime:  time.Time{},
			},
		},
		{
			name: "key_elements_2",
			key:  GetKey("node_2", "client_2", "some_channel_*,|,>%$#*Q1", "", "", ""),
			item: Item{
				MsgCount:   0,
				MsgSize:    100.0,
				CacheHit:   0,
				CacheMiss:  0,
				Errors:     0,
				Latency:    0,
				LastUpdate: 0,
			},
			expected: &ChannelSummary{
				Node:            "node_2",
				Channel:         "some_channel_*,|,>%$#*Q1",
				Group:           "",
				ClientID:        "client_2",
				Kind:            "",
				TotalMsgCount:   0,
				TotalMsgSize:    100,
				AvgMsgSize:      0,
				TotalCacheHits:  0,
				TotalCacheMiss:  0,
				CacheHitsRatio:  0,
				TotalErrors:     0,
				AvgLatency:      0,
				SuccessRate:     100,
				ErrorRate:       0,
				LastUpdatedUnix: 0,
				LastUpdateTime:  time.Time{},
			},
		},
		{
			name: "key_elements_3",
			key:  GetKey("node_2", "client_2", "some_channel_*,|,>%$#*Q1", "q1", "", ""),
			item: Item{
				MsgCount:   4,
				MsgSize:    0,
				CacheHit:   2,
				CacheMiss:  3,
				Errors:     4,
				Latency:    0,
				LastUpdate: 0,
			},
			expected: &ChannelSummary{
				Node:            "node_2",
				Channel:         "some_channel_*,|,>%$#*Q1",
				Group:           "q1",
				ClientID:        "client_2",
				Kind:            "",
				TotalMsgCount:   4,
				TotalMsgSize:    0,
				AvgMsgSize:      0,
				TotalCacheHits:  2,
				TotalCacheMiss:  3,
				CacheHitsRatio:  0.4,
				TotalErrors:     4,
				AvgLatency:      0,
				SuccessRate:     0,
				ErrorRate:       100,
				LastUpdatedUnix: 0,
				LastUpdateTime:  time.Time{},
			},
		},
		{
			name: "key_elements_4",
			key:  GetKey("node_2", "client_2", "some_channel_*,|,>%$#*Q1", "q1", "publish", "subscribe"),
			item: Item{
				MsgCount:   0,
				MsgSize:    0,
				CacheHit:   2,
				CacheMiss:  3,
				Errors:     4,
				Latency:    2 * time.Millisecond,
				LastUpdate: 1000,
			},
			expected: &ChannelSummary{
				Node:            "node_2",
				Channel:         "some_channel_*,|,>%$#*Q1",
				Group:           "q1",
				ClientID:        "client_2",
				Kind:            "publish_subscribe",
				TotalMsgCount:   8,
				TotalMsgSize:    0,
				AvgMsgSize:      0,
				TotalCacheHits:  2,
				TotalCacheMiss:  3,
				CacheHitsRatio:  0.4,
				TotalErrors:     4,
				AvgLatency:      2,
				SuccessRate:     50,
				ErrorRate:       50,
				LastUpdatedUnix: 1000,
				LastUpdateTime:  time.Time{},
			},
		},
	}
	s, err := Init(WithExportInterval(10*time.Millisecond), WithInternalExporter())
	require.NoError(t, err)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.key.Record(test.item)
			require.NoError(t, err)
			time.Sleep(100 * time.Millisecond)
			resultMap, _ := s.GetMetricsMap()
			metric, ok := resultMap[string(test.key)]
			require.True(t, ok)
			assert.EqualValues(t, test.expected, metric)
			resultMap, _ = s.GetMetricsMap()
			require.Zero(t, len(resultMap))

		})
	}
}
func TestKey_SingleKeyMultiItems(t *testing.T) {
	tests := []struct {
		name              string
		key               Key
		item              []Item
		expChannelSummary *ChannelSummary
		expSummary        Summary
	}{
		{
			name: "one_key_multi_items",
			key:  GetKey("node_1", "client_multi_1", "some_channel", "", "publish", ""),
			item: []Item{
				Item{
					MsgCount:   1,
					MsgSize:    50,
					CacheHit:   1,
					CacheMiss:  3,
					Errors:     2,
					Latency:    2 * time.Millisecond,
					LastUpdate: 8000,
				},
				Item{
					MsgCount:   3,
					MsgSize:    150,
					CacheHit:   4,
					CacheMiss:  2,
					Errors:     1,
					Latency:    3 * time.Millisecond,
					LastUpdate: 10000,
				},
			},
			expChannelSummary: &ChannelSummary{
				Node:            "node_1",
				Channel:         "some_channel",
				Group:           "",
				ClientID:        "client_multi_1",
				Kind:            "publish",
				TotalMsgCount:   4,
				TotalMsgSize:    200,
				TotalCacheHits:  5,
				TotalCacheMiss:  5,
				TotalErrors:     3,
				AvgLatency:      2.5,
				LastUpdatedUnix: 10000,
			},
			expSummary: Summary{
				Node:                "node_1",
				TotalMsgCount:       4,
				TotalMsgSize:        200,
				AvgMsgSize:          50,
				TotalCacheHits:      5,
				TotalCacheMiss:      5,
				CacheHitsRatio:      0.5,
				TotalErrors:         3,
				TotalActiveChannels: 1,
				TotalActiveClients:  1,
				SuccessRate:         3 / 4,
				ErrorRate:           100 - (3 / 4),
			},
		},
	}
	s, err := Init(WithExportInterval(10*time.Millisecond), WithInternalExporter())
	require.NoError(t, err)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.key.Record(test.item...)
			require.NoError(t, err)
			time.Sleep(100 * time.Millisecond)
			resultMap, sum := s.GetMetricsMap()
			metric, ok := resultMap[string(test.key)]
			require.True(t, ok)
			assert.EqualValues(t, test.expChannelSummary, metric)
			assert.EqualValues(t, test.expSummary, sum)
			resultMap, _ = s.GetMetricsMap()
			require.Zero(t, len(resultMap))

		})
	}
}

func TestKey_SingleSetSingleItem(t *testing.T) {
	tests := []struct {
		name     string
		keys     []Key
		item     []Item
		expected []*ChannelSummary
	}{
		{
			name: "set_1",
			keys: []Key{
				GetKey("node_1", "client_set_1", "", "", "", ""),
				GetKey("node_2", "client_set_2", "some_channel_*,|,>%$#*Q1", "", "", ""),
				GetKey("node_2", "client_set_2", "some_channel_*,|,>%$#*Q1", "q1", "", ""),
				GetKey("node_2", "client_set_2", "some_channel_*,|,>%$#*Q1", "q1", "publish", "subscribe"),
			},
			item: []Item{
				{
					MsgCount:   2,
					MsgSize:    100.0,
					CacheHit:   2,
					CacheMiss:  3,
					Errors:     4,
					Latency:    10 * time.Second,
					LastUpdate: 10000,
				},
				{
					MsgCount:   3,
					MsgSize:    200.0,
					CacheHit:   4,
					CacheMiss:  5,
					Errors:     6,
					Latency:    1 * time.Second,
					LastUpdate: 11000,
				},
				{
					MsgCount:   3,
					MsgSize:    200.0,
					CacheHit:   4,
					CacheMiss:  5,
					Errors:     6,
					Latency:    1 * time.Second,
					LastUpdate: 12000,
				},
			},
			expected: []*ChannelSummary{
				&ChannelSummary{
					TotalMsgCount:   2,
					TotalMsgSize:    100.0,
					TotalCacheHits:  2,
					TotalCacheMiss:  3,
					TotalErrors:     4,
					AvgLatency:      10000,
					LastUpdatedUnix: 10000,
				},
				&ChannelSummary{
					TotalMsgCount:   3,
					TotalMsgSize:    200.0,
					TotalCacheHits:  4,
					TotalCacheMiss:  5,
					TotalErrors:     6,
					AvgLatency:      1000,
					LastUpdatedUnix: 11000,
				},
				&ChannelSummary{
					TotalMsgCount:   3,
					TotalMsgSize:    200.0,
					TotalCacheHits:  4,
					TotalCacheMiss:  5,
					TotalErrors:     6,
					AvgLatency:      1000,
					LastUpdatedUnix: 12000,
				},
			},
		},
	}
	s, err := Init(WithExportInterval(10*time.Millisecond), WithInternalExporter())
	require.NoError(t, err)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			set := NewSet("set_1")
			set.Add(test.keys...)
			err := set.Record(test.item[0])
			require.NoError(t, err)
			time.Sleep(100 * time.Millisecond)
			resultMap, _ := s.GetMetricsMap()
			//			require.Equal(t, len(test.keys), len(resultMap))
			for _, key := range test.keys {
				expectedMetric := createResultMetric(key, test.expected[0])
				resultMetric, ok := resultMap[string(key)]
				require.True(t, ok)
				assert.EqualValues(t, expectedMetric, resultMetric)
			}
			err = set.Record(test.item[1])
			require.NoError(t, err)
			time.Sleep(100 * time.Millisecond)
			resultMap, _ = s.GetMetricsMap()
			require.Equal(t, len(test.keys), len(resultMap))
			for _, key := range test.keys {
				expectedMetric := createResultMetric(key, test.expected[1])
				resultMetric, ok := resultMap[string(key)]
				require.True(t, ok)
				assert.EqualValues(t, expectedMetric, resultMetric)
			}
			fmt.Println("remove")
			set.Remove(test.keys[0])
			err = set.Record(test.item[2])
			require.NoError(t, err)
			time.Sleep(100 * time.Millisecond)
			resultMap, _ = s.GetMetricsMap()
			require.Equal(t, len(test.keys)-1, len(resultMap))
			for i := 1; i < len(test.keys); i++ {
				expectedMetric := createResultMetric(test.keys[i], test.expected[2])
				resultMetric, ok := resultMap[string(test.keys[i])]
				require.True(t, ok)
				assert.EqualValues(t, expectedMetric, resultMetric)
			}
		})
	}
}
