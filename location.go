// Copyright 2016 R. Stanley Hum
//
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
package cdw

import "strings"

// Location is the structure for holding distinct locations
type Location struct {
	LocationCode string
	Room         string
	Bed          string
}

// ByLocation is the sorting algorithm for Location
type ByLocation []Location

func (a ByLocation) Len() int      { return len(a) }
func (a ByLocation) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByLocation) Less(i, j int) bool {
	if a[i].LocationCode < a[j].LocationCode {
		return true
	}

	if a[i].LocationCode > a[j].LocationCode {
		return false
	}

	if a[i].Room < a[j].Room {
		return true
	}

	if a[i].Room > a[j].Room {
		return false
	}

	if a[i].Bed < a[j].Bed {
		return true
	}

	return false
}

// ToCSVRecord converts back to CSV record
func (l *Location) ToCSVRecord() (record []string) {
	record = append(record,
		l.LocationCode,
		l.Room,
		l.Bed)
	return
}

// FromFlowsheetCSVRecord parses the location field in flowsheet data
func (l *Location) FromFlowsheetCSVRecord(record string) {
	var locationstrings = strings.Split(record, "-")
	if locationstrings[0] == "" {
		return
	}

	l.LocationCode = locationstrings[0]

	if len(locationstrings) > 1 {
		l.Room = locationstrings[1]
	}

	if len(locationstrings) > 2 {
		l.Bed = locationstrings[2]
	}

	return
}
