## DaoDst BlankChain

Golang execution layer implementation of the NXNOS protocol.

[![CLI Reference](
https://pkg.go.dev/badge/github.com/ethereum/go-ethereum
)](https://docs.daodst.com/protocol/stcd/)
[![Modules](https://goreportcard.com/badge/github.com/ethereum/go-ethereum)](https://docs.daodst.com/modules_accounts/)
[![Chat](https://app.travis-ci.com/ethereum/go-ethereum.svg?branch=master)](https://docs.daodst.com/chat/installation/)

[![FAQ](https://img.shields.io/badge/discord-join%20chat-blue.svg)](https://docs.daodst.com/faq/)

Automated builds are available for stable releases and the unstable master branch. Binary
archives are published at https://www.daodst.com/#/download.

## Building the source

For prerequisites and detailed build instructions please read the [Installation Instructions](https://docs.daodst.com/protocol/stcd/).

Building `stcd` requires both a Go (version 1.20 or later) and a C compiler. You can install
them using your favourite package manager. Once the dependencies are installed, run

```shell
make stcd --version
```

or, to build the full suite of utilities:

```shell
make all
```

## Executables

The daodst project comes with several wrappers/executables found in the `cmd`
directory.

|  Command   | Description                                                                                                                                                      |
|:----------:|------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **`stcd`** | Blockchain node program, completes blockchain synchronization processing, and all related functions                                                              |
|   `chat`   | Provide communication layer processing projects                                                                                                                  |
| `gateway`  | The mobile application connects to the gateway endpoint, provides block query services, and connects to the communication relay                                  |
|   `dvm`    | DVM (DST VM) is a dedicated module used to run smart contracts, supports virtual machines running on the EVM platform, and adds support for high-level languages. |
|  `NXNOS`   | Device Interconnect Protocol Operating System.|

## Running `stcd`

Going through all the possible command line flags is out of scope here (please consult our
[CLI Wiki page](https://docs.daodst.com/protocol/stcd/command/)),
but we've enumerated a few common parameter combos to get you up to speed quickly
on how you can run your own `stcd` instance.

### Hardware Requirements

Minimum:

* CPU with 8+ cores
* 64GB RAM
* 1TB free storage space to sync the Mainnet
* 8 MBit/sec download Internet service
* Port Open 8545 8546


### Full node on the main DaoDst network

By far the most common scenario is people wanting to simply interact with the DaoDst
network: create accounts; transfer funds; deploy and interact with contracts. For this
particular use case, the user doesn't care about years-old historical data, so we can
sync quickly to the current state of the network. To do so:

```shell
$ stcd console
```

This command will:
* Start `stcd` in snap sync mode (default, can be changed with the `--trust-node` flag),
causing it to download more data in exchange for avoiding processing the entire history
of the DaoDst network, which is very CPU intensive.


### A Full node on the  test network

Transitioning towards developers, if you'd like to play around with creating DaoDst
contracts, you almost certainly would like to do that without any real money involved until
you get the hang of the entire system. In other words, instead of attaching to the main
network, you want to join the **test** network with your node, which is fully equivalent to
the main network, but with play-DaoDst only.

```shell
$ stcd start --chain-id=testnet
```

The `console` subcommand has the same meaning as above and is equally
useful on the testnet too.

Specifying the `--chain-id=testnet` flag, however, will reconfigure your `stcd` instance a bit:

* Instead of connecting to the main DaoDst network, the client will connect to the 
test network, which uses different P2P bootnodes, different network IDs and genesis
states.

*Note: Although some internal protective measures prevent transactions from
crossing over between the main network and test network, you should always
use separate accounts for play and real money. Unless you manually move
accounts, `stcd` will by default correctly separate the two networks and will not make any
accounts available between them.*

### Configuration

As an alternative to passing the numerous flags to the `stcd` binary, you can also pass a
configuration file via:

```shell
$ stcd start --home /path/to/your_config
```

To get an idea of how the file should look like you can use the `ll /path/to/your_config/.stcd/config/` subcommand to
export your existing configuration files:

*Note: This works only with `stcd` v23.1.3 and above.*

### Programmatically interfacing `stcd` nodes

As a developer, sooner rather than later you'll want to start interacting with `stcd` and the
DaoDst network via your own programs and not manually through the console. To aid
this, `geth` has built-in support for a JSON-RPC based APIs ([standard APIs](https://ethereum.github.io/execution-apis/api-documentation/).
These can be exposed via HTTP, WebSockets and IPC (UNIX sockets on UNIX based
platforms, and named pipes on Windows).

The IPC interface is enabled by default and exposes all the APIs supported by `stcd`,
whereas the HTTP and WS interfaces need to manually be enabled and only expose a
subset of APIs due to security reasons. These can be turned on/off and configured as
you'd expect.

HTTP based JSON-RPC API options:

* `--http` Enable the HTTP-RPC server
* `--http.addr` HTTP-RPC server listening interface (default: `localhost`)
* `--http.port` HTTP-RPC server listening port (default: `8545`)
* `--http.api` API's offered over the HTTP-RPC interface (default: `eth,net,web3`)
* `--http.corsdomain` Comma separated list of domains from which to accept cross origin requests (browser enforced)
* `--ws` Enable the WS-RPC server
* `--ws.addr` WS-RPC server listening interface (default: `localhost`)
* `--ws.port` WS-RPC server listening port (default: `8546`)
* `--ws.api` API's offered over the WS-RPC interface (default: `eth,net,web3`)
* `--ws.origins` Origins from which to accept WebSocket requests
* `--ipcdisable` Disable the IPC-RPC server
* `--ipcapi` API's offered over the IPC-RPC interface (default: `admin,debug,eth,miner,net,personal,txpool,web3`)
* `--ipcpath` Filename for IPC socket/pipe within the datadir (explicit paths escape it)

You'll need to use your own programming environments' capabilities (libraries, tools, etc) to
connect via HTTP, WS or IPC to a `stcd` node configured with the above flags and you'll
need to speak [JSON-RPC](https://www.jsonrpc.org/specification) on all transports. You
can reuse the same connection for multiple requests!

**Note: Please understand the security implications of opening up an HTTP/WS based
transport before doing so! Hackers on the internet are actively trying to subvert
DaoDst nodes with exposed APIs! Further, all browser tabs can access locally
running web servers, so malicious web pages could try to subvert locally available
APIs!**

### Operating a private network

Maintaining your own private network is more involved as a lot of configurations taken for
granted in the official networks need to be manually set up.

#### Defining the private genesis state

First, you'll need to create the genesis state of your networks, which all nodes need to be
aware of and agree upon. This consists of a small JSON file (e.g. call it `genesis.json`):

```json
{
  "genesis_time": "2024-02-28T01:10:22.3276221Z",
  "chain_id": "daodst_9000-1",
  "initial_height": "1",
  "consensus_params": {
    "block": {
      "max_bytes": "22020096",
      "max_gas": "-1",
      "time_iota_ms": "1000"
    },
    "evidence": {
      "max_age_num_blocks": "100000",
      "max_age_duration": "172800000000000",
      "max_bytes": "1048576"
    },
    "validator": {
      "pub_key_types": [
        "ed25519"
      ]
    },
    "version": {}
  },
  "app_hash": "",
  "app_state": {
    "auth": {
      "params": {
        "max_memo_characters": "256",
        "tx_sig_limit": "7",
        "tx_size_cost_per_byte": "10",
        "sig_verify_cost_ed25519": "590",
        "sig_verify_cost_secp256k1": "1000"
      }
    },
    "authz": {
      "authorization": []
    },
    "bank": {
      "params": {
        "send_enabled": [],
        "default_send_enabled": true
      },
      "denom_metadata": []
    },
    "capability": {
      "index": "1",
      "owners": []
    },
    "chat": {
      "params": {
        "max_phone_number": "10",
        "destroy_phone_number_coin": {
          "denom": "dst",
          "amount": "10000000000000000000"
        },
        "chat_fee": {
          "denom": "dst",
          "amount": "1000000000000000000"
        },
        "min_register_burn_amount": {
          "denom": "dst",
          "amount": "100000000000000000000"
        },
        "min_burn_amount": {
          "denom": "dst",
          "amount": "100000000000000"
        }
      }
    },
    "claims": {
      "params": {
        "enable_claims": true,
        "airdrop_start_time": "0001-01-01T00:00:00Z",
        "duration_until_decay": "2629800s",
        "duration_of_decay": "5259600s",
        "claims_denom": "aevmos",
        "authorized_channels": [
          "channel-0",
          "channel-3"
        ],
        "evm_channels": [
          "channel-2"
        ]
      },
      "claims_records": []
    },
    "contract": {
      "params": {}
    },
    "crisis": {
      "constant_fee": {
        "denom": "nxn",
        "amount": "1000"
      }
    },
    "dao": {
      "params": {
        "rate": "0.100000000000000000",
        "min_salary_reward_ratio": "0.000000000000000000",
        "max_salary_reward_ratio": "0.999000000000000000",
        "burn_get_power_ratio": "100.000000000000000000",
        "max_cluster_members": "666",
        "min_create_cluster_pledge_amount": "5000",
        "dao_reward_percent": "0.100000000000000000",
        "dpos_reward_percent": "0.100000000000000000",
        "burn_current_gate_ratio": "0.057500000000000000",
        "burn_register_gate_ratio": "0.037500000000000000",
        "day_mint_amount": "360000000000000000000000.000000000000000000",
        "power_gas_ratio": "100.000000000000000000",
        "ad_price": "0.050000000000000000",
        "ad_rate": "0.200000000000000000",
        "burn_reward_fee_rate": "0.300000000000000000",
        "receive_dao_ratio": "0.075000000000000000",
        "connectivity_dao_ratio": "0.700000000000000000",
        "burn_dao_pool": "0.100000000000000000",
        "min_dao_reward_ratio": "0.000000000000000000",
        "max_dao_reward_ratio": "0.900000000000000000",
        "max_online_ratio": "0.700000000000000000"
      }
    },
    "distribution": {
      "params": {
        "community_tax": "0.010000000000000000",
        "base_proposer_reward": "0.010000000000000000",
        "bonus_proposer_reward": "0.040000000000000000",
        "withdraw_addr_enabled": true,
        "fee_burn_ratio": "0.500000000000000000"
      },
      "fee_pool": {
        "community_pool": []
      },
      "delegator_withdraw_infos": [],
      "previous_proposer": "",
      "outstanding_rewards": [],
      "validator_accumulated_commissions": [],
      "validator_historical_rewards": [],
      "validator_current_rewards": [],
      "delegator_starting_infos": [],
      "validator_slash_events": []
    },
    "erc20": {
      "params": {
        "enable_erc20": true,
        "enable_evm_hook": true
      },
      "token_pairs": []
    },
    "evidence": {
      "evidence": []
    }
    "feegrant": {
      "allowances": []
    },
    "feemarket": {
      "params": {
        "no_base_fee": false,
        "base_fee_change_denominator": 8,
        "elasticity_multiplier": 2,
        "enable_height": "0",
        "base_fee": "1000000000",
        "min_gas_price": "0.000000000000000000",
        "min_gas_multiplier": "0.500000000000000000"
      },
      "block_gas": "0"
    },
    "gateway": {
      "params": {
        "index_num_height": "100",
        "redeem_fee_height": "432000",
        "redeem_fee": "0.100000000000000000",
        "min_delegate": {
          "denom": "nxn",
          "amount": "10000000000000000000"
        },
        "validity": "5256000"
      }
    },
 
    "gov": {
      "starting_proposal_id": "1",
      "deposits": [],
      "votes": [],
      "proposals": [],
      "deposit_params": {
        "min_deposit": [
          {
            "denom": "nxn",
            "amount": "1000000000000000000"
          }
        ],
        "max_deposit_period": "300s"
      },
      "voting_params": {
        "voting_period": "300s"
      },
      "tally_params": {
        "quorum": "0.334000000000000000",
        "threshold": "0.500000000000000000",
        "veto_threshold": "0.334000000000000000"
      }
    },
    "group": {
      "group_seq": "0",
      "groups": [],
      "group_members": [],
      "group_policy_seq": "0",
      "group_policies": [],
      "proposal_seq": "0",
      "proposals": [],
      "votes": []
    },
    "ibc": {
      "client_genesis": {
        "clients": [],
        "clients_consensus": [],
        "clients_metadata": [],
        "params": {
          "allowed_clients": [
            "06-solomachine",
            "07-tendermint"
          ]
        },
        "create_localhost": false,
        "next_client_sequence": "0"
      },
      "connection_genesis": {
        "connections": [],
        "client_connection_paths": [],
        "next_connection_sequence": "0",
        "params": {
          "max_expected_time_per_block": "30000000000"
        }
      },
      "channel_genesis": {
        "channels": [],
        "acknowledgements": [],
        "commitments": [],
        "receipts": [],
        "send_sequences": [],
        "recv_sequences": [],
        "ack_sequences": [],
        "next_channel_sequence": "0"
      }
    },
    "mint": {
      "minter": {
        "inflation": "0.100000000000000000",
        "annual_provisions": "0.000000000000000000"
      },
      "params": {
        "mint_denom": "nxn",
        "inflation_rate_change": "1.000000000000000000",
        "inflation_max": "0.070000000000000000",
        "inflation_min": "0.030000000000000000",
        "goal_bonded": "0.670000000000000000",
        "blocks_per_year": "5259600"
      }
    },
    "params": null,
    "recovery": {
      "params": {
        "enable_recovery": true,
        "packet_timeout_duration": "14400s"
      }
    },
    "slashing": {
      "params": {
        "signed_blocks_window": "100",
        "min_signed_per_window": "0.500000000000000000",
        "downtime_jail_duration": "600s",
        "slash_fraction_double_sign": "0.050000000000000000",
        "slash_fraction_downtime": "0.010000000000000000"
      },
      "signing_infos": [],
      "missed_blocks": []
    },
    "staking": {
      "params": {
        "unbonding_time": "1814400s",
        "max_validators": 100,
        "max_entries": 7,
        "historical_entries": 10000,
        "bond_denom": "nxn",
        "min_commission_rate": "0.050000000000000000"
      },
      "last_total_power": "0",
      "last_validator_powers": [],
      "validators": [],
      "delegations": [],
      "unbonding_delegations": [],
      "redelegations": [],
      "exported": false
    },
    "transfer": {
      "port_id": "transfer",
      "denom_traces": [],
      "params": {
        "send_enabled": true,
        "receive_enabled": true
      }
    },
    "upgrade": {},
    "vesting": {}
  }
}
```

The above fields should be fine for most purposes, although we'd recommend changing
the `nonce` to some random value so you prevent unknown remote nodes from being able
to connect to you. If you'd like to pre-fund some accounts for easier testing, create
the accounts and populate the `alloc` field with their addresses.

```json
  "bank": {
    "params": {
      "send_enabled": [],
      "default_send_enabled": true
    },
    "balances": [
      {
        "address": "your_dst_address",
        "coins": [
          {
            "denom": "dst",
            "amount": "11111111"
          },
        ]
      }
    ]
  }
```

With the genesis state defined in the above JSON file, you'll need to initialize **every**
`stcd` node with it prior to starting it up to ensure all blockchain parameters are correctly
set:

```shell
$ init.sh
```

#### Creating the rendezvous point

With all nodes that you want to run initialized to the desired genesis state, you'll need to
start a bootstrap node that others can use to find each other in your network and/or over
the internet. The clean way is to configure and run a dedicated bootnode:

```shell
$ bootnode --genkey=boot.key
$ bootnode --nodekey=boot.key
```

With the bootnode online, it will display an
that other nodes can use to connect to it and exchange peer information. Make sure to
replace the displayed IP address information (most probably `[::]`) with your externally
accessible IP to get the actual `enode` URL.

*Note: You could also use a full-fledged `stcd` node as a bootnode, but it's the less
recommended way.*

#### Starting up your member nodes

With the bootnode operational and externally reachable (you can try
`telnet <ip> <port>` to ensure it's indeed reachable), start every subsequent `stcd`
node pointed to the bootnode for peer discovery via the `--p2p.seeds` flag. It will
probably also be desirable to keep the data directory of your private network separated, so
do also specify a custom `--home` flag.

```shell
$ stcd --home=path/to/custom/data/folder --p2p.seeds=<bootnode-enode-url-from-above>
```

*Note: Since your network will be completely cut off from the main and test networks, you'll
also need to configure a miner to process transactions and create new blocks for you.*

## Contribution

Thank you for considering helping out with the source code! We welcome contributions
from anyone on the internet, and are grateful for even the smallest of fixes!

If you'd like to contribute to Daodst, please fork, fix, commit and send a pull request
for the maintainers to review and merge into the main code base. 
to ensure those changes are in line with the general philosophy of the project and/or get
some early feedback which can make both your efforts much lighter as well as our review
and merge procedures quick and simple.

Please make sure your contributions adhere to our coding guidelines:

* Code must adhere to the official Go [formatting](https://golang.org/doc/effective_go.html#formatting)
guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt/)).
* Code must be documented adhering to the official Go [commentary](https://golang.org/doc/effective_go.html#commentary)
guidelines.
* Pull requests need to be based on and opened against the `master` branch.
* Commit messages should be prefixed with the package(s) they modify.
* E.g. "daodst, rpc: make trace configs optional"



## License

The stcd library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html),
also included in our repository in the `COPYING.LESSER` file.

The go-ethereum binaries (i.e. all code inside of the `cmd` directory) are licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also
included in our repository in the `COPYING` file.
