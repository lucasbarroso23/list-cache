package main

import (
	"fmt"
	"sync"
)

type MetricsCache struct {
	sync.RWMutex
	metrics []MetricDataPoint
}

type MetricDataPoint struct {
	MetricName  string
	DeviceUUID  string
	MetricValue int
}

func main() {

	Cache := MetricsCache{
		metrics: []MetricDataPoint{
			{"aurora_node_expired", "test_uuid1", 2},
			{"aurora_node_registration_expired", "test_uuid1", 2},
			{"aurora_node_expired", "test_uuid3", 2},
			{"aurora_node_registration_expired", "test_uuid3", 2},
		},
	}

	dataPointsDB := MetricsCache{
		metrics: []MetricDataPoint{
			{"aurora_node_expired", "test_uuid1", 2},
			{"aurora_node_expired", "test_uuid2", 2},
			{"aurora_node_registration_expired", "test_uuid1", 2},
			{"aurora_node_registration_expired", "test_uuid2", 2},
		},
	}

	deletedMetrics := MetricsCache{}

	fmt.Printf("metrics cache: %v\n", Cache.metrics)

	for value := range Cache.Iter() {

		fmt.Print(value)
		if metric := dataPointsDB.CheckForMetric(value); metric {
			fmt.Println(" metric exists")
		} else {
			fmt.Println(" metric doesn`t exists")
			deletedMetrics.Append(value)
		}
	}

	Cache.metrics = dataPointsDB.metrics

	fmt.Printf("\nUpdated metrics cache\n: %v", Cache.metrics)
	fmt.Printf("\nDeleted metrics cache\n: %v", deletedMetrics.metrics)

}

func (cache *MetricsCache) Append(metric MetricDataPoint) {
	cache.Lock()
	defer cache.Unlock()

	cache.metrics = append(cache.metrics, metric)
}

func (cache *MetricsCache) Iter() <-chan MetricDataPoint {
	c := make(chan MetricDataPoint)

	f := func() {
		cache.Lock()
		defer cache.Unlock()
		for _, metric := range cache.metrics {
			c <- MetricDataPoint{
				MetricName:  metric.MetricName,
				DeviceUUID:  metric.DeviceUUID,
				MetricValue: metric.MetricValue,
			}
		}
		close(c)
	}
	go f()

	return c
}

func (cache *MetricsCache) CheckForMetric(metric MetricDataPoint) bool {
	for _, value := range cache.metrics {
		if value.DeviceUUID == metric.DeviceUUID && value.MetricName == metric.MetricName {
			return true
		}
	}
	return false
}
