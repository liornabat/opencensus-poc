package stats

import (
	"strings"
	"sync"
	"time"
)

type aggregator interface {
	insert(values ...interface{})
	aggregate() (Key, statType, interface{})
	isTouched() bool
}

type aggCount struct {
	st        statType
	key       Key
	touched   bool
	prevValue int64
	lastValue int64
}

func (a *aggCount) isTouched() bool {
	return a.touched
}

func newAggCount(key Key, st statType) *aggCount {
	a := &aggCount{
		st:        st,
		key:       key,
		prevValue: 0,
		lastValue: 0,
	}
	return a
}

func (a *aggCount) insert(values ...interface{}) {
	if len(values) == 1 {
		count, ok := values[0].(int64)
		if ok && count != a.lastValue {
			a.lastValue = count
			a.touched = true
		}
	}
}

func (a *aggCount) aggregate() (Key, statType, interface{}) {
	diff := a.lastValue - a.prevValue
	a.prevValue = a.lastValue
	a.touched = false
	if diff > 0 {
		return a.key, a.st, diff
	}
	return a.key, a.st, int64(0)
}

type aggSum struct {
	st        statType
	key       Key
	touched   bool
	prevValue float64
	lastValue float64
}

func newAggSum(key Key, st statType) *aggSum {
	a := &aggSum{
		st:  st,
		key: key,
	}
	return a
}

func (a *aggSum) isTouched() bool {
	return a.touched
}

func (a *aggSum) insert(values ...interface{}) {
	if len(values) == 1 {
		sum, ok := values[0].(float64)
		if ok && sum != a.lastValue {
			a.lastValue = sum
			a.touched = true
		}
	}
}

func (a *aggSum) aggregate() (Key, statType, interface{}) {
	diff := a.lastValue - a.prevValue
	a.prevValue = a.lastValue
	a.touched = false
	if diff > 0 {
		return a.key, a.st, diff
	}
	return a.key, a.st, float64(0)
}

type aggLastValue struct {
	st        statType
	key       Key
	touched   bool
	lastValue float64
}

func newAggLastValue(key Key, st statType) *aggLastValue {
	return &aggLastValue{
		st:        st,
		key:       key,
		lastValue: 0,
	}
}

func (a *aggLastValue) isTouched() bool {
	return a.touched
}

func (a *aggLastValue) insert(values ...interface{}) {
	if len(values) == 1 {
		lastValue, ok := values[0].(float64)
		if ok && lastValue != a.lastValue {
			a.lastValue = lastValue
			a.touched = true
		}
	}
}

func (a *aggLastValue) aggregate() (Key, statType, interface{}) {
	lastValue := a.lastValue
	//	a.lastValue = 0
	a.touched = false
	return a.key, a.st, lastValue
}

type ageDistribution struct {
	st        statType
	key       Key
	touched   bool
	prevCount int64
	prevSum   float64
	lastCount int64
	lastSum   float64
}

func newAgeDistribution(key Key, st statType) *ageDistribution {
	return &ageDistribution{
		st:  st,
		key: key,
	}
}

func (a *ageDistribution) isTouched() bool {
	return a.touched
}

func (a *ageDistribution) insert(values ...interface{}) {
	if len(values) == 2 {
		lastCount, ok := values[0].(int64)
		if ok && lastCount != a.lastCount {
			a.lastCount = lastCount
			a.touched = true
		}
		lastSum, ok := values[1].(float64)
		if ok && lastSum != a.lastSum {
			a.lastSum = lastSum
			a.touched = true
		}

	}
}

func (a *ageDistribution) aggregate() (Key, statType, interface{}) {
	diffCount := a.lastCount - a.prevCount
	a.prevCount = a.lastCount

	diffSum := a.lastSum - a.prevSum
	a.prevSum = a.lastSum
	a.touched = false
	if diffCount > 0 {
		return a.key, a.st, diffSum / float64(diffCount)
	}

	return a.key, a.st, float64(0)
}

type aggMap struct {
	sync.Mutex
	m map[string]aggregator
}

func newAggMap() *aggMap {
	return &aggMap{
		m: map[string]aggregator{},
	}
}

func (a *aggMap) insert(index string, values ...interface{}) {
	a.Lock()
	defer a.Unlock()
	agg, ok := a.m[index]
	if !ok {
		params := strings.Split(index, "@@")
		key := Key(params[0])
		statName := params[1]
		switch statName {
		case "total_messages":
			agg = newAggSum(key, typeMsgCount)
		case "total_cache_hits":
			agg = newAggCount(key, typeCacheHits)
		case "total_cache_miss":
			agg = newAggCount(key, typeCacheMiss)
		case "total_errors":
			agg = newAggCount(key, typeErrors)
		case "total_message_size":
			agg = newAggSum(key, typeMsgSize)
		case "total_latency":
			agg = newAgeDistribution(key, typeLatency)
		case "LastUpdatedUnix":
			agg = newAggLastValue(key, typeLastUpdate)
		}
		a.m[index] = agg
	}
	agg.insert(values...)
}

func (a *aggMap) GetChannelSummaryMap() (map[string]*ChannelSummary, Summary) {
	a.Lock()
	defer a.Unlock()
	metricsMap := make(map[string]*ChannelSummary)
	summery := Summary{}
	for _, agg := range a.m {
		if agg.isTouched() {
			key, st, value := agg.aggregate()
			metric, ok := metricsMap[key.String()]
			if !ok {
				metric = NewChannelSummary(key)
				metricsMap[key.String()] = metric
			}
			switch st {
			case typeMsgCount:
				metric.TotalMsgCount = value.(float64)
			case typeMsgSize:
				metric.TotalMsgSize = value.(float64)
			case typeCacheHits:
				metric.TotalCacheHits = value.(int64)
			case typeCacheMiss:
				metric.TotalCacheMiss = value.(int64)
			case typeErrors:
				metric.TotalErrors = value.(int64)
			case typeLatency:
				metric.AvgLatency = value.(float64)
			case typeLastUpdate:
				metric.LastUpdatedUnix = int64(value.(float64))
				metric.LastUpdateTime = time.Unix(metric.LastUpdatedUnix, 0)
			}
			if metric.TotalMsgCount > 0 {
				metric.AvgMsgSize = metric.TotalMsgSize / metric.TotalMsgCount
				metric.ErrorRate = float64(metric.TotalErrors) / metric.TotalMsgCount * 100
			}

			metric.SuccessRate = 100 - metric.ErrorRate
			if metric.TotalCacheHits+metric.TotalCacheMiss > 0 {
				metric.CacheHitsRatio = float64(metric.TotalCacheHits) / float64(metric.TotalCacheHits+metric.TotalCacheMiss)
			}
			summery = summery.AddSummary(metric)
		}

	}
	return metricsMap, summery
}
