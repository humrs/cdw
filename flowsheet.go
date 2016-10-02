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

import (
	"log"
	"strconv"
	"time"
)

// FlowsheetEvent holds a single from the flowsheet table
type FlowsheetEvent struct {
	PatientID       int
	RecordedTime    time.Time
	ObsDocGUID      int
	ClientVisitGUID int
	UpdatedTime     time.Time
	Location
	FlowsheetName    string
	ItemName         string
	ItemDescription  string
	ObservationValue string
}

// FromCSVRecord is a function which loads a flowsheet event from a CSV Record
func (f *FlowsheetEvent) FromCSVRecord(record []string) {
	var err error

	if len(record) != 10 {
		log.Fatal("Problem with Record Length\n")
	}

	f.PatientID, err = strconv.Atoi(record[0])
	if err != nil {
		log.Fatalf("PatientID: %v", err)
	}

	f.RecordedTime, err = time.Parse(Timeformat, record[1])
	if err != nil {
		log.Fatalf("RecordedDTM: %v", err)
	}

	f.ObsDocGUID, err = strconv.Atoi(record[2])
	if err != nil {
		log.Fatalf("ObsDocGUID: %v", err)
	}

	f.ClientVisitGUID, err = strconv.Atoi(record[3])
	if err != nil {
		log.Fatalf("ClientVisitGUID: %v", err)
	}

	f.UpdatedTime, err = time.Parse(Timeformat, record[4])
	if err != nil {
		log.Fatalf("UpdatedTime: %v", err)
	}

	f.Location.FromFlowsheetCSVRecord(record[5])

	f.FlowsheetName = record[6]
	f.ItemName = record[7]
	f.ItemDescription = record[8]
	f.ObservationValue = record[9]
}

// FromCSVRecordToLocationFiltered is a function which loads a flowsheet event from a CSV Record
func (f *FlowsheetEvent) FromCSVRecordToLocationFiltered(record []string) {
	var err error

	if len(record) != 10 {
		log.Fatal("Problem with Record Length\n")
	}

	f.PatientID, err = strconv.Atoi(record[0])
	if err != nil {
		log.Fatalf("PatientID: %v", err)
	}

	f.RecordedTime, err = time.Parse(Timeformat, record[1])
	if err != nil {
		log.Fatalf("RecordedDTM: %v", err)
	}

	f.UpdatedTime, err = time.Parse(Timeformat, record[4])
	if err != nil {
		log.Fatalf("UpdatedTime: %v", err)
	}

	f.Location.FromFlowsheetCSVRecord(record[5])
}

// ToLocationFilteredCSVRecord filters only the location and some identifiers
func (f *FlowsheetEvent) ToLocationFilteredCSVRecord() (record []string) {
	record = append(record,
		strconv.Itoa(f.PatientID),
		f.RecordedTime.Format(Timeformat),
		strconv.Itoa(f.ClientVisitGUID),
		f.LocationCode,
		f.Room,
		f.Bed)
	return
}
