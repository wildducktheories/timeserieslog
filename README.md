#NAME
timeserieslog - a go API optimized for processing timeseries data

#DESCRIPTION
timeserieslog is a go API that is optimized for processing timeseries data.

Timeseries data is data which tends to arrive in an almost, but not completely sorted order. For some applications, it is important to be able to access the data in a strictly sorted order. This API provides
abstractions which provide efficient means to sort such data in a continuous, streaming manner taking
advantage of the mostly sorted nature of most timeseries data.

#EXAMPLE

The [example](https://github.com/wildducktheories/timeserieslog/blob/master/examples/toy-tsl-sort/main.go) contains an implementation of a sort utility that makes use of the timeserieslog API to efficiently sort timeseries data.

toy-tsl-sort reads lines from stdin and writes them into a `tsl.UnsortedRange`. Every so often (based on the number of records configured with the --window parameter), the reading goroutine takes a snapshot of the records read since the last snapshot, and partitions them into two ranges, using a split record to identify the partition boundary. Records _younger_ than the `split` record are put aside in a local variable called `hold`. Records _older_ than the `split` record (but younger than the `prevsplit` record) are merged with records kept around from the last iteration (`keep`). The resulting merged range is written into the `sorted` channel. Any records _older_ than the `prevsplit` record are written into a spill channel.
After each snapshot is processed `keep`, `prevsplit` and `split` are updated based on the values of `hold`, `split` and the current record.

The `--progressive` option determines whether the sort proceeds optimistically or conservatively. If `--progressive` is true, the sort proceeds with the assumption that the configured window size is sufficient
to absorb any input records that arrive out of order and so it writes sorted data as soon as it falls out
of the current window. Conversely, if `--progressive` is false (the default), the sort processes all the input before generating any output.

`--progressive=true` favours performance over correctness; `--progressive=false` favours correctness over performance. In particular, if out of order data that exceeds the buffering window arrives, the progressive sort will write such data into a spill buffer which is tacked onto the end of the otherwise sorted stream and write a warning message to stderr and set a non-zero exit code. On the otherhand, the non-progressive sort will retain the data and perform a final merge sort with the bulk of the sorted input and then write the fully (and correctly) sorted result into the output stream.

[timestamp.txt.gz](https://github.com/wildducktheories/timeserieslog/blob/master/examples/toy-tsl-sort/timestamp.txt.gz) contains 4.2M records from a real webserver's access log. Each line consists of timestamp at which a request started. The data was written to the log at the time each request ended. This is example of timeseries data that is almost, but not actually sorted.

The following shows the sort performance of toy-tsl-sort compared to the standard OSX sort utility.

First, of all, we build the utility:

    $ (cd examples/toy-tsl-sort; go build)

Next, we baseline the performance of the OSX gnu sort utility:

    $ (cd examples/toy-tsl-sort/;
       gzip -dc timestamps.txt.gz | \
       time /usr/bin/sort | \
       wc)
    59.15 real        58.67 user         0.30 sys
     4284042 8568084 102817008

Next, we run ./toy-tsl-sort without any parameters.

    $ (cd examples/toy-tsl-sort/; gzip -dc timestamps.txt.gz | \
       time ./toy-tsl-sort | \
       wc)
    12.60 real        14.59 user         3.79 sys
    4284042 8568084 102817008

This is ~ 5x faster than the OSX sort utility. The apparent reason for the discrepancy is because the OSX sort is writing work files into the /tmp directory.

toy-tsl-sort supports a --progressive output option which optimistically assumes that the
default window size is sufficient to cope with any unsorted data in the input stream.
If the assumption is too optimistic the sort will fail with a non-zero exit code and
print a message indicating how many records were written out of order.

    $ (cd examples/toy-tsl-sort/;
       gzip -dc timestamps.txt.gz | \
       time ./toy-tsl-sort --progressive | \
       tail -20)
    last 11 records written out of order. increase --window to at least 8192
            8.41 real        11.56 user         2.64 sys
    2016-03-12 10:28:51.464
    2016-03-12 10:28:54.288
    2016-03-12 10:28:58.480
    2016-03-12 10:29:01.354
    2016-03-12 10:29:04.482
    2016-03-12 10:29:08.477
    2016-03-12 10:29:11.341
    2016-03-12 10:29:14.447
    2016-03-12 10:29:18.480
    2016-03-06 22:25:53.499
    2016-03-06 22:25:54.465
    2016-03-06 22:25:54.499
    2016-03-06 22:26:41.499
    2016-03-06 22:26:41.499
    2016-03-06 22:26:43.499
    2016-03-06 22:26:44.093
    2016-03-06 22:27:31.329
    2016-03-06 22:27:34.499
    2016-03-06 22:27:34.499
    2016-03-06 22:27:34.499

A `--window` option may be used to increase the default window size.

    $ (cd examples/toy-tsl-sort/;
       gzip -dc timestamps.txt.gz | \
       time ./toy-tsl-sort --progressive --window=8192 | \
       wc)
    7.34 real        12.84 user         3.34 sys
    4284042 8568084 102817008

Notice that the elapsed time of the progressive sort is about 5 seconds faster than the non-progressive sort. The reason is that the optimitistic progressive sort can write output as it goes whereas the conservative non-progressive sort must sort all the data before writing any of it and so there is no possibilty to
take advantage of available concurrency between the CPU and IO paths.
