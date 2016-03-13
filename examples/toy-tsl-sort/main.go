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
