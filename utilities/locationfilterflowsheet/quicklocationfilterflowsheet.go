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
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
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

	w := bufio.NewWriter(outfp)

	defer func() {
		fp.Close()
		outfp.Close()
	}()

	total := 0

	scanner := bufio.NewScanner(fp)

	var line string
	var items []string
	var newrow []string
	var newline string

	for scanner.Scan() {
		line = scanner.Text()
		total++
		items = strings.Split(line, "|")
		newrow = []string{
			items[0],
			items[1],
			items[4],
			items[5],
			items[6],
		}
		newline = strings.Join(newrow, "|")

		_, err = w.WriteString(newline)
		if err != nil {
			log.Fatalf("Writelin")
		}
		_, err = w.WriteString("\n")
		if err != nil {
			log.Fatalf("Writeln newline char")
		}

		w.Flush()
	}

	fmt.Printf("Total Count = %v\n", total)
}
