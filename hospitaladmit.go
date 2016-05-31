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
	"errors"
    "strconv"
)

// HospitalAdmit is an inpatient stay
type HospitalAdmit struct {
	admits        []UnitAdmit
    HasPICUAdmit    bool
}

// NewHospitalAdmit is the HospitalAdmit constructor
func NewHospitalAdmit() HospitalAdmit {
    var result HospitalAdmit
    result.HasPICUAdmit = false
    return result
}

// GetPatientID returns the Patient ID
func (h HospitalAdmit) GetPatientID() (patientID int, returnOk bool) {
	if len(h.admits) == 0 {
        return patientID, false
    }
    patientID = h.admits[0].GetPatientID()
    returnOk = true
	return
}

// GetAccount returns the Account
func (h HospitalAdmit) GetAccount() (account int, returnOk bool) {
    if len(h.admits) == 0 {
        return account, false
    }
    account = h.admits[0].GetAccount()
    returnOk = true
	return
}

// GetLocations returns a list of locations
func (h HospitalAdmit) GetLocations() (locs []Location, returnOk bool) {
    numAdmits := len(h.admits)
    if numAdmits == 0 {
        return locs, false
    }
    locs = make([]Location, numAdmits, numAdmits)
    for i, v := range h.admits {
        locs[i] = v.GetLocation()
    }
    returnOk = true
    return
}

// GetUnits returns all of the units in a stay
func (h HospitalAdmit) GetUnits() (units []string, returnOk bool) {
    locs, ok := h.GetLocations()
    if !ok {
        return units, false
    }
    numLocs := len(locs)
    units = make([]string, numLocs, numLocs)
    for i, v := range locs {
        units[i] = v.LocationCode
    }
    if len(units) == 0 {
        returnOk = false
    } else {
        returnOk = true
    }
    return
}

// isPICUUnit checks if location code is a PICU or PCICU
func isPICUUnit(u string) bool {
    switch u {
        case "B09N":
            return true
        case "B09S":
            return true
        case "B09C":
            return true
        case "B09T":
            return true
        case "B11C":
            return true
    }
    return false
}

// GetPICUAdmits returns UnitSummaries about PICU Admits
func (h HospitalAdmit) GetPICUAdmits() (admits []UnitSummary) {
    isInPICU := false
    if len(h.admits) == 0 {
        return
    }
    
    var currentPICUAdmit UnitAdmit
    
    for _, v := range h.admits {
        loc := v.AdmitEvent.LocationCode
        if isPICUUnit(loc) {
            if !isInPICU {
                isInPICU = true
                currentPICUAdmit = v
            }
        } else {
            if isInPICU {
                var summary UnitSummary
                summary.Account = currentPICUAdmit.GetAccount()
                summary.PatientID = currentPICUAdmit.GetPatientID()
                summary.AdmitTime = currentPICUAdmit.GetAdmitTime()
                summary.DischargeTime = v.GetAdmitTime()
                summary.Location = currentPICUAdmit.GetLocation()
                admits = append(admits, summary)
            }
            isInPICU = false
        }
    }
    if isInPICU {
        var summary UnitSummary
        summary.Account = currentPICUAdmit.GetAccount()
        summary.PatientID = currentPICUAdmit.GetPatientID()
        summary.AdmitTime = currentPICUAdmit.GetAdmitTime()
        summary.DischargeTime = currentPICUAdmit.GetDischargeTime()
        summary.Location = currentPICUAdmit.GetLocation()
        admits = append(admits, summary)
    }
    return
}

// AddAdmit adds an admisison to the hospital stay
func (h *HospitalAdmit) AddAdmit(a UnitAdmit) error {
    if isPICUUnit(a.AdmitEvent.LocationCode) {
        h.HasPICUAdmit = true
    }
    
    if len(h.admits) == 0 {
        h.admits = append(h.admits, a)
        return nil
    }
    
    account, ok := h.GetAccount()
    if !ok {
        return errors.New("Problem getting Account")
    }
    
    patientID, ok := h.GetPatientID()
    if !ok {
        return errors.New("Problem getting PatientID")
    }
    
    if account != a.GetAccount() || patientID != a.GetPatientID() {
        return errors.New("Account and PatientID does not match")
    }
    
    h.admits = append(h.admits, a)
    return nil
}

// ToCSVRecord creates CSV record for writing
func (h HospitalAdmit) ToCSVRecord() (record []string) {
    if len(h.admits) == 0 {
        return
    }
    isInPICU := false
    account, _ := h.GetAccount()
    patientID, _ := h.GetPatientID()
    record = append(record,
        strconv.Itoa(account),
        strconv.Itoa(patientID))
    
    var currentPICUAdmit UnitAdmit
    var lastPICUDischargeTime time.Time
    
    for _, v := range h.admits {
        loc := v.AdmitEvent.LocationCode
        if isPICUUnit(loc) {
            if !isInPICU {
                isInPICU = true
                currentPICUAdmit = v
            }
        } else {
            if isInPICU {
                records := createRecord(currentPICUAdmit.AdmitEvent, v.AdmitEvent, lastPICUDischargeTime)
                record = append(record, records...)
                lastPICUDischargeTime = v.AdmitEvent.PrimaryTime
            }
            isInPICU = false
        }
    }
    
    if isInPICU {
        records := createRecord(currentPICUAdmit.AdmitEvent, currentPICUAdmit.DischargeEvent, lastPICUDischargeTime)
        record = append(record, records...)
    }
    return
}

func createRecord(admit Event, discharge Event, lastPICUDischargeTime time.Time) (record []string) {
    admitTime := admit.PrimaryTime
    dischargeTime := discharge.PrimaryTime
    location := admit.LocationCode
    room := admit.Room
    bed := admit.Bed
    duration := dischargeTime.Sub(admitTime)
    record = append(record,
        location,
        room,
        bed,
        strconv.FormatFloat(duration.Hours(), 'f', 2, 32),
        admitTime.Format(Timeformat),
        dischargeTime.Format(Timeformat))
    if !lastPICUDischargeTime.IsZero() {
        bounceback := admitTime.Sub(lastPICUDischargeTime)
        record = append(record, strconv.FormatFloat(bounceback.Hours(), 'f', 2, 32))
    }
    return
}