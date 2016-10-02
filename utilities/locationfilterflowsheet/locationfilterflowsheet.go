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
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/humrs/cdw"
)

// locationfilterflowsheet takes a csv file and sorts the adt events and outputs a csv file
//
//  Usage: locationfilterflowsheet < infile.csv > outfile.csv
//
func main() {

	infilename := flag.String("infile", "", "Input Filename")

	outfilename := flag.String("outfile", "", "Outputfilename")

	flag.Parse()

	fmt.Println("infile: ", *infilename)
	fmt.Println("outfile: ", *outfilename)

	fp, err := os.Open(*infilename)
	if err != nil {
		log.Fatalf("Unable to open input filename.")
	}

	outfp, err := os.Create(*outfilename)
	if err != nil {
		log.Fatalf("Unable to open output filename.")
	}

	defer func() {
		fp.Close()
		outfp.Close()
	}()

	r := csv.NewReader(fp)
	r.Comma = '|'
	total := 0
	isFirstRow := true
	finished := false
	csvWriter := csv.NewWriter(outfp)

	for {
		fsevents := make(map[cdw.FlowsheetEvent]int)

		for i := 0; i < 10000000; i++ {

			record, err := r.Read()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Reading problem: %v", err)
			}

			if isFirstRow {
				isFirstRow = false
				continue
			}

			var fsevent cdw.FlowsheetEvent
			fsevent.FromCSVRecordToLocationFiltered(record)
			fsevents[fsevent]++
		}

		for k := range fsevents {
			if err := csvWriter.Write(k.ToLocationFilteredCSVRecord()); err != nil {
				log.Fatalf("Writing problem: %v", err)
			}
		}
		csvWriter.Flush()

		fmt.Printf("written: %v\n", len(fsevents))

		total = total + len(fsevents)

		if err = csvWriter.Error(); err != nil {
			log.Fatalf("Writer error: %v", err)
		}

		if finished {
			break
		}
	}

	fmt.Println("Total Count = %v", total)
}
