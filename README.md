# Repeat command line invocations

## Overview

`rpt` (pronounced repeat) is a CLI tool for repeating command line invocations.

`rpt` exists to make repeating a command in a Bash shell more convenient. Let's
say we want to print "The task" 10 times with a one second delay, to do so we
could write the following in Bash.

```bash
for i in $(seq 1 10) ; do echo "The task"; sleep 1 ; done
```

With `rpt` we can instead write

```bash
rpt --delay 1s 10 echo "The task"
```

## Installation

Download the `.deb`, `.rpm` or `.apk` packages from the the [releases
page](https://github.com/benarmston/rpt/releases) and install them with the
appropriate tools.

## Usage

See [docs/usage.md](docs/usage.md) for usage.

## Contributing

If you found a bug or have a feature request, [create a new
issue](https://github.com/benarmston/rpt/issues/new).

## Copyright and License

Copyright (C) 2025 Ben Armston.  Licensed under the MIT License, see
[LICENSE](LICENSE) for details.
