package config

import "freemasonry.cc/blockchain/core"

var ClientToml = `
chain-id = "` + core.ChainID + `"
# The keyring's backend, where the keys are stored (os|file|kwallet|pass|test|memory)
keyring-backend = "os"
# CLI output format (text|json)
output = "text"
# <host>:<port> to Tendermint RPC interface for this chain
node = "tcp://localhost:` + core.RpcPort + `"
# Transaction broadcasting mode (sync|async|block)
broadcast-mode = "sync"
`
