# Jonathan's OnLine / Streaming Algorithm Tool: jolsat

`jolsat` is a command line utility modelled after the venerable `cut` tool, but with the added ability to display summary statistics about its input. It uses "online" and "streaming" algorithms to calculate these statistics, which allows it to print them as it processes its input, without having to hold the entire dataset in memory. Because of this choice, it is able to process input streams indefinitely, without increasing its memory use in an unbounded fashion. This is useful when tailing live log files, for example.

## Features

- Produces [exactly the same output](#compatibility-with-cut) as `cut`, given only `-d` and `-f` parameters
- Field transformation before statistical processing
- Timestamp insertion based on record arrival time
- Automatic and tuneable training mode produces full-fidelity statistics for smaller datasets

### Compatibility with `cut`

It is a bug for `jolsat`:

- given **any** input
- with **any** invocation containing only `-d` and `-f` parameters, properly formed
- to produce different output to `cut`, called with the same parameters and input.

NB all short-form parameters must have a space seperating them from their values. i.e. `jolsat -d ' ' -f 1-10` as opposed to `jolsat -d' ' -f1-10`. Both forms are acceptable to `cut`; compatibility is only assured with the space-containing form.

There is one exception to this: where the `-f` parameter uses `jolsat`'s enhanced ability to repeat or re-order output fields. For `cut` compatibility the `-f` parameter must only reference increasing field numbers.

This also explicitly rules out the cases where `cut` is invoked with no parameters, or`cut -d <CHAR>` is invoked *without* a `-f` parameter. Both of these `cut` invocations produce an error, but `jolsat` treats both as a request to behave exactly like `cut -f 1-`: to mirror stdin to stdout, without erroring. This allows for easier incremental exploration of an input stream, where the field list input param may be built up over time, and is a further byte-for-byte compatibility guarantee:

It is a bug for `jolsat`:

- with **any** input
- when invoked with no parameters
- or invoked with only a properly formed `-d <CHAR>` parameter
- to produce different output to `cut -f1-`, given the same input steam.

## Usage

### Basic arguments

```
$ seq 1 20 | xargs -n5 | jolsat
1 2 3 4 5
6 7 8 9 10
11 12 13 14 15
16 17 18 19 20
$ seq 1 20 | xargs -n5 | jolsat -d ' '
1 2 3 4 5
6 7 8 9 10
11 12 13 14 15
16 17 18 19 20
```

### Specifying fields and field ranges, compatible with `cut`

```
$ seq 1 20 | xargs -n5 | jolsat -d ' ' -f 2,4-5
2 4 5
7 9 10
12 14 15
17 19 20
```
