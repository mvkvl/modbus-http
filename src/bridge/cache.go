package bridge

import (
	"fmt"
	"mbridge/model"
	"slices"
	"time"
)

type MetricCache interface {
	Key(channel *model.Channel, register *model.Register) string
	Get(reference string) *model.Metric
	Set(reference string, value *model.Metric)
	List() []*model.Metric
	Flush()
}

func CreateMetricCache(ttl time.Duration) MetricCache {
	return &metricCacheImpl{
		metrics: make(map[string]*model.Metric, 0),
		ttl:     ttl,
	}
}

type metricCacheImpl struct {
	ttl     time.Duration
	metrics map[string]*model.Metric
}

func (mc *metricCacheImpl) Key(channel *model.Channel, register *model.Register) string {
	return fmt.Sprintf("%s:%s:%s", channel.Title, register.Device.Title, register.Title)
}
func (mc *metricCacheImpl) Get(reference string) *model.Metric {
	v, ok := mc.metrics[reference]
	if !ok {
		return nil
	}
	if v.IsExpired(mc.ttl) {
		return nil
	}
	return v
}
func (mc *metricCacheImpl) Set(reference string, value *model.Metric) {
	mc.metrics[reference] = value
}
func (mc *metricCacheImpl) Flush() {
	for k := range mc.metrics {
		delete(mc.metrics, k)
	}
}
func (mc *metricCacheImpl) List() []*model.Metric {
	var keys []string = make([]string, 0)
	for k, _ := range mc.metrics {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	var result []*model.Metric
	for _, k := range keys {
		m := mc.metrics[k]
		if m.IsExpired(mc.ttl) {
			continue
		}
		result = append(result, m)
	}
	return result
}
