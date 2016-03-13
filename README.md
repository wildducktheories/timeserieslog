#NAME
timeserieslog - a go API optimized for processing timeseries data

#DESCRIPTION
timeserieslog is a go API that is optimized for processing timeseries data.

Timeseries data is data which tends to arrive in an almost, but not completely sorted order. For some applications, it is important to be able to access the data in a strictly sorted order. This API provides
abstractions which provide efficient means to sort such data in a continuous, streaming manner taking
advantage of the mostly sorted nature of most timeseries data.

#USAGE

The following example, found in examples/toy-tsl-sort, shows how to use the API to strictly sort
almost sorted timeseries data.

The basic idea is that data is buffered into a `tsl.UnsortedRange` until a configurable number of
records have been read (controlled by the `--window` option). That range is then frozen and the
data is then partitioned into two sorted ranges - the records before a split record (`write`) and the records afterwards (`hold`). The records before the split record are merged with any records kept from a previous iteration (`keep`) and then immediately written to the output. The records after the split record are retained (`keep`) and the split record is updated to the current record. The process continues until all the records have been read.

The process detects if the configured window is too small to successfully sort the data and will
will die with a non-zero exit status if this is the case.

	package main

	import (
		"bufio"
		"flag"
		"fmt"
		"io"
		"os"
		"strings"
		"sync"

		"github.com/wildducktheories/timeserieslog"
	)

	type element struct {
		line    string
		ordinal int
	}

	func (e *element) Less(o tsl.Element) bool {
		oe := o.(*element)
		if e.line == oe.line {
			return e.ordinal < oe.ordinal
		} else {
			return e.line < oe.line
		}
	}

	func main() {
		window := 0
		flag.IntVar(&window, "window", 4000, "The number of records to buffer prior to writing.")
		flag.Parse()

		r := bufio.NewReader(os.Stdin)
		count := 0
		buffer := tsl.NewUnsortedRange()
		var keep tsl.SortedRange = tsl.EmptyRange
		var split string
		var prevsplit string

		wg := sync.WaitGroup{}
		ch := make(chan tsl.SortedRange)

		wg.Add(1)
		go func() {
			defer wg.Done()

			for r := range ch {
				c := r.Open()
				for {
					e := c.Next()
					if e == nil {
						break
					}
					os.Stdout.WriteString(e.(*element).line + "\n")
				}
			}
		}()

		for {
			count++
			if l, err := r.ReadString('\n'); err != nil {
				if err == io.EOF {
					break
				}
				fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
				os.Exit(1)
			} else {
				l = strings.TrimSpace(l)
				if count == 1 {
					split = l
				}
				buffer.Add([]tsl.Element{&element{
					line:    l,
					ordinal: count,
				}})
				if count%window == 0 {
					write, hold := buffer.Freeze().Partition(&element{
						line:    split,
						ordinal: 0,
					}, tsl.LessOrder)

					if prevsplit != "" {

						check, _ := write.Partition(&element{
							line:    prevsplit,
							ordinal: 0,
						}, tsl.LessOrder)

						if check.Limit() > 0 {

							// if check has any records, then it means the window is too small to succesfully sort
							// the data.

							fmt.Fprintf(os.Stderr, "window too small: %s: %s", tsl.AsSlice(check)[0].(*element).line, prevsplit)
							os.Exit(1)
						}
					}
					ch <- tsl.Merge(keep, write)

					keep = hold
					prevsplit = split
					split = l
					buffer = tsl.NewUnsortedRange()
				}
			}
		}

		ch <- tsl.Merge(keep, buffer.Freeze())
		close(ch)
		wg.Wait()
	}

examples/toy-tsl-sort/timestamp.txt.gz contains 4.2M records from a real webserver's access log. Each line consists of timestamp at which a request started. The data is written at the time each request ends. This
is example of timeseries data that is almost, but not actually sorted. The following shows the
sort performance of toy-tsl-sort compared to the standard OSX sort utility.

	$ (cd examples/toy-tsl-sort; go build)
	$ (cd examples/toy-tsl-sort/; gzip -dc timestamps.txt.gz | time ./toy-tsl-sort --window=4000 | wc)
	        7.35 real        10.06 user         3.43 sys
	 4284042 8568084 102817008

	$ (cd examples/toy-tsl-sort/; gzip -dc timestamps.txt.gz | time sort | wc)
	       59.15 real        58.67 user         0.30 sys
	 4284042 8568084 102817008


