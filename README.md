# tmtop

![Latest release](https://img.shields.io/github/v/release/QuokkaStake/tmtop)
[![Actions Status](https://github.com/QuokkaStake/tmtop/workflows/test/badge.svg)](https://github.com/QuokkaStake/tmtop/actions)

tmtop is a tool that can visualize the consensus on any Tendermint-based chain in a top-like way.
It is heavily inspired by [pvtop](https://github.com/blockpane/pvtop) and is in fact a from-scratch rewrite of it.

It can do the following:
- display the consensus state of the current block (the percentage, who prevoted/precommitted and who didn't etc.)
- display chain info (chain-id, block time, Tendermint version etc.)
- display chain upgrade info and estimated time
- work with non cosmos-sdk chains (for instance, Nomic; it won't be able to display the validators' monikers then)
- work with ICS (fetching the validators list from the provider chain while taking the consensus from the consumer chain)

## How can I set it up?

Download the latest release from [the releases page](https://github.com/QuokkaStake/tmtop/releases/). After that, you should unzip it, and you are ready to go:

```sh
wget <the link from the releases page>
tar <downloaded file>
./tmtop --config <path to config>
```

Alternatively, install `golang` (>1.18), clone the repo and build it:
```
git clone https://github.com/QuokkaStake/tmtop
cd tmtop
# This will generate a `tmtop` binary file in the repository folder
make build
# This will generate a `tmtop` binary file in $GOPATH/bin
```

Then just run it:

```
./tmtop [args]
```

For possible arguments, see `./tmtop --help`.

## How does it work?

It queries Tendermint's RPC node to get the following data:
- consensus state
- validators list and their voting power
- blocks and their time difference
and uses this data to build a consensus state to visualise.

Additionally, if it's a cosmos-sdk chain, it can also fetch the following data via the abci_query query:
- chain upgrade info
- validators list (to show validators' monikers instead of addresses)

## How can I configure it?

All configuration is done via flags, see `./tmtop --help` for the list of flags.
Additionally, the app itself has a few shortcuts allowing you to control it.
You can press the [h] button to display the help message, which will show you the shortcuts and when/how to use them.


## Which networks this is guaranteed to work?

In theory, it should work on any Tendermint-based network with any node that has its RPC accessible.

## How can I contribute?

Bug reports and feature requests are always welcome! If you want to contribute, feel free to open issues or PRs.
