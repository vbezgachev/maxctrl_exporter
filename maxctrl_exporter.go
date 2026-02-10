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

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
)

const (
	metricsPath = "/metrics"
	localIP     = "0.0.0.0"
)

var (
	maxScaleUrl                   string // URL of maxscale instance
	maxScaleUsername              string // Username for maxscale REST API authentication
	maxScalePassword              string // Password for maxscale REST API authentication
	maxScaleExporterPort          string // Port for this exporter to run on
	maxScaleCACertificate         string // File containing CA certificate
	maxctrlExporterConfigFile     string // File containing exporter config
	maxScaleTLSInsecureSkipVerify bool   // Disable TLS verify
)

type ConfigValues struct {
	Url                   string `yaml:"url"`
	Username              string `yaml:"username"`
	Password              string `yaml:"password"`
	ExporterPort          string `yaml:"exporter_port"`
	CACertificate         string `yaml:"caCertificate"`
	TLSInsecureSkipVerify bool   `yaml:"tlsInsecureSkipVerify" default:"false"`
}

// MaxScale contains connection parameters to the server and metric maps
type MaxScale struct {
	url                   string
	username              string
	password              string
	transport             *http.Transport
	up                    prometheus.Gauge
	totalScrapes          prometheus.Counter
	serverMetrics         map[string]Metric
	serviceMetrics        map[string]Metric
	monitorMetrics        map[string]Metric
	maxscaleStatusMetrics map[string]Metric
	statusMetrics         map[string]Metric
}

// NewExporter creates a new instance of the MaxScale
func NewExporter(url string, username string, password string, caCertificate string, tlsInsecureSkipVerify bool) (*MaxScale, error) {
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	if len(caCertificate) > 0 {
		// Read in the cert file
		certs, err := os.ReadFile(caCertificate)
		if err != nil {
			log.Fatalf("Failed to open CA certificate file %q: %v", caCertificate, err)
		}

		// Append our cert to the system pool
		if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
			log.Fatalf("Could not append certificate to the root store from file %s", caCertificate)
		}
	}

	transport := &http.Transport{TLSClientConfig: &tls.Config{
		RootCAs:            rootCAs,
		InsecureSkipVerify: tlsInsecureSkipVerify,
	}}

	return &MaxScale{
		url:       url,
		username:  username,
		password:  password,
		transport: transport,
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "up",
			Help:      "Was the last scrape of MaxScale successful?",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "exporter_total_scrapes",
			Help:      "Current total MaxScale scrapes",
		}),
		serverMetrics:         ServerMetrics,
		serviceMetrics:        ServiceMetrics,
		maxscaleStatusMetrics: MaxscaleStatusMetrics,
		statusMetrics:         StatusMetrics,
		monitorMetrics:        MonitorMetrics,
	}, nil
}

// Describe describes all the metrics ever exported by the MaxScale exporter. It
// implements prometheus.Collector.
func (m *MaxScale) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range m.serverMetrics {
		ch <- m.Desc
	}

	for _, m := range m.serviceMetrics {
		ch <- m.Desc
	}

	for _, m := range m.maxscaleStatusMetrics {
		ch <- m.Desc
	}

	for _, m := range m.statusMetrics {
		ch <- m.Desc
	}

	ch <- m.up.Desc()
	ch <- m.totalScrapes.Desc()
}

// Collect fetches the stats from configured MaxScale location and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (m *MaxScale) Collect(ch chan<- prometheus.Metric) {
	m.totalScrapes.Inc()

	var parseErrors = false

	if err := m.parseServers(ch); err != nil {
		parseErrors = true
		log.Print(err)
	}

	if err := m.parseServices(ch); err != nil {
		parseErrors = true
		log.Print(err)
	}

	if err := m.parseMaxscaleStatus(ch); err != nil {
		parseErrors = true
		log.Print(err)
	}

	if err := m.parseThreadStatus(ch); err != nil {
		parseErrors = true
		log.Print(err)
	}

	if err := m.parseMonitors(ch); err != nil {
		parseErrors = true
		log.Print(err)
	}

	if parseErrors {
		m.up.Set(0)
	} else {
		m.up.Set(1)
	}

	ch <- m.up
	ch <- m.totalScrapes
}

func (m *MaxScale) getStatistics(path string, v interface{}) error {
	var err error
	req, err := http.NewRequest("GET", m.url+"/v1"+path, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(m.username, m.password)

	client := &http.Client{Transport: m.transport}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error while getting %v: %v", path, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("the MaxScale statistic request failed with a status: %s", resp.Status)
		return err
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

func serverUp(status string) int {
	if strings.Contains(status, ",Down,") {
		return 0
	}
	if strings.Contains(status, ",Running,") {
		return 1
	}
	return 0
}

func (m *MaxScale) createMetricForPrometheus(metricsMap map[string]Metric, metricKey string,
	value int, ch chan<- prometheus.Metric, labelValues ...string) {

	metric := metricsMap[metricKey]
	ch <- prometheus.MustNewConstMetric(
		metric.Desc,
		metric.ValueType,
		float64(value),
		labelValues...,
	)
}

func (m *MaxScale) parseServers(ch chan<- prometheus.Metric) error {
	var servers Servers
	err := m.getStatistics("/servers", &servers)

	if err != nil {
		return err
	}

	for _, server := range servers.Data {
		serverID := server.ID
		serverAddress := server.Attributes.Parameters.Address
		m.createMetricForPrometheus(m.serverMetrics, "server_connections",
			server.Attributes.Statistics.Connections, ch, serverID, serverAddress)

		// We surround the separated list with the separator as well. This way regular expressions
		// in labeling don't have to consider satus positions.
		normalizedStatus := "," + strings.Replace(server.Attributes.State, ", ", ",", -1) + ","
		m.createMetricForPrometheus(m.serverMetrics, "server_up",
			serverUp(normalizedStatus), ch, serverID, serverAddress, normalizedStatus)
	}

	return nil
}

func (m *MaxScale) parseServices(ch chan<- prometheus.Metric) error {
	var services Services
	err := m.getStatistics("/services", &services)

	if err != nil {
		return err
	}

	for _, service := range services.Data {
		m.createMetricForPrometheus(m.serviceMetrics, "service_current_sessions",
			service.Attributes.Connections, ch, service.ID, service.Attributes.Router)

		m.createMetricForPrometheus(m.serviceMetrics, "service_sessions_total",
			service.Attributes.Connections, ch, service.ID, service.Attributes.Router)

		m.createMetricForPrometheus(m.serviceMetrics, "service_max_connections",
			service.Attributes.Parameters.MaxConnections, ch, service.ID, service.Attributes.Router)
	}

	return nil
}

func (m *MaxScale) parseMaxscaleStatus(ch chan<- prometheus.Metric) error {
	var maxscaleStatus MaxscaleStatus
	err := m.getStatistics("/maxscale", &maxscaleStatus)

	if err != nil {
		return err
	}

	m.createMetricForPrometheus(m.maxscaleStatusMetrics, "status_uptime",
		maxscaleStatus.Data.Attributes.Uptime, ch)

	m.createMetricForPrometheus(m.maxscaleStatusMetrics, "status_threads",
		maxscaleStatus.Data.Attributes.Parameters.Threads, ch)

	m.createMetricForPrometheus(m.maxscaleStatusMetrics, "status_writeq_high_water",
		maxscaleStatus.Data.Attributes.Parameters.WriteqHighWater, ch)

	m.createMetricForPrometheus(m.maxscaleStatusMetrics, "status_writeq_low_water",
		maxscaleStatus.Data.Attributes.Parameters.WriteqLowWater, ch)

	passiveMode := 0
	if maxscaleStatus.Data.Attributes.Parameters.Passive {
		passiveMode = 1
	}

	m.createMetricForPrometheus(m.maxscaleStatusMetrics, "status_passive", passiveMode, ch)

	return nil
}

func (m *MaxScale) parseMonitors(ch chan<- prometheus.Metric) error {
	var monitors Monitors
	err := m.getStatistics("/monitors", &monitors)

	if err != nil {
		return err
	}

	for _, monitor := range monitors.Data {

		primary := 0
		if monitor.Attributes.MonitorDiagnostics.Primary {
			primary = 1
		}
		m.createMetricForPrometheus(m.monitorMetrics, "monitor_primary", primary, ch, monitor.ID, monitor.Attributes.Parameters.CooperativeMonitoringLocks)

		auto_failover := 0
		// Use the Bool() method of AutoFailoverValue which handles specific enum values
		// "true" and "safe" are considered true, all others are false
		if monitor.Attributes.Parameters.AutoFailover.Bool() {
			auto_failover = 1
		}
		m.createMetricForPrometheus(m.monitorMetrics, "monitor_auto_failover", auto_failover, ch, monitor.ID, monitor.Attributes.Parameters.CooperativeMonitoringLocks)

		auto_rejoin := 0
		if monitor.Attributes.Parameters.AutoRejoin {
			auto_rejoin = 1
		}
		m.createMetricForPrometheus(m.monitorMetrics, "monitor_auto_rejoin", auto_rejoin, ch, monitor.ID, monitor.Attributes.Parameters.CooperativeMonitoringLocks)

	}

	return nil
}

func (m *MaxScale) parseThreadStatus(ch chan<- prometheus.Metric) error {
	var threadStatus ThreadStatus
	err := m.getStatistics("/maxscale/threads", &threadStatus)

	if err != nil {
		return err
	}

	for _, threadStatus := range threadStatus.Data {
		m.createMetricForPrometheus(m.statusMetrics, "status_read_events",
			threadStatus.Attributes.Stats.Reads, ch, threadStatus.ID)
		m.createMetricForPrometheus(m.statusMetrics, "status_write_events",
			threadStatus.Attributes.Stats.Writes, ch, threadStatus.ID)
		m.createMetricForPrometheus(m.statusMetrics, "status_error_events",
			threadStatus.Attributes.Stats.Errors, ch, threadStatus.ID)
		m.createMetricForPrometheus(m.statusMetrics, "status_hangup_events",
			threadStatus.Attributes.Stats.Hangups, ch, threadStatus.ID)
		m.createMetricForPrometheus(m.statusMetrics, "status_accept_events",
			threadStatus.Attributes.Stats.Accepts, ch, threadStatus.ID)
		m.createMetricForPrometheus(m.statusMetrics, "status_avg_event_queue_length",
			threadStatus.Attributes.Stats.AvgEventQueueLength, ch, threadStatus.ID)
		m.createMetricForPrometheus(m.statusMetrics, "status_max_event_queue_length",
			threadStatus.Attributes.Stats.MaxEventQueueLength, ch, threadStatus.ID)
		m.createMetricForPrometheus(m.statusMetrics, "status_max_event_exec_time",
			threadStatus.Attributes.Stats.MaxExecTime, ch, threadStatus.ID)
		m.createMetricForPrometheus(m.statusMetrics, "status_max_event_queue_time",
			threadStatus.Attributes.Stats.MaxQueueTime, ch, threadStatus.ID)
		m.createMetricForPrometheus(m.statusMetrics, "status_current_descriptors",
			threadStatus.Attributes.Stats.CurrentDescriptors, ch, threadStatus.ID)
		m.createMetricForPrometheus(m.statusMetrics, "status_total_descriptors",
			threadStatus.Attributes.Stats.TotalDescriptors, ch, threadStatus.ID)
		m.createMetricForPrometheus(m.statusMetrics, "status_load_last_second",
			threadStatus.Attributes.Stats.Load.LastSecond, ch, threadStatus.ID)
		m.createMetricForPrometheus(m.statusMetrics, "status_load_last_minute",
			threadStatus.Attributes.Stats.Load.LastMinute, ch, threadStatus.ID)
		m.createMetricForPrometheus(m.statusMetrics, "status_load_last_hour",
			threadStatus.Attributes.Stats.Load.LastHour, ch, threadStatus.ID)
		m.createMetricForPrometheus(m.statusMetrics, "status_query_classifier_cache_size",
			threadStatus.Attributes.Stats.QueryClassifierCache.Size, ch, threadStatus.ID)
		m.createMetricForPrometheus(m.statusMetrics, "status_query_classifier_cache_inserts",
			threadStatus.Attributes.Stats.QueryClassifierCache.Inserts, ch, threadStatus.ID)
		m.createMetricForPrometheus(m.statusMetrics, "status_query_classifier_cache_hits",
			threadStatus.Attributes.Stats.QueryClassifierCache.Hits, ch, threadStatus.ID)
		m.createMetricForPrometheus(m.statusMetrics, "status_query_classifier_cache_misses",
			threadStatus.Attributes.Stats.QueryClassifierCache.Misses, ch, threadStatus.ID)
		m.createMetricForPrometheus(m.statusMetrics, "status_query_classifier_cache_evictions",
			threadStatus.Attributes.Stats.QueryClassifierCache.Evictions, ch, threadStatus.ID)
	}

	return nil
}

// GetEnvVar - retrieves values of environment variables. If nothing is set, return the default value
func GetEnvVar(envName string, defaultValue string) string {
	envVal := os.Getenv(envName)
	if envVal == "" {
		return defaultValue
	}
	return envVal
}

func readConfigFile(fname string) {
	yamlFile, err := os.ReadFile(fname)
	if err != nil {
		// If the file doesn't exist, just return
		if errors.Is(err, fs.ErrNotExist) {
			return
		} else {
			log.Printf("Could not open configuration file '%s': %v", maxctrlExporterConfigFile, err)
		}
	}
	parseConfigFile(yamlFile)
}

func parseConfigFile(contents []byte) {
	var config ConfigValues
	err := yaml.Unmarshal(contents, &config)
	if err != nil {
		log.Fatalf("Could not parse config file contents: %v", err)
	}
	if config.Url != "" {
		maxScaleUrl = config.Url
	}
	if config.Username != "" {
		maxScaleUsername = config.Username
	}
	if config.Password != "" {
		maxScalePassword = config.Password
	}
	if config.ExporterPort != "" {
		maxScaleExporterPort = config.ExporterPort
	}
	if config.CACertificate != "" {
		maxScaleCACertificate = config.CACertificate
	}
	maxScaleTLSInsecureSkipVerify = config.TLSInsecureSkipVerify
}

func setConfigFromEnvironmentVars() {
	var err error
	maxScaleUrl = GetEnvVar("MAXSCALE_URL", "http://127.0.0.1:8989")
	maxScaleUsername = GetEnvVar("MAXSCALE_USERNAME", "admin")
	maxScalePassword = GetEnvVar("MAXSCALE_PASSWORD", "mariadb")
	maxScaleExporterPort = GetEnvVar("MAXSCALE_EXPORTER_PORT", "8080")
	if maxScaleTLSInsecureSkipVerify, err = strconv.ParseBool(GetEnvVar("MAXSCALE_TLS_INSECURE_SKIP_VERIFY", "false")); err != nil {
		maxScaleTLSInsecureSkipVerify = false
	}
	maxScaleCACertificate = GetEnvVar("MAXSCALE_CA_CERTIFICATE", "")
	maxctrlExporterConfigFile = GetEnvVar("MAXCTRL_EXPORTER_CFG_FILE", "maxctrl_exporter.yaml")
}

func main() {
	setConfigFromEnvironmentVars()
	readConfigFile(maxctrlExporterConfigFile)

	log.Print("Starting MaxScale exporter")
	log.Printf("Scraping MaxScale JSON API at: %s", maxScaleUrl)

	exporter, err := NewExporter(maxScaleUrl, maxScaleUsername, maxScalePassword, maxScaleCACertificate, maxScaleTLSInsecureSkipVerify)
	if err != nil {
		log.Fatalf("Failed to start maxscale exporter: %v\n", err)
	}

	prometheus.MustRegister(exporter)
	http.Handle(metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>MaxScale Exporter</title></head>
			<body>
			<h1>MaxScale Exporter</h1>
			<p><a href="` + metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})
	log.Printf("Started MaxScale exporter, listening on port: %v", maxScaleExporterPort)
	log.Fatal(http.ListenAndServe(localIP+":"+maxScaleExporterPort, nil))
}
