// © 2026 Nokia
// Licensed under the Apache License, Version 2.0 (the "License");
// SPDX-License-Identifier: Apache-2.0
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
	"encoding/json"
	"errors"
	"strings"
)

// AutoFailoverValue is a custom type that can handle both bool and string values
// String values are expected to be one of: "true", "false", or "safe"
type AutoFailoverValue struct {
	BoolValue bool
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (a *AutoFailoverValue) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as bool first
	var boolValue bool
	if err := json.Unmarshal(data, &boolValue); err == nil {
		a.BoolValue = boolValue
		return nil
	}

	// If not a bool, try as string
	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err == nil {
		// Validate that the string is one of the expected enum values
		lowerValue := strings.ToLower(stringValue)
		// Set BoolValue to true for "true" or "safe", false otherwise
		a.BoolValue = lowerValue == "true" || lowerValue == "safe"
		return nil
	}

	// If neither, return error
	return errors.New("Cannot find AutoFailoverValue in JSON")
}

// MarshalJSON implements the json.Marshaler interface
func (a AutoFailoverValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.BoolValue)
}

// Bool returns the boolean value
func (a AutoFailoverValue) Bool() bool {
	return a.BoolValue
}
