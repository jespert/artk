module artk.dev/x/eventlog

go 1.22.0

require (
	artk.dev v0.1.0
	artk.dev/x/testlog v0.1.0
)

require (
	github.com/lmittmann/tint v1.0.4 // indirect
	github.com/neilotoole/slogt v1.1.0 // indirect
)

replace artk.dev => ../../

replace artk.dev/x/testlog => ../testlog
