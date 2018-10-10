package stats

import (
	"context"
	"strings"

	"go.opencensus.io/tag"
)

type Key string

func MakeKey(clientID, channel, group, kind, subKind string) Key {
	fields := []string{clientID, channel, group, kind, subKind}
	return Key(strings.Join(fields, "|"))
}
func (k Key) getElement(n int) string {
	fields := strings.Split(string(k), "|")
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
func (k Key) ClientID() string {
	return k.getElement(0)
}
func (k Key) Channel() string {
	return k.getElement(1)
}
func (k Key) Group() string {
	return k.getElement(2)
}
func (k Key) Kind() string {
	return k.getElement(3)
}
func (k Key) SubKind() string {
	return k.getElement(4)
}

func (k Key) context(ctx context.Context) (context.Context, error) {
	mut := []tag.Mutator{tag.Insert(KeyNode, node)}
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
