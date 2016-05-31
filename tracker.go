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

// Tracker is an inpatient environment tracker
type Tracker struct {
    openEvents map[UniquePatientID] Event
}

// NewTracker is a constructor for Tracker type
func NewTracker() *Tracker {
    return &Tracker{openEvents: make(map[UniquePatientID] Event)}
}

// ProcessEvent takes an Event and creates an UnitAdmit if possible
func (t *Tracker) ProcessEvent(e Event) (a UnitAdmit, returnOk bool){
    returnOk = false
    if initialEvent, ok := t.openEvents[e.UniquePatientID]; ok {
        a = NewUnitAdmit(initialEvent, e)
        returnOk = true
    }
    t.setEvent(e)
    return a, returnOk
}

// setEvent takes an Event and puts it into the proper place in the openEvents in the environment
func (t *Tracker) setEvent(e Event) {
    if e.EventCode != DischargeAPatientCode {
        t.openEvents[e.UniquePatientID] = e
    } else {
        delete(t.openEvents, e.UniquePatientID) 
    }
}