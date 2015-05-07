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

package mojang

// Profile contains information about the player required
// to connect to a server
type Profile struct {
	Username    string
	ID          string
	AccessToken string
}

// IsComplete returns whether the profile is enough to connect
// with.
func (p Profile) IsComplete() bool {
	return p.Username != "" && p.ID != "" && p.AccessToken != ""
}
