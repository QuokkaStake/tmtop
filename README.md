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

See how it looks like (on Sentinel chain, which is cosmos-sdk based chain, and Nomic, which uses Tendermint but not cosmos-sdk):
[![asciicast](https://asciinema.org/a/pnmH6j1MHGNdUY8y4eF2Ut21M.svg)](https://asciinema.org/a/pnmH6j1MHGNdUY8y4eF2Ut21M)

## How can I set it up?

Download the latest release from [the releases page](https://github.com/QuokkaStake/tmtop/releases/). After that, you should unzip it, and you are ready to go:

```sh
wget <the link from the releases page>
tar <downloaded file>
./tmtop --config <path to config>
```

Alternatively, install `golang` (>=1.21), clone the repo and build it:
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

To run it for a sovereign Cosmos-based chain, something like this should be enough:
```
./tmtop --rpc-host <RPC host address>
```

To run it for a sovereign chain that is not Cosmos-based (for example, Nomic), something like this should be enough
(this will limit the app possibilities, as in, it won't display validators monikers,
upgrades status etc.):
```
/tmtop  --rpc-host <RPC host address> --chain-type tendermint
```

To run it for a Cosmos-based consumer chains (like Stride or Neutron),
something like this should be enough:
```
/tmtop  --rpc-host <RPC host address> --provider-rpc-host <provider RPC host> --consumer-chain-id <consumer chain ID>
```

There are more parameters to tweak, for all the possible arguments, see `./tmtop --help`.


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

## Troubleshooting

If you have issues with the app, try pressing the D button to open the debug panel.
Most likely, the app cannot connect to one of the hosts it needs to connect.
If there's something unusual, feel free to report a bug on this repository.

Some common errors:

Q: The app displays nothing and is stuck.

![Telegram](https://raw.githubusercontent.com/QuokkaStake/tmtop/main/images/rpc-host-issues.jpg)

A: It's likely the app cannot connect to the RPC host to get data from.

Q: The app displays hashes instead of validators' monikers.

![Telegram](https://raw.githubusercontent.com/QuokkaStake/tmtop/main/images/provider-host-issues.jpg)

A: For sovereign chains, it's likely that the app cannot connect to the provider host,
or it's not specified (so the app thinks it's a sovereign chain and tries to get validators
out of the consumer chain). Additionally, this won't work with non-cosmos-sdk chains, like Nomic.

Q: The app displays some monikers as hashes on a consumer chain.

![Telegram](https://raw.githubusercontent.com/QuokkaStake/tmtop/main/images/wrong-chain-id.jpg)

A: Verify you have a correct `--consumer-chain-id` specified.

## Which networks this is guaranteed to work?

In theory, it should work on any Tendermint-based network with any node that has its RPC accessible.

## How can I contribute?

Bug reports and feature requests are always welcome! If you want to contribute, feel free to open issues or PRs.
