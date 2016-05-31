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
    "time"
	"strconv"
)

// UnitAdmit is stores unit admission data
type UnitAdmit struct {
	AdmitEvent         Event
	DischargeEvent     Event
}

// NewUnitAdmit creates a new UnitAdmit
func NewUnitAdmit(a Event, d Event) (result UnitAdmit) {
	result.AdmitEvent = a
	result.DischargeEvent = d
	return
}

// GetPatientID returns the Patient ID
func (a *UnitAdmit) GetPatientID() (PatientID int) {
	PatientID = a.AdmitEvent.PatientID
	return
}

// GetAccount returns the Account
func (a *UnitAdmit) GetAccount() (Account int) {
	Account = a.AdmitEvent.Account
	return
}

// GetAdmitTime returns the time of UnitAdmit event
func (a *UnitAdmit) GetAdmitTime() (admitTime time.Time) {
	admitTime = a.AdmitEvent.PrimaryTime
	return
}

// GetDischargeTime returns the time of discharge event
func (a *UnitAdmit) GetDischargeTime() (dischargeTime time.Time) {
	dischargeTime = a.DischargeEvent.PrimaryTime
	return
}

// GetAdmitEventCode returns the event code for admit event
func (a *UnitAdmit) GetAdmitEventCode() (EventCode int) {
	EventCode = a.AdmitEvent.EventCode
	return
}

// GetLocation returns the location of the admit event
func (a *UnitAdmit) GetLocation() (loc Location) {
	loc = a.AdmitEvent.Location
	return
}

// ToCSVRecord converts back to CSV record
func (a *UnitAdmit) ToCSVRecord() (record []string) {
    record = append(record,
        strconv.Itoa(a.AdmitEvent.Account),
        strconv.Itoa(a.AdmitEvent.PatientID),
        a.AdmitEvent.PrimaryTime.Format(Timeformat),
		strconv.Itoa(a.AdmitEvent.EventCode),
        a.AdmitEvent.LocationCode,
        a.AdmitEvent.Room,
        a.AdmitEvent.Bed,
		strconv.Itoa(a.DischargeEvent.EventCode))
	return
}