// Copyright © 2021 Elis Lulja
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

package pollresult

// Type represents the type of the event
type Type string

const (
	// EventSuccess represents a sucess
	EventSuccess Type = "success"
	// EventUnexpected means that something unexpected
	// happened or that last scrape contained unexpected
	// data.
	EventUnexpected Type = "unexpected"
	// EventError means that an error occurred during last
	// poll
	EventError Type = "error"
	// ... Define any other event type
)

// Event holds data about the poll
type Event struct {
	// Type of the event
	Type Type `json:"type" yaml:"type"`
	// WebsiteName is the name of the website that was polled
	WebsiteName string `json:"websiteName" yaml:"websiteName"`
	// ID is the ID of the page the was polled
	ID string `json:"id" yaml:"id"`
	// Message for the event
	Message string `json:"message" yaml:"message"`
	// Data to give more information about this event
	Data map[string]string `json:"data" yaml:"data"`
	// ... Define any other data relevant to you
}
