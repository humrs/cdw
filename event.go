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

// Event holds the information for a ADT event
type Event struct {
	UniquePatientID
	PrimaryTime time.Time
	EventCode   int
	Location
	MedicalServiceCode string
	ProviderID         string
	SequenceNo         string
	AlternateID        string
	SqNum              int
}

// ByPrimaryTime is the sorting algorithm for Events
type ByPrimaryTime []Event

func (a ByPrimaryTime) Len() int           { return len(a) }
func (a ByPrimaryTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPrimaryTime) Less(i, j int) bool { return a[i].PrimaryTime.Before(a[j].PrimaryTime) }

// FromCSVRecord is a function which loads an event from a CSV record
func (e *Event) FromCSVRecord(record []string) {
	var err error

	if len(record) != 12 {
		log.Fatalf("Problem with Record Length\n")
	}

	e.Account, err = strconv.Atoi(record[0])
	if err != nil {
		log.Fatalf("Account: %v", err)
	}
	e.PatientID, err = strconv.Atoi(record[1])
	if err != nil {
		log.Fatalf("PatientID: %v", err)
	}
	e.PrimaryTime, err = time.Parse(Timeformat, record[2])
	if err != nil {
		log.Fatalf("Time: %v", err)
	}
	e.EventCode, err = strconv.Atoi(record[3])
	if err != nil {
		log.Fatalf("EventCode: %v", err)
	}
	e.LocationCode = record[4]
	e.Room = record[5]
	e.Bed = record[6]
	e.MedicalServiceCode = record[7]
	e.ProviderID = record[8]
	e.SequenceNo = record[9]
	e.AlternateID = record[10]
	e.SqNum, err = strconv.Atoi(record[11])
	if err != nil {
		log.Fatalf("SQ Num: %v", err)
	}
}

// ToCSVRecord converts back to CSV record
func (e *Event) ToCSVRecord() (record []string) {
	record = append(record,
		strconv.Itoa(e.Account),
		strconv.Itoa(e.PatientID),
		e.PrimaryTime.Format(Timeformat),
		strconv.Itoa(e.EventCode),
		e.LocationCode,
		e.Room,
		e.Bed,
		e.MedicalServiceCode,
		e.ProviderID,
		e.SequenceNo,
		e.AlternateID,
		strconv.Itoa(e.SqNum))
	return
}
