# Jonathan's OnLine / Streaming Algorithm Tool: jolsat

`jolsat` is a command line utility modelled after the venerable `cut` tool, but with the added ability to display summary statistics about its input. It uses "online" and "streaming" algorithms to calculate these statistics, which allows it to print them as it processes its input, without having to hold the entire dataset in memory. Because of this choice, it is able to process input streams indefinitely, without increasing its memory use in an unbounded fashion. This is useful when tailing live log files, for example.

## Features

- Produces [exactly the same output](#compatibility-with-cut) as `cut`, given only `-d` and `-f` parameters
- Field transformation before statistical processing
- Timestamp insertion based on record arrival time
- Automatic and tuneable training mode produces full-fidelity statistics for smaller datasets

### Compatibility with `cut`

It is a bug for `jolsat` given **any** input and any invocation containing only `-d` and `-f` to produce different output from `cut` with the same parameters.

There is one exception to this: where the `-f` parameter uses `jolsat`'s enhanced ability to repeat or re-order output fields. For `cut` compatibility the `-f` parameter must only reference increasing field numbers.

This also explicitly rules out the cases where `cut` is invoked with no parameters, or`cut -d <CHAR>` is invoked *without* a `-f` parameter. Both of these `cut` invocations produce an error, but `jolsat` treats both as a request to behave exactly like `cat`: to mirror stdin to stdout, without erroring. This allows for easier incremental exploration of an input stream, where the field list input param may be built up over time, and is a further byte-for-byte compatibility guarantee: it is a bug for `jolsat` or any `jolsat -d <CHAR>` invocations to produce different output than `cat` does, given the same input steam.

## Usage

```
$ seq 1 20 | xargs -n5 | jolsat
1 2 3 4 5
6 7 8 9 10
11 12 13 14 15
16 17 18 19 20
```
