package stats

import (
	"context"
	"fmt"
	"testing"
	"time"
)

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

func Benchmark_Helpers(b *testing.B) {
	benchmarks := []struct {
		name string
		keys int
	}{
		{
			name: "publish",
		},
		{
			name: "subscribe_1000",
			keys: 1000,
		},
		{
			name: "subscribe_10000",
			keys: 10000,
		},
	}
	Init(WithExportInterval(10000*time.Millisecond), WithInternalExporter())

	for _, bm := range benchmarks {
		b.Run("publish"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ReportPublish("some_node", "some_client_id", "some_channe", 1, 100)
			}
		})
		b.Run("subscribe"+bm.name, func(b *testing.B) {
			m := make(map[Key]struct{})
			for i := 0; i < bm.keys; i++ {
				m[GetKey("some_node", fmt.Sprintf("client_id_%d", i), "some_channel", "", "subscribe", "messages")] = struct{}{}
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ReportMessageSubscribe(m, 1, 100)
			}
		})

	}
}
