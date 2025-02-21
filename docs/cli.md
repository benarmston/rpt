# CLI

# NAME

rpt - repeat running a command a number of times

# SYNOPSIS

rpt

```
[--delay|-d]=[value]
[--times|-t]=[value]
[--verbose|-v]
```

# DESCRIPTION

Repeatedly run COMMAND with ARGUMENTS.  The number of times to run COMMAND
is determined by OPTIONS.

**Usage**:

```
rpt [OPTIONS] COMMAND [ARGUMENTS...]
```

# OPTIONS

**--delay, -d**="": wait `DURATION` between runs (default: 0s)

**--times, -t**="": number of `TIMES` to run COMMAND (default: 1)

**--verbose, -v**: print debugging messages

