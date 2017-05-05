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
package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/rstanleyhum/cdw"
)

// createadmissions creates an executable which you can pipe a csv file and outputs a csv file
//  with the picu admissions listed as a csv record which has a variable number of entries
//  Each line is a hospital admission and each line has several picu admissions (looks like it is
//  max of 8 for 2011-2015 admissions). Each admission gives a location, admit time, discharge time,
//  duration, and an optional bounceback time.
//
// Usage: createadmissions < input.csv > outfile.csv
//
func main() {
	log.SetOutput(os.Stderr)

	fp := os.Stdin
	outfp := os.Stdout

	chonyStays := make(map[cdw.UniquePatientID]*cdw.HospitalAdmit)

	hospitalTracker := cdw.NewTracker()

	r := csv.NewReader(fp)
	r.Comma = '|'
	firstRow := true

	csvWriter := csv.NewWriter(outfp)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if firstRow {
			firstRow = false
			continue
		}

		var event cdw.Event
		event.FromCSVRecord(record)
		admission, ok := hospitalTracker.ProcessEvent(event)
		if ok {
			if _, ok := chonyStays[admission.AdmitEvent.UniquePatientID]; !ok {
				var stay cdw.HospitalAdmit
				stay = cdw.NewHospitalAdmit()
				chonyStays[admission.AdmitEvent.UniquePatientID] = &stay
			}
			chonyStays[admission.AdmitEvent.UniquePatientID].AddAdmit(admission)
		}
	}

	for k := range chonyStays {
		if !chonyStays[k].HasPICUAdmit {
			continue
		}
		if err := csvWriter.Write(chonyStays[k].ToCSVRecord()); err != nil {
			log.Fatal(err)
		}
	}

	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		log.Fatal(err)
	}
}
