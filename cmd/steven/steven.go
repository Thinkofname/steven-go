// Copyright 2015 Matthew Collins
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

package main

import (
	"os"
	"runtime"

	"github.com/thinkofdeath/steven"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	// Can't use flags as we need to support a weird flag
	// format
	var username, uuid, accessToken, server string

	for i, arg := range os.Args {
		switch arg {
		case "--username":
			username = os.Args[i+1]
		case "--uuid":
			uuid = os.Args[i+1]
		case "--accessToken":
			accessToken = os.Args[i+1]
		case "--server":
			server = os.Args[i+1]
		}
	}
	steven.Main(username, uuid, accessToken, server)
}
