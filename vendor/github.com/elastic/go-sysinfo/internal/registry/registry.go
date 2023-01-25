// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package registry

import (
	"fmt"

	"github.com/elastic/go-sysinfo/types"
)

var (
	hostProvider    HostProvider
	processProvider ProcessProvider
)

type HostProvider interface {
	Host() (types.Host, error)
}

type ProcessProvider interface {
	Processes() ([]types.Process, error)
	Process(pid int) (types.Process, error)
	Self() (types.Process, error)
}

func Register(provider interface{}) {
	if h, ok := provider.(HostProvider); ok {
		if hostProvider != nil {
			panic(fmt.Sprintf("HostProvider already registered: %v", hostProvider))
		}
		hostProvider = h
	}

	if p, ok := provider.(ProcessProvider); ok {
		if processProvider != nil {
			panic(fmt.Sprintf("ProcessProvider already registered: %v", processProvider))
		}
		processProvider = p
	}
}

func GetHostProvider() HostProvider       { return hostProvider }
func GetProcessProvider() ProcessProvider { return processProvider }
