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
	"bufio"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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

	locfp, err := os.Create("locationindex.txt")
	if err != nil {
		log.Fatalf("Unable to open location index file")
	}

	itemfp, err := os.Create("itemindex.txt")
	if err != nil {
		log.Fatalf("Unable to open item index file")
	}

	fsfp, err := os.Create("flowsheetindex.txt")
	if err != nil {
		log.Fatalf("Unable to open flowsheet index filename.")
	}

	obsfp, err := os.Create("C:\\Users\\stanley\\Desktop\\observations.txt")
	if err != nil {
		log.Fatalf("Unable to open observations index file")
	}

	outfp, err := os.Create(*outfilename)
	if err != nil {
		log.Fatalf("Unable to open output filename.")
	}

	enc := gob.NewEncoder(outfp)

	defer func() {
		fp.Close()
		obsfp.Close()
		locfp.Close()
		itemfp.Close()
		fsfp.Close()
		outfp.Close()
	}()

	total := 0

	scanner := bufio.NewScanner(fp)

	var line string
	var items []string
	isFirst := true

	for scanner.Scan() {
		line = scanner.Text()

		if isFirst {
			isFirst = false
			continue
		}

		total++
		items = strings.Split(line, "|")

		var binaryFSevent cdw.BinaryFlowsheetEvent
		binaryFSevent.FromLine(items[:9], total)

		observationItem := strings.Join(items[9:len(items)], "|")
		_, _ = obsfp.WriteString(strconv.Itoa(total))
		_, _ = obsfp.WriteString("|")
		_, _ = obsfp.WriteString(observationItem)

		err = enc.Encode(binaryFSevent)
		if err != nil {
			log.Fatalf("Encoding problem: %v", err)
		}
	}

	data, _ := json.Marshal(cdw.LocationIndex)
	locfp.Write(data)

	data, _ = json.Marshal(cdw.ItemIndex)
	itemfp.Write(data)

	data, _ = json.Marshal(cdw.FlowsheetIndex)
	fsfp.Write(data)

	fmt.Printf("Total Count = %v\n", total)
}
