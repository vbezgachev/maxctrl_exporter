// Copyright 2019, Vitaly Bezgachev, vitaly.bezgachev [the_at_symbol] gmail.com, Kadir Tugan, kadir.tugan [the_at_symbol] gmail.com
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import "github.com/prometheus/client_golang/prometheus"

// Namespace for Prometheus
const (
	Namespace = "maxctrl"
)

// Metric for Prometheus consists of value desciption and type
type Metric struct {
	Desc      *prometheus.Desc
	ValueType prometheus.ValueType
}

var (
	serverLabelNames         = []string{"server", "address"}
	serverUpLabelNames       = []string{"server", "address", "status"}
	serviceLabelNames        = []string{"name", "router"}
	monitorLabelNames        = []string{"name", "cooperative_monitoring_locks"}
	maxscaleStatusLabelNames = []string{}
	statusLabelNames         = []string{"id"}
)

type metrics map[string]Metric

func newDesc(subsystem string, name string, help string, variableLabels []string, t prometheus.ValueType) Metric {
	return Metric{
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, name),
			help, variableLabels, nil),
		ValueType: t,
	}
}

// Exported MaxScale metrics for Prometheus
var (
	ServerMetrics = metrics{
		"server_connections": newDesc("server", "connections", "Amount of connections to the server", serverLabelNames, prometheus.GaugeValue),
		"server_up":          newDesc("server", "up", "Is the server up", serverUpLabelNames, prometheus.GaugeValue),
	}
	ServiceMetrics = metrics{
		"service_current_sessions": newDesc("service", "current_sessions", "Amount of sessions currently active", serviceLabelNames, prometheus.GaugeValue),
		"service_sessions_total":   newDesc("service", "total_sessions", "Total amount of sessions", serviceLabelNames, prometheus.CounterValue),
		"service_max_connections":  newDesc("service", "max_connections", "Max connections allowed", serviceLabelNames, prometheus.GaugeValue),
	}

	MaxscaleStatusMetrics = metrics{
		"status_uptime":  newDesc("status", "uptime", "How long has the server been running", maxscaleStatusLabelNames, prometheus.GaugeValue),
		"status_threads": newDesc("status", "threads", "Number of worker threads", maxscaleStatusLabelNames, prometheus.GaugeValue),
		"status_writeq_high_water": newDesc("status", "writeq_high_water", "High water mark for network write buffer", maxscaleStatusLabelNames, prometheus.GaugeValue),
		"status_writeq_low_water": newDesc("status", "writeq_low_water", "Low water mark for network write buffer", maxscaleStatusLabelNames, prometheus.GaugeValue),
		"status_passive": newDesc("status", "passive", "Has passive mode", maxscaleStatusLabelNames, prometheus.GaugeValue),
	}

	StatusMetrics = metrics{
		"status_read_events":                      newDesc("status", "read_events", "How many read events happened", statusLabelNames, prometheus.CounterValue),
		"status_write_events":                     newDesc("status", "write_events", "How many write events happened", statusLabelNames, prometheus.CounterValue),
		"status_error_events":                     newDesc("status", "error_events", "How many error events happened", statusLabelNames, prometheus.CounterValue),
		"status_hangup_events":                    newDesc("status", "hangup_events", "How many hangup events happened", statusLabelNames, prometheus.CounterValue),
		"status_accept_events":                    newDesc("status", "accept_events", "How many accept events happened", statusLabelNames, prometheus.CounterValue),
		"status_avg_event_queue_length":           newDesc("status", "avg_event_queue_length", "The average length of the event queue", statusLabelNames, prometheus.GaugeValue),
		"status_max_event_queue_length":           newDesc("status", "max_event_queue_length", "The maximum length of the event queue", statusLabelNames, prometheus.GaugeValue),
		"status_max_event_exec_time":              newDesc("status", "max_event_exec_time", "The maximum event execution time", statusLabelNames, prometheus.GaugeValue),
		"status_max_event_queue_time":             newDesc("status", "max_event_queue_time", "The maximum event queue time", statusLabelNames, prometheus.GaugeValue),
		"status_current_descriptors":              newDesc("status", "current_descriptors", "How many current descriptors there are", statusLabelNames, prometheus.GaugeValue),
		"status_total_descriptors":                newDesc("status", "total_descriptors", "How many total descriptors there are", statusLabelNames, prometheus.CounterValue),
		"status_load_last_second":                 newDesc("status", "load_last_second", "The load during the last measured second", statusLabelNames, prometheus.GaugeValue),
		"status_load_last_minute":                 newDesc("status", "load_last_minute", "The load during the last measured minute", statusLabelNames, prometheus.GaugeValue),
		"status_load_last_hour":                   newDesc("status", "load_last_hour", "The load during the last measured hour", statusLabelNames, prometheus.GaugeValue),
		"status_query_classifier_cache_size":      newDesc("status", "query_classifier_cache_size", "The query classifier cache size", statusLabelNames, prometheus.GaugeValue),
		"status_query_classifier_cache_inserts":   newDesc("status", "query_classifier_cache_inserts", "The number of inserts into the query classifier cache", statusLabelNames, prometheus.GaugeValue),
		"status_query_classifier_cache_hits":      newDesc("status", "query_classifier_cache_hits", "The number of hits in the query classifier cache", statusLabelNames, prometheus.GaugeValue),
		"status_query_classifier_cache_misses":    newDesc("status", "query_classifier_cache_misses", "The number of misses in the query classifier cache", statusLabelNames, prometheus.GaugeValue),
		"status_query_classifier_cache_evictions": newDesc("status", "query_classifier_cache_evictions", "The number of evictions in the query classifier cache", statusLabelNames, prometheus.GaugeValue),
	}

	MonitorMetrics = metrics{
		"monitor_primary":       newDesc("monitor", "primary", "Is a primary node", monitorLabelNames, prometheus.GaugeValue),
		"monitor_auto_failover": newDesc("monitor", "auto_failover", "Is auto-failover enabled (1 for true/safe, 0 for false)", monitorLabelNames, prometheus.CounterValue),
		"monitor_auto_rejoin":   newDesc("monitor", "auto_rejoin", "Is auto-rejoin enable", monitorLabelNames, prometheus.GaugeValue),
	}
)
