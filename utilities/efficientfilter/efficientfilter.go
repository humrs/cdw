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
	fsevent cdw.BinaryFlowsheetEvent
	number int
	observation string
}

// efficientfilter takes a csv file and sorts the adt events and outputs a csv file
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

	fp, err := os.Open(*infilename)
	if err != nil {
		log.Fatalf("Unable to open input filename.")
	}

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
		fp.Close()
		locfp.Close()
		itemfp.Close()
		fsfp.Close()
	}()

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

	fmt.Printf("Starting: %v\n", time.Now())
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

			msgch <- message{number: lines, line: line}

			lines++
		}

		fmt.Printf("Pass: %v  Total: %v  Lines: %v  LinesCorrected: %v\n", pass, total, lines, linesCorrected)

	}
	close(msgch)
	wg.Wait()

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

	fmt.Printf("Total Count = %v\n", total)
	close(indexc)
	close(obsc)

}

func fsoutfileservice(cs chan cdw.BinaryFlowsheetEvent) {
	outfp, err := os.Create("C:\\Users\\stanley\\Desktop\\fsindex.bin")
	if err != nil {
		log.Fatalf("Unable to open output filename.")
	}

	defer outfp.Close()

	for i := range cs {
		_, err = outfp.Write(i.ToBytes())
		if err != nil {
			log.Fatalf("fsout: %v\n", err)
		}
		_, err = outfp.WriteString("\n")
		if err != nil {
			log.Fatalf("fsout: %v\n", err)
		}
	}
}

func fsobsoutfileservice(obs chan string) {
	outfp, err := os.Create("C:\\Users\\stanley\\Desktop\\observations.txt")
	if err != nil {
		log.Fatalf("Unable to open observation file.")
	}

	defer outfp.Close()

	for i := range obs {
		outfp.WriteString(i)
		outfp.WriteString("\n")
	}
}

func processLine(msgc chan message, indexc chan cdw.BinaryFlowsheetEvent, obsc chan string, wg sync.WaitGroup) {
	defer wg.Done()

	for msg := range msgc {
		items := strings.Split(msg.line, "|")

		var binaryfsevent cdw.BinaryFlowsheetEvent
		binaryfsevent.FromLine(items[:9], msg.number)

		observationArray := []string{strconv.Itoa(msg.number)}
		observationArray = append(observationArray, items[9:len(items)]...)
		observationItem := strings.Join(observationArray, "|")

		indexc <- binaryfsevent
		obsc <- observationItem
	}

}
