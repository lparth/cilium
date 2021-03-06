// Copyright 2016-2018 Authors of Cilium
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

package option

var (
	specPolicyTracing = Option{
		Description: "Enable tracing when resolving policy (Debug)",
	}

	daemonLibrary = OptionLibrary{
		PolicyTracing: &specPolicyTracing,
	}

	DaemonMutableOptionLibrary = OptionLibrary{
		ConntrackAccounting: &specConntrackAccounting,
		ConntrackLocal:      &specConntrackLocal,
		Conntrack:           &specConntrack,
		Debug:               &specDebug,
		DebugLB:             &specDebugLB,
		DropNotify:          &specDropNotify,
		TraceNotify:         &specTraceNotify,
		NAT46:               &specNAT46,
	}
)

func init() {
	for k, v := range DaemonMutableOptionLibrary {
		daemonLibrary[k] = v
	}
}

// ParseDaemonOption parses a string as daemon option
func ParseDaemonOption(opt string) (string, bool, error) {
	return ParseOption(opt, &daemonLibrary)
}
