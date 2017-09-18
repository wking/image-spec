// Copyright 2016 The Linux Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"encoding/json"
	"fmt"
)

const (
	// ImageLayoutFile is the file name of oci image layout file
	ImageLayoutFile = "oci-layout"
	// ImageLayoutVersion is the version of ImageLayout
	ImageLayoutVersion = "1.1.0"
)

type EngineConfig struct {
	Protocol string `json:"protocol"`
	Data     map[string]interface{}
}

// ImageLayout is the structure in the "oci-layout" file, found in the root
// of an OCI Image-layout directory.
type ImageLayout struct {
	Version string `json:"imageLayoutVersion"`

	RefEngines []EngineConfig `json:"refEngines,omitempty"`
	CasEngines []EngineConfig `json:"casEngines,omitempty"`
}

func (c *EngineConfig) UnmarshalJSON(b []byte) (err error) {
	var dataInterface interface{}
	if err := json.Unmarshal(b, &dataInterface); err != nil {
		return err
	}

	data, ok := dataInterface.(map[string]interface{})
	if !ok {
		return fmt.Errorf("engine config is not a JSON object: %v", dataInterface)
	}

	protocolInterface, ok := data["protocol"]
	if !ok {
		return fmt.Errorf("engine config missing required 'protocol' entry: %v", data)
	}

	c.Protocol, ok = protocolInterface.(string)
	if !ok {
		return fmt.Errorf("engine config protocol is not a string: %v", protocolInterface)
	}

	delete(data, "protocol")
	c.Data = data
	return nil
}

func (c EngineConfig) MarshalJSON() ([]byte, error) {
	var data map[string]interface{}
	data = make(map[string]interface{})
	for key, value := range c.Data {
		data[key] = value
	}
	data["protocol"] = c.Protocol
	return json.Marshal(data)
}
