package stats

import "time"

const (
	KindPublish            = "publish"
	KindPublishPersistence = "publish_persistence"
)

func ReportPublish(node, clientID, channel string, msgCount, msgSize float64) {
	key := GetKey(node, clientID, channel, "", KindPublish, "")
	key.Record(Item{
		MsgCount:   msgCount,
		MsgSize:    msgSize,
		CacheHit:   0,
		CacheMiss:  0,
		Errors:     0,
		Latency:    0,
		LastUpdate: time.Now().UTC().UnixNano(),
	})
}

func ReportPublishError(node, clientID, channel string) {
	key := GetKey(node, clientID, channel, "", KindPublish, "")
	key.Record(Item{
		Errors:     1,
		LastUpdate: time.Now().UTC().UnixNano(),
	})
}
func ReportPublishPersistence(node, clientID, channel string, msgCount, msgSize float64) {
	key := GetKey(node, clientID, channel, "", KindPublishPersistence, "")
	key.Record(Item{
		MsgCount:   msgCount,
		MsgSize:    msgSize,
		CacheHit:   0,
		CacheMiss:  0,
		Errors:     0,
		Latency:    0,
		LastUpdate: time.Now().UTC().UnixNano(),
	})
}

func ReportPublishPersistenceError(node, clientID, channel string) {
	key := GetKey(node, clientID, channel, "", KindPublishPersistence, "")
	key.Record(Item{
		Errors:     1,
		LastUpdate: time.Now().UTC().UnixNano(),
	})
}

func ReportMessageSubscribe(keys map[Key]struct{}, msgCount, msgSize float64) {
	item := Item{
		MsgCount:   msgCount,
		MsgSize:    msgSize,
		LastUpdate: time.Now().UTC().UnixNano(),
	}
	go func() {
		for key, _ := range keys {
			key.Record(item)
		}
	}()

}
