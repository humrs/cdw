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
	"sort"
	"time"
	"github.com/humrs/cdw"
)

var events []cdw.Event

// sortevents takes a csv file and sorts the adt events and outputs a csv file
//
//  Usage: sortevents < infile.csv > outfile.csv
//
func main() {
	log.SetOutput(os.Stderr)
	
	fp := os.Stdin
	outfp := os.Stdout

	r := csv.NewReader(fp)
	isFirstRow := true

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		if isFirstRow {
			isFirstRow = false
			continue
		}

		var event cdw.Event
		event.FromCSVRecord(record)
		events = append(events, event)
	}

	sort.Sort(cdw.ByPrimaryTime(events))

	oldtime, err := time.Parse(cdw.Timeformat, cdw.Timeformat)

	for _, v := range events {
		if v.PrimaryTime.After(oldtime) || v.PrimaryTime.Equal(oldtime) {
			oldtime = v.PrimaryTime
		} else {
			log.Fatalf("Needs sorting again: Primary Time = %v, oldtime = %v\n", v.PrimaryTime, oldtime)
		}
	}

	csvWriter := csv.NewWriter(outfp)
	for _, v := range events {
		if err := csvWriter.Write(v.ToCSVRecord()); err != nil {
			log.Fatal(err)
		}
		csvWriter.Flush()
	}

	if err = csvWriter.Error(); err != nil {
		log.Fatal(err)
	}

}
