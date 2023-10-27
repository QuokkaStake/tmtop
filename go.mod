module main

go 1.21.3

require (
	github.com/gdamore/tcell/v2 v2.6.0
	github.com/rivo/tview v0.0.0-20231022175332-f7f32ad28104
	github.com/rs/zerolog v1.31.0
	github.com/spf13/cobra v1.7.0
)

require (
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/rivo/uniseg v0.4.3 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/term v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/tendermint/tendermint => github.com/cometbft/cometbft v0.34.29
)
