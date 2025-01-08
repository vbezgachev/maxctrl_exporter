// Copyright 2023 Dylan Northrup [@dylan-tock on github]
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
	"fmt"
	"log"
	"os"
	"testing"
)

// Verify GetEnvVar does what it says on the tin
func TestGetInitializedEnvVar(t *testing.T) {
	// Save possibly pre-existing value of env variable for restoration later
	envVariableName := "maxctrlExporterTestEnvVariable"
	previousValue := os.Getenv(envVariableName)
	defer func() { os.Setenv(envVariableName, previousValue) }()
	fallback := "fallback value"
	envValue := "Test Value"
	want := "Test Value"

	os.Setenv(envVariableName, envValue)
	got := GetEnvVar(envVariableName, fallback)

	if want != got {
		t.Fatalf("Did not get expected value from env variable. Wanted: '%s'. Got '%s'", want, got)
	}
	if previousValue == "" {
		os.Unsetenv(envVariableName)
	} else {
		os.Setenv(envVariableName, previousValue)
	}
}

func TestGetEmptyEnvVar(t *testing.T) {
	// Save possibly pre-existing value of env variable for restoration later
	envVariableName := "maxctrlExporterTestEnvVariable"
	previousValue := os.Getenv(envVariableName)
	defer func() { os.Setenv(envVariableName, previousValue) }()
	fallback := "fallback value"
	envValue := ""
	want := fallback

	os.Setenv(envVariableName, envValue)
	got := GetEnvVar(envVariableName, fallback)

	if want != got {
		t.Fatalf("Did not get expected value from env variable. Wanted: '%s'. Got '%s'", want, got)
	}
	if previousValue == "" {
		os.Unsetenv(envVariableName)
	} else {
		os.Setenv(envVariableName, previousValue)
	}
}

func TestGettingConfigFromEnvironment(t *testing.T) {
	maxScaleUrl = ""
	maxScaleUsername = ""
	maxScalePassword = ""
	maxScaleExporterPort = ""
	maxScaleCACertificate = ""
	maxctrlExporterConfigFile = ""
	maxScaleMaxConnections = ""

	keys := []string{"MAXSCALE_URL", "MAXSCALE_USERNAME", "MAXSCALE_PASSWORD", "MAXSCALE_EXPORTER_PORT",
		"MAXSCALE_CA_CERTIFICATE", "MAXCTRL_EXPORTER_CFG_FILE", "MAXSCALE_MAX_CONNECTIONS"}

	want := map[string]string{
		keys[0]: "http://10.10.10.1:8989",
		keys[1]: "userMcUserFace",
		keys[2]: "secretPassword",
		keys[3]: "8080",
		keys[4]: "cert.pem",
		keys[5]: "exporterConfig.yml",
		keys[6]: "",
	}

	for k, v := range want {
		os.Setenv(k, v)
	}

	setConfigFromEnvironmentVars()
	got := make(map[string]string)
	got[keys[0]] = maxScaleUrl
	got[keys[1]] = maxScaleUsername
	got[keys[2]] = maxScalePassword
	got[keys[3]] = maxScaleExporterPort
	got[keys[4]] = maxScaleCACertificate
	got[keys[5]] = maxctrlExporterConfigFile
	got[keys[6]] = maxScaleMaxConnections

	for _, k := range keys {
		if want[k] != got[k] {
			log.Fatalf("Config key '%s' had unexpected value. wanted '%s' and got '%s'", k, want[k], got[k])
		}
		// Unset these as we test
		os.Unsetenv(k)
	}

	for k, v := range want {
		if k == "MAXSCALE_EXPORTER_PORT" || k == "MAXCTRL_EXPORTER_CFG_FILE" {
			continue
		}
		os.Setenv(k, v)
	}

	for _, k := range keys {
		if want[k] != got[k] {
			log.Fatalf("Config key '%s' had unexpected value. wanted '%s' and got '%s'", k, want[k], got[k])
		}
		// Unset these as we test
		os.Unsetenv(k)
	}
}

// Test parsing config contents
func TestConfigParsing(t *testing.T) {
	// Pre-initialize the variables
	setConfigFromEnvironmentVars()

	keys := []string{"url", "username", "password", "exporter_port", "caCertificate", "maxConnections"}

	want := map[string]string{
		keys[0]: "http://10.10.10.1:8989",
		keys[1]: "userMcUserFace",
		keys[2]: "secretPassword",
		keys[3]: "8080",
		keys[4]: "",
		keys[5]: "",
	}

	contents := ""
	for k, v := range want {
		contents = fmt.Sprintf("%s%s: %s\n", contents, k, v)
	}

	// We have to call this to set up proper default values
	parseConfigFile([]byte(contents))

	got := make(map[string]string)
	got[keys[0]] = maxScaleUrl
	got[keys[1]] = maxScaleUsername
	got[keys[2]] = maxScalePassword
	got[keys[3]] = maxScaleExporterPort
	got[keys[4]] = maxScaleCACertificate
	got[keys[5]] = maxScaleMaxConnections

	for _, k := range keys {
		if want[k] != got[k] {
			log.Fatalf("Config key '%s' had unexpected value. wanted '%s' and got '%s'", k, want[k], got[k])
		}
	}

	// Redo the test, but with some values missing from the config file
	contents = ""
	for k, v := range want {
		if k == "exporter_port" || k == "caCertificate" || k == "maxConnections" {
			continue
		}
		contents = fmt.Sprintf("%s%s: %s\n", contents, k, v)
	}

	// We have to call this to set up proper default values
	setConfigFromEnvironmentVars()
	parseConfigFile([]byte(contents))

	got[keys[0]] = maxScaleUrl
	got[keys[1]] = maxScaleUsername
	got[keys[2]] = maxScalePassword
	got[keys[3]] = maxScaleExporterPort
	got[keys[4]] = maxScaleCACertificate
	got[keys[5]] = maxScaleMaxConnections

	for _, k := range keys {
		if want[k] != got[k] {
			log.Fatalf("Config key '%s' had unexpected value. wanted '%s' and got '%s'", k, want[k], got[k])
		}
	}
}
