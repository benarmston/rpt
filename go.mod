module github.com/benarmston/rpt

go 1.25.0

require github.com/urfave/cli/v3 v3.8.0

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/urfave/cli-docs/v3 v3.1.0
)

require (
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sergi/go-diff v1.4.0
)

replace github.com/urfave/cli-docs/v3 => github.com/benarmston/cli-docs/v3 v3.0.0-20260329190107-f9d064b7a39c
