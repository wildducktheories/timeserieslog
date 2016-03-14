package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/wildducktheories/timeserieslog"
)

// element represents a sortable element.
type element struct {
	line    string
	ordinal int
}

// Less for elements compares the line. Otherwise
// identical lines are distinguished by their position in the
// input stream.
func (e *element) Less(o tsl.Element) bool {
	oe := o.(*element)
	if e.line == oe.line {
		return e.ordinal < oe.ordinal
	} else {
		return e.line < oe.line
	}
}

// statistics is a JSON encodable type which contains
// observable statistics for the sort operation.
type statistics struct {
	Read            int      `json:"read"`
	Duration        int64    `json:"duration"`
	DurationSeconds float64  `json:"durationSeconds"`
	Args            []string `json:"args"`
	Window          int      `json:"window"`
	SpillLimit      int      `json:"spillLimit"`
}

// process encapsulates the processing state of a tsl-sort process.
type process struct {
	input       *bufio.Reader        // the source of records to be sorted
	output      *os.File             // the sink for sorted records
	sorted      chan tsl.SortedRange // a channel used to accumulate sorted records
	spills      chan tsl.SortedRange // a channel used to accumulate "old" records that need to be spilled
	progressive bool                 // true if sorted output is to be written progressively
	window      int                  // the number of records between splits
	windowCount int                  // the total number of records since the last spill
	buffer      tsl.UnsortedRange    // a buffer for records as they are read
	keep        tsl.SortedRange      // the younger part of the last split
	split       *element             // record to be used for the next split
	prevsplit   *element             // the record used on the last split
	stats       statistics
}

// dump Writes a SortedRange to the process' Writer by iterating
// over the specified range's cursor.
func (p *process) dump(r tsl.SortedRange) {
	c := r.Open()
	for {
		e := c.Next()
		if e == nil {
			break
		}
		p.output.WriteString(e.(*element).line + "\n")
	}
}

// accumulator reads SortedRanges from the sorted channel and depending on
// p.progressive either accumulates them in a new UnsortedRange or dumps
// them directly to the output.
//
// Since the p.sorted channel is guaranteed to contain ranges that
// arrive in strictly sorted order the unsorted range will simply
// append to the sorted half of a mutableRange.
func (p *process) accumulator(final chan<- tsl.SortedRange) {
	acc := tsl.NewUnsortedRange()
	for r := range p.sorted {
		if p.progressive {
			p.dump(r)
		} else {
			acc.Add(tsl.AsSlice(r))
		}
	}
	final <- acc.Freeze()
}

// spill accumulator accumulates spilled record into an UnsortedRange
func (p *process) spillAccumulator(final chan<- tsl.SortedRange) {
	acc := tsl.NewUnsortedRange()
	for r := range p.spills {
		acc.Add(tsl.AsSlice(r))
	}
	final <- acc.Freeze()
}

// Snapshot takes a snapshot of the input accumulated in the buffer
// splits it into two parts - a part older than p.split that is to be merged
// with the keep range from a previous iteration and
// a younger part that is to be kept around for the next cycle.
//
// We check that we haven't got any really old parts that would cause
// the sort order. If we have, then we write these to a spills channel
//
// Note that most of the expensive work associated with processing the
// ranges produced by this method is done in the goroutines that
// service the channels. On the Partition operations may involve an
// O(n.log(n)) sort of the unsorted arm of the frozen partition but
// given the assumption that the stream is mostly sorted, this should
// be relatively rare.
func (p *process) snapshot(current *element) {

	write, hold := p.buffer.Freeze().Partition(p.split, tsl.LessOrder)

	if p.prevsplit != nil {
		// check that we don't have anything that sorts before the previous split.
		var older tsl.SortedRange
		older, newer := write.Partition(p.prevsplit, tsl.LessOrder)
		if older.Limit() > 0 {
			// we can't write this to the sorted channel, because it violate
			// the invariant about never writing anything into p.sorted that
			// is older than p.prevsplit
			p.window = 2 * p.window
			p.windowCount = 0
			p.spills <- older
			write = newer // only write the newer portion of the write partition
		}
	}

	p.sorted <- tsl.Merge(p.keep, write)

	p.buffer, p.keep, p.prevsplit, p.split = tsl.NewUnsortedRange(), hold, p.split, current
}

// run copies records from the input reader into an unsorted range and occasionally
// takes snapshots of this unsorted range and rights the sorted results into
// the sorted channel which either accumulates the records for writing out later
// or progressively writes the records to stdout.
func (p *process) run() (int, int) {
	started := time.Now()

	p.buffer = tsl.NewUnsortedRange()
	p.keep = tsl.EmptyRange

	p.sorted = make(chan tsl.SortedRange)
	p.spills = make(chan tsl.SortedRange)

	final := make(chan tsl.SortedRange)
	finalSpills := make(chan tsl.SortedRange)

	go p.accumulator(final)
	go p.spillAccumulator(finalSpills)

	for {
		p.stats.Read++
		p.windowCount++
		if line, err := p.input.ReadString('\n'); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
			os.Exit(1)
		} else {
			line = strings.TrimSpace(line)
			current := &element{line: line, ordinal: p.stats.Read}
			if p.stats.Read == 1 {
				p.split = current
			}
			p.buffer.Add([]tsl.Element{current})
			if p.windowCount%p.window == 0 {
				p.snapshot(current)
			}
		}
	}
	p.sorted <- tsl.Merge(p.keep, p.buffer.Freeze())

	close(p.sorted)
	close(p.spills)

	spill := <-finalSpills
	p.dump(tsl.Merge(<-final, spill))

	p.stats.Duration = int64(time.Now().Sub(started))
	p.stats.DurationSeconds = float64(p.stats.Duration) / float64(time.Second)
	p.stats.Args = os.Args[1:]
	p.stats.Window = p.window
	p.stats.SpillLimit = spill.Limit()

	return spill.Limit(), p.window
}

func main() {
	process := &process{
		input:  bufio.NewReader(os.Stdin),
		output: os.Stdout,
	}

	dumpStatistics := false
	comment := ""

	flag.IntVar(&process.window, "window", 1024, "The number of records to buffer prior to writing.")
	flag.BoolVar(&dumpStatistics, "statistics", false, "Dump a statistics record to stdout on exit.")
	flag.StringVar(&comment, "comment", "", "Arbitrary text to be logged as an argument.")
	flag.BoolVar(&process.progressive, "progressive", false, "Progressively write sorted output with finite probability that data will be written out of sort order.")
	flag.Parse()

	spillLimit, finalWindow := process.run()

	if dumpStatistics {
		json.NewEncoder(os.Stderr).Encode(process.stats)
	}

	if process.progressive && spillLimit > 0 {
		fmt.Fprintf(os.Stderr, "last %d records written out of order. increase --window to at least %d\n", spillLimit, finalWindow)
		os.Exit(1)
	}
}
