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
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/humrs/cdw"
)

type inmessage struct {
	number int
	line   string
}

type outmessage struct {
	fsevent     cdw.BinaryFlowsheetEvent
	number      int
	observation string
}

// newefficientfilter takes a csv file and sorts the adt events and outputs a csv file
//
//  Usage: efficientfilter < infile.csv > outfile.csv
//
func main() {
	out := make(chan outmessage)

	go save(out)

	cdw.LocIndex = cdw.NewLocationIndex()
	cdw.ItmIndex = cdw.NewItemIndex()
	cdw.FSIndex = cdw.NewFlowsheetIndex()

	infilename := flag.String("infile", "", "Input Filename")

	flag.Parse()

	fmt.Println("infile: ", *infilename)

	locfp, err := os.Create("locations.txt")
	if err != nil {
		log.Fatalf("Unable to open locations file")
	}

	itemfp, err := os.Create("items.txt")
	if err != nil {
		log.Fatalf("Unable to open items file")
	}

	fsfp, err := os.Create("flowsheet.txt")
	if err != nil {
		log.Fatalf("Unable to open flowsheet file")
	}

	defer func() {
		locfp.Close()
		itemfp.Close()
		fsfp.Close()
		close(out)
	}()

	fmt.Printf("Starting: %v\n", time.Now())

	frominfile := infileservice(*infilename)

	c1 := processLine(frominfile)
	c2 := processLine(frominfile)
	c3 := processLine(frominfile)
	c4 := processLine(frominfile)
	c5 := processLine(frominfile)
	c6 := processLine(frominfile)
	c7 := processLine(frominfile)
	c8 := processLine(frominfile)
	c9 := processLine(frominfile)
	c10 := processLine(frominfile)

	for n := range merge(c1, c2, c3, c4, c5, c6, c7, c8, c9, c10) {
		out <- n
	}

	fmt.Printf("Ending: %v\n", time.Now())

	locdata := cdw.LocIndex.GetValues()
	for k, v := range locdata {
		locfp.WriteString(strconv.Itoa(k))
		locfp.WriteString("|")
		locfp.WriteString(v)
		locfp.WriteString("\n")
	}

	itemdata := cdw.ItmIndex.GetValues()
	for k, v := range itemdata {
		itemfp.WriteString(strconv.Itoa(k))
		itemfp.WriteString("|")
		itemfp.WriteString(v)
		itemfp.WriteString("\n")
	}

	fsdata := cdw.FSIndex.GetValues()
	for k, v := range fsdata {
		fsfp.WriteString(strconv.Itoa(k))
		fsfp.WriteString("|")
		fsfp.WriteString(v)
		fsfp.WriteString("\n")
	}

}

func infileservice(filename string) <-chan inmessage {
	out := make(chan inmessage)

	go func() {
		fp, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Unable to open input file")
		}
		defer fp.Close()

		var largeBuf []byte
		buffersize := 1024 * 1024 * 1024
		var n int
		var total int64
		var pass int
		var lines int
		var linesCorrected int

		last := ""
		isFinished := false
		isFirst := true

		for {
			if isFinished {
				break
			}

			largeBuf = make([]byte, buffersize)

			n, err = fp.Read(largeBuf)
			if err != nil {
				log.Fatalf("Error in Reading: %v", err)
			}

			if n < buffersize {
				isFinished = true
			}

			total = total + int64(n)
			pass++
			lread := bytes.NewBuffer(largeBuf)
			for {
				line, err := lread.ReadString('\n')
				if err == io.EOF {
					last = line
					break
				}

				if isFirst {
					isFirst = false
					continue
				}

				if err != nil {
					log.Fatalf("First error: %v: %v", lines, err)
				}

				if last != "" {
					var buffer bytes.Buffer
					buffer.WriteString(last)
					buffer.WriteString(line)
					line = buffer.String()
					last = ""
					linesCorrected++
				}

				inmsg := inmessage{line: line, number: lines}
				out <- inmsg

				lines++
			}

			fmt.Printf("Pass: %v  Total: %v  Lines: %v  LinesCorrected: %v  Time: %v\n", pass, total, lines, linesCorrected, time.Now())

		}
		close(out)
	}()

	return out
}

func save(in <-chan outmessage) {
	outfp, err := os.Create("C:\\Users\\stanley\\Desktop\\fsindex.bin")
	if err != nil {
		log.Fatalf("Unable to open output filename.")
	}

	obsfp, err := os.Create("C:\\Users\\stanley\\Desktop\\observations.txt")
	if err != nil {
		log.Fatalf("Unable to open observation file.")
	}

	defer func() {
		outfp.Close()
		obsfp.Close()
	}()

	lineno := 0
	for i := range in {
		_, err := outfp.Write(i.fsevent.ToBytes())
		if err != nil {
			log.Fatalf("outfp: %v\n", err)
		}

		_, err = obsfp.WriteString(i.observation)
		if err != nil {
			log.Fatalf("obsfp: %v\n", err)
		}

		if (lineno % 10000000) == 0 {
			fmt.Printf("Lines so far: %v\n", lineno)
		}

		lineno++
	}

	fmt.Printf("Total lines: %v", lineno)

}

func processLine(in <-chan inmessage) <-chan outmessage {
	out := make(chan outmessage)

	go func() {
		for msg := range in {
			items := strings.Split(msg.line, "|")

			var binfsevent cdw.BinaryFlowsheetEvent
			binfsevent.FromLine(items[:9], msg.number)

			observationArray := []string{strconv.Itoa(msg.number)}
			observationArray = append(observationArray, items[9:len(items)]...)
			observationItem := strings.Join(observationArray, "|")

			outmsg := outmessage{fsevent: binfsevent, number: msg.number, observation: observationItem}
			out <- outmsg
		}
		close(out)
	}()

	return out
}

func merge(cs ...<-chan outmessage) <-chan outmessage {
	var wg sync.WaitGroup

	out := make(chan outmessage)

	output := func(c <-chan outmessage) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
