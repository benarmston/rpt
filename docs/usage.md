# NAME

rpt - run the given command the given number of times

# SYNOPSIS

```
rpt [OPTIONS] TIMES COMMAND [ARGUMENTS...]
```

# DESCRIPTION

Run `COMMAND ARGUMENTS` TIMES times.

# OPTIONS

**--delay=DURATION, -d=DURATION**: wait `DURATION` between runs (default: 0s)

**--leading-edge**: if given, any provided delay is between the
start of one command invocation and the start the next. If not given,
any provided delay is between the end of one invocation and the start of
the next

**--fail-fast**: if COMMAND fails, exit immediately with the same exit code
as COMMAND

**--verbose, -v**: print debugging messages



# AUTHOR

Ben Armston

# COPYRIGHT

Copyright 2025 Ben Armston. Licensed under the MIT License.
