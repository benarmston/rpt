# NAME

rpt - run the given command the given number of times

# SYNOPSIS

```
rpt [OPTIONS] TIMES COMMAND [-- ARGUMENTS...]
```

# DESCRIPTION

Run `COMMAND ARGUMENTS` TIMES times.

If the '--delay' option is given, there will be a delay of the given
DURATION after one run ends and the next starts. This provides a
guaranteed delay between runs.

If the '--every' option is given, COMMAND will be run every DURATION. If
COMMAND takes longer to run than the given DURATION the next run will
start immediately once the current run has completed. This provides a
predictable start time for each run (provided COMMAND consistently
completes in under DURATION).

# OPTIONS

**--fail-fast**: if COMMAND fails, exit immediately with the same exit code
as COMMAND

**--verbose, -v**: print debugging messages

**--delay=DURATION, -d=DURATION**: wait `DURATION` between one run ending and
the next starting (default: 0s)

**--every=DURATION, -e=DURATION**: run COMMAND every `DURATION`. The next run
will start DURATION after the previous run started or as soon as the
previous run ends if it takes longer than DURATION (default: 0s)



# AUTHOR

Ben Armston

# COPYRIGHT

Copyright 2025 Ben Armston. Licensed under the MIT License.
