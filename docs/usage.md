# NAME

rpt - repeat running a command a number of times

# SYNOPSIS

```
rpt [OPTIONS] COMMAND [ARGUMENTS...]
```

# DESCRIPTION

Repeatedly run COMMAND with ARGUMENTS.  The number of times to run COMMAND
is determined by OPTIONS.

# OPTIONS

**-d=DURATION, --delay=DURATION**
: wait `DURATION` between runs (default: 0s)

**-t=TIMES, --times=TIMES**
: number of `TIMES` to run COMMAND (default: 1)

**-v, --verbose**
: print debugging messages



# AUTHOR

Ben Armston

# COPYRIGHT

Copyright 2025 Ben Armston. Licensed under the MIT License.
