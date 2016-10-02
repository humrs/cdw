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
	"bytes"
	"encoding/binary"
	"log"
	"strconv"
	"time"
)

// BinaryFlowsheetEvent holds a single from the flowsheet table
type BinaryFlowsheetEvent struct {
	PatientID        int32
	RecordedTimeNano int64
	ObsDocGUID       int32
	ClientVisitGUID  int32
	UpdatedTimeNano  int64
	LocationNo       int32
	FlowsheetNo      int32
	ItemNo           int32
	ObservationNo    int32
}

var LocIndex *LocationIndex

var ItmIndex *ItemIndex

var FSIndex *FlowsheetIndex

// ToBytes gives the BinaryFlowsheetEvent in bytes
func (b BinaryFlowsheetEvent) ToBytes() (result []byte) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.BigEndian, b)
	if err != nil {
		log.Fatalf("ToBytes: %v\n", err)
	}
	result = buf.Bytes()
	return
}

// FromLine converts the scanned Text to a BinaryFlowsheetEvent
func (b *BinaryFlowsheetEvent) FromLine(items []string, observationNo int) {
	var err error
	var value int

	if len(items) != 9 {
		log.Fatalf("Problem with Record length: %v: %v", len(items), items)
	}

	value, err = strconv.Atoi(items[0][1 : len(items[0])-1])
	if err != nil {
		log.Fatalf("PatientID: %v", err)
	}
	b.PatientID = int32(value)

	recordedtime, err := time.Parse(Timeformat, items[1][1:len(items[1])-1])
	if err != nil {
		log.Fatalf("RecordedDTM: %v", err)
	}
	b.RecordedTimeNano = recordedtime.UnixNano()

	value, err = strconv.Atoi(items[2][1 : len(items[2])-1])
	if err != nil {
		log.Fatalf("ObsDocGUID: %v", err)
	}
	b.ObsDocGUID = int32(value)

	value, err = strconv.Atoi(items[3][1 : len(items[3])-1])
	if err != nil {
		log.Fatalf("ClientVisitGUID: %v", err)
	}
	b.ClientVisitGUID = int32(value)

	updatedtime, err := time.Parse(Timeformat, items[4][1:len(items[4])-1])
	if err != nil {
		log.Fatalf("UpdatedTime: %v", err)
	}
	b.UpdatedTimeNano = updatedtime.UnixNano()

	b.LocationNo = LocIndex.Get(items[5])
	b.FlowsheetNo = FSIndex.Get(items[6])

	fullitemname := joinItems(items[7], items[8])
	b.ItemNo = ItmIndex.Get(fullitemname)
	b.ObservationNo = int32(observationNo)

}

func joinItems(a string, b string) (result string) {
	var buffer bytes.Buffer
	buffer.WriteString(a)
	buffer.WriteString("|")
	buffer.WriteString(b)
	result = buffer.String()
	return
}
