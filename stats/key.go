package stats

import (
	"context"
	"strings"

	ocstats "go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

const separator = "<|>"

type Key string

func GetKey(node, clientID, channel, group, kind, subKind string) Key {
	fields := []string{node, clientID, channel, group, kind, subKind}
	return Key(strings.Join(fields, separator))
}
func makeKeyFromTags(tags []tag.Tag) Key {
	m := make(map[string]string, len(tags))
	for _, t := range tags {
		m[t.Key.Name()] = t.Value
	}
	fields := []string{
		m["node"],
		m["client_id"],
		m["channel"],
		m["group"],
		m["kind"],
		m["sub_kind"],
	}
	return Key(strings.Join(fields, separator))
}

func (k Key) getElement(n int) string {
	fields := strings.Split(string(k), separator)
	if n+1 > len(fields) {
		return ""
	}
	return fields[n]
}
func EqualKeys(a, b Key) bool {
	return string(a) == string(b)
}

func (k Key) String() string {
	return string(k)
}

func (k Key) Node() string {
	return k.getElement(0)
}
func (k Key) ClientID() string {
	return k.getElement(1)
}
func (k Key) Channel() string {
	return k.getElement(2)
}
func (k Key) Group() string {
	return k.getElement(3)
}
func (k Key) Kind() string {
	return k.getElement(4)
}
func (k Key) SubKind() string {
	return k.getElement(5)
}

func (k Key) context(ctx context.Context) (context.Context, error) {
	var mut []tag.Mutator
	node := k.Node()
	if node != "" {
		mut = append(mut, tag.Insert(KeyNode, node))
	}
	clientID := k.ClientID()
	if clientID != "" {
		mut = append(mut, tag.Insert(KeyClientID, clientID))
	}
	channel := k.Channel()
	if channel != "" {
		mut = append(mut, tag.Insert(KeyChannel, channel))
	}

	group := k.Group()
	if group != "" {
		mut = append(mut, tag.Insert(KeyGroup, group))
	}

	kind := k.Kind()
	if kind != "" {
		mut = append(mut, tag.Insert(KeyKind, kind))
	}

	subKind := k.SubKind()
	if subKind != "" {
		mut = append(mut, tag.Insert(KeySubKind, subKind))
	}

	return tag.New(ctx, mut...)
}

func (k Key) Record(items ...*Item) error {
	var ms []ocstats.Measurement
	for i := 0; i < len(items); i++ {
		for j := 0; j < int(items[i].MsgCount); j++ {
			ms = append(ms, typeIntMeasures[typeMsgCount].M(1))
		}

		if items[i].MsgSize > 0 {
			ms = append(ms, typeFloatMeasures[typeMsgSize].M(items[i].MsgSize))
		}
		for j := 0; j < int(items[i].Errors); j++ {
			ms = append(ms, typeIntMeasures[typeErrors].M(1))
		}

		for j := 0; j < int(items[i].CacheHit); j++ {
			ms = append(ms, typeIntMeasures[typeCacheHits].M(1))
		}
		for j := 0; j < int(items[i].CacheMiss); j++ {
			ms = append(ms, typeIntMeasures[typeCacheMiss].M(1))
		}

		if items[i].Latency > 0 {
			ms = append(ms, typeFloatMeasures[typeLatency].M(float64(items[i].Latency)/1e6))
		}
		if items[i].LastUpdate > 0 {
			ms = append(ms, typeIntMeasures[typeLastUpdate].M(items[i].LastUpdate))
		}
	}
	ctx, err := k.context(context.Background())
	if err != nil {
		return err
	}
	ocstats.Record(ctx, ms...)
	return nil
}
