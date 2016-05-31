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

// Timeformat is the constant for the format of the csv time
const Timeformat = "2006-01-02 15:04:05.000000"

// Event Codes
const (
    AdmissionEvaluationCode           = 32467
    TransferAPatientCode              = 32468
    DischargeAPatientCode             = 32469
    TransferOutpatientToInpatientCode = 32472
    SwapAPatientCode                  = 32483
    AmbulatorySurgeryVisitCode        = 48677
    EmergencyRoomVisitCode            = 48678
    DpoVisitCode                      = 48679
    ClinicVisitCode                   = 48661
    TherapeuticEncounterCode          = 56317
)
