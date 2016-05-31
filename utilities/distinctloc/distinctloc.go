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
	"github.com/humrs/cdw"
	"sort"
)

// distinctloc takes a csv file and creates a file distinct locations (list)
//
//  Usage: distinctloc < infile.csv > outfile.csv
//
func main() {
	log.SetOutput(os.Stderr)
	
	fp := os.Stdin
	outfp := os.Stdout

	locations := make(map[cdw.Location] int)
	
	r := csv.NewReader(fp)
	firstRow := true
    
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
		loc := event.Location
		
		locations[loc]++
	}

	keys := make([]cdw.Location, len(locations))
	i := 0
	for k := range locations {
		keys[i] = k
		i++
	}
	
	sort.Sort(cdw.ByLocation(keys))
	
	csvWriter := csv.NewWriter(outfp)
	for _, v := range keys {
		if err := csvWriter.Write(v.ToCSVRecord()); err != nil {
			log.Fatal(err)
		}
		csvWriter.Flush()
	}
	
	if err := csvWriter.Error(); err != nil {
		log.Fatal(err)
	}
}
