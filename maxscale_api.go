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

// Servers structure reflects JSON object returned by MaxScale REST API
// <maxscale url>/v1/servers
type Servers struct {
	Links interface {
	} `json:"links"`
	Data []struct {
		ID string `json:"id"`
		//nolint
		Type string `json:"type"`
		//nolint
		Relationships interface {
		} `json:"relationships"`
		Attributes struct {
			Parameters struct {
				Address string `json:"address"`
				// add other parameters if needed
			} `json:"parameters"`
			State string `json:"state"`
			// add other parameters if needed
			Statistics struct {
				Connections int `json:"connections"`
				// add other parameters if needed
			} `json:"statistics"`
		} `json:"attributes"`
		//nolint
		Links interface {
		} `json:"links"`
	} `json:"data"`
}

// Services structure reflects JSON object returned by MaxScale REST API
// <maxscale url>/v1/services
type Services struct {
	Links interface {
	} `json:"links"`
	Data []struct {
		ID string `json:"id"`
		// add other parameters if needed
		Attributes struct {
			Router      string `json:"router"`
			Connections int    `json:"connections"`
			//nolint
			RouterDiagnostics interface {
			} `json:"router_diagnostics"`
			//nolint
			Parameters struct {
				MaxConnections int `json:"max_connections"`
			} `json:"parameters"`
			//nolint
			Listeners []interface {
			} `json:"listeners"`
			// add other parameters if needed
		} `json:"attributes"`
		//nolint
		Relationships interface {
		} `json:"relationships"`
		//nolint
		Links interface {
		} `json:"links"`
	} `json:"data"`
}

// Monitors structure reflects JSON object returned by MaxScale REST API
// <maxscale url>/v1/monitors
type Monitors struct {
	Links interface {
	} `json:"links"`
	Data []struct {
		ID string `json:"id"`
		// add other parameters if needed
		Attributes struct {
			Module string `json:"module"`
			//nolint
			MonitorDiagnostics struct {
				Primary bool   `json:"primary"`
				Master  string `json:"master"`
			} `json:"monitor_diagnostics"`
			//nolint
			Parameters struct {
				CooperativeMonitoringLocks string           `json:"cooperative_monitoring_locks"`
				AutoFailover               AutoFailoverValue `json:"auto_failover"`
				AutoRejoin                 bool             `json:"auto_rejoin"`
			} `json:"parameters"`
			//nolint
		} `json:"attributes"`
		//nolint
		Relationships interface {
		} `json:"relationships"`
		//nolint
		Links interface {
		} `json:"links"`
	} `json:"data"`
}

// MaxscaleStatus structure reflects JSON object returned by MaxScale REST API
// <maxscale url>/v1/maxscale
type MaxscaleStatus struct {
	Links interface {
	} `json:"links"`
	Data struct {
		Attributes struct {
			Parameters struct {
				Threads int  `json:"threads"`
				WriteqHighWater int  `json:"writeq_high_water"`
				WriteqLowWater int  `json:"writeq_low_water"`
				Passive bool `json:"passive"`
				// add other parameters if needed
			} `json:"parameters"`
			Uptime int `json:"uptime"`
			// add other parameters if needed
		} `json:"attributes"`
		// add other parameters if needed
	} `json:"data"`
}

// ThreadStatus structure reflects JSON object returned by MaxScale REST API
// <maxscale url>/v1/maxscale/threads
type ThreadStatus struct {
	Links interface {
	} `json:"links"`
	Data []struct {
		ID string `json:"id"`
		// add other parameters if needed
		Attributes struct {
			Stats struct {
				Reads               int `json:"reads"`
				Writes              int `json:"writes"`
				Errors              int `json:"errors"`
				Hangups             int `json:"hangups"`
				Accepts             int `json:"accepts"`
				AvgEventQueueLength int `json:"avg_event_queue_length"`
				MaxEventQueueLength int `json:"max_event_queue_length"`
				MaxExecTime         int `json:"max_exec_time"`
				MaxQueueTime        int `json:"max_queue_time"`
				CurrentDescriptors  int `json:"current_descriptors"`
				TotalDescriptors    int `json:"total_descriptors"`
				Load                struct {
					LastSecond int `json:"last_second"`
					LastMinute int `json:"last_minute"`
					LastHour   int `json:"last_hour"`
				} `json:"load"`
				QueryClassifierCache struct {
					Size      int `json:"size"`
					Inserts   int `json:"inserts"`
					Hits      int `json:"hits"`
					Misses    int `json:"misses"`
					Evictions int `json:"evictions"`
				} `json:"query_classifier_cache"`
			} `json:"stats"`
		} `json:"attributes"`
		//nolint
		Links struct {
		} `json:"links"`
	} `json:"data"`
}
