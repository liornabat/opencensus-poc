package stats

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func createResultMetric(key Key, base *Metric) *Metric {
	m := NewMetric(key)
	m.CacheHit = base.CacheHit
	m.MsgCount = base.MsgCount
	m.MsgSize = base.MsgSize
	m.LastUpdate = base.LastUpdate
	m.Latency = base.Latency
	m.Errors = base.Errors
	m.CacheMiss = base.CacheMiss
	return m
}
func TestKey_SingleKeySingleItem(t *testing.T) {
	tests := []struct {
		name     string
		key      Key
		item     Item
		expected *Metric
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
			expected: &Metric{
				Node:       "node_1",
				Channel:    "",
				Group:      "",
				ClientID:   "client_1",
				Kind:       "",
				MsgCount:   2,
				MsgSize:    0,
				CacheHit:   0,
				CacheMiss:  0,
				Errors:     0,
				Latency:    0,
				LastUpdate: 0,
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
			expected: &Metric{
				Node:       "node_2",
				Channel:    "some_channel_*,|,>%$#*Q1",
				Group:      "",
				ClientID:   "client_2",
				Kind:       "",
				MsgCount:   0,
				MsgSize:    100,
				CacheHit:   0,
				CacheMiss:  0,
				Errors:     0,
				Latency:    0,
				LastUpdate: 0,
			},
		},
		{
			name: "key_elements_3",
			key:  GetKey("node_2", "client_2", "some_channel_*,|,>%$#*Q1", "q1", "", ""),
			item: Item{
				MsgCount:   0,
				MsgSize:    0,
				CacheHit:   2,
				CacheMiss:  3,
				Errors:     4,
				Latency:    0,
				LastUpdate: 0,
			},
			expected: &Metric{
				Node:       "node_2",
				Channel:    "some_channel_*,|,>%$#*Q1",
				Group:      "q1",
				ClientID:   "client_2",
				Kind:       "",
				MsgCount:   0,
				MsgSize:    0,
				CacheHit:   2,
				CacheMiss:  3,
				Errors:     4,
				Latency:    0,
				LastUpdate: 0,
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
			expected: &Metric{
				Node:       "node_2",
				Channel:    "some_channel_*,|,>%$#*Q1",
				Group:      "q1",
				ClientID:   "client_2",
				Kind:       "publish_subscribe",
				MsgCount:   0,
				MsgSize:    0,
				CacheHit:   2,
				CacheMiss:  3,
				Errors:     4,
				Latency:    2,
				LastUpdate: 1000,
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
			resultMap := s.GetMetricsMap()
			metric, ok := resultMap[string(test.key)]
			require.True(t, ok)
			assert.EqualValues(t, test.expected, metric)
			resultMap = s.GetMetricsMap()
			require.Zero(t, len(resultMap))

		})
	}
}
func TestKey_SingleKeyMultiItems(t *testing.T) {
	tests := []struct {
		name     string
		key      Key
		item     []Item
		expected *Metric
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
			expected: &Metric{
				Node:       "node_1",
				Channel:    "some_channel",
				Group:      "",
				ClientID:   "client_multi_1",
				Kind:       "publish",
				MsgCount:   4,
				MsgSize:    200,
				CacheHit:   5,
				CacheMiss:  5,
				Errors:     3,
				Latency:    2.5,
				LastUpdate: 10000,
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
			resultMap := s.GetMetricsMap()
			metric, ok := resultMap[string(test.key)]
			require.True(t, ok)
			assert.EqualValues(t, test.expected, metric)
			resultMap = s.GetMetricsMap()
			require.Zero(t, len(resultMap))

		})
	}
}
func TestKey_SingleSetSingleItem(t *testing.T) {
	tests := []struct {
		name     string
		keys     []Key
		item     []Item
		expected []*Metric
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
			expected: []*Metric{
				&Metric{
					MsgCount:   2,
					MsgSize:    100.0,
					CacheHit:   2,
					CacheMiss:  3,
					Errors:     4,
					Latency:    10000,
					LastUpdate: 10000,
				},
				&Metric{
					MsgCount:   3,
					MsgSize:    200.0,
					CacheHit:   4,
					CacheMiss:  5,
					Errors:     6,
					Latency:    1000,
					LastUpdate: 11000,
				},
				&Metric{
					MsgCount:   3,
					MsgSize:    200.0,
					CacheHit:   4,
					CacheMiss:  5,
					Errors:     6,
					Latency:    1000,
					LastUpdate: 12000,
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
			resultMap := s.GetMetricsMap()
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
			resultMap = s.GetMetricsMap()
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
			resultMap = s.GetMetricsMap()
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

func Benchmark_Set_Insert(b *testing.B) {
	benchmarks := []struct {
		name string
		keys int
	}{
		{
			name: "1_keys",
			keys: 1,
		},
		{
			name: "10_keys",
			keys: 10,
		},
		{
			name: "100_keys",
			keys: 100,
		},
		{
			name: "1000_keys",
			keys: 1000,
		},
		{
			name: "10000_keys",
			keys: 10000,
		},
	}
	Init(WithExportInterval(10000*time.Millisecond), WithInternalExporter())

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			set := NewSet("test")
			for i := 0; i < bm.keys; i++ {
				key := GetKey("some_node", fmt.Sprintf("clinet_id_%d", i), "", "", "", "")
				set.Add(key)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {

				set.Record(
					Item{
						MsgCount:   100,
						MsgSize:    100,
						CacheHit:   1,
						CacheMiss:  1,
						Errors:     1,
						Latency:    1000,
						LastUpdate: 1000,
					})

			}
		})
	}
}
func Benchmark_Key_Insert(b *testing.B) {
	benchmarks := []struct {
		name string
		item Item
	}{
		{
			name: "one_item_one_field",
			item: Item{
				MsgCount:   1,
				MsgSize:    200,
				CacheHit:   1,
				CacheMiss:  0,
				Errors:     0,
				Latency:    0,
				LastUpdate: 0,
			},
		},
	}
	Init(WithExportInterval(10000*time.Millisecond), WithInternalExporter())

	for _, bm := range benchmarks {
		b.Run("check_context"+bm.name, func(b *testing.B) {
			key := GetKey("some_node", "clinet_id", "some_channel", "some_group", "some_kind", "sub_kind")

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key.context(context.Background())
			}
		})
		b.Run("get_key"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				GetKey("some_node", "client_id", "some_channel", "some_group", "some_kind", "sub_kind")

			}
		})
		b.Run("run_Record"+bm.name, func(b *testing.B) {
			key := GetKey("some_node", "clinet_id", "some_channel", "some_group", "some_kind", "sub_kind")
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key.Record(bm.item)
			}
		})
	}
}
