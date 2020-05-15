# whissage (whisper messenger)
The backend of blockchain-based messenger built on ethereum whisper.

**Сurrently in active development**

[Android client here](https://github.com/fvckyounimu/Whissage)

## Run
You need [go v1.14 minimum](https://golang.org/dl/) & [geth](https://geth.ethereum.org/docs/install-and-build/installing-geth)

#### Private network
[more detailed in official wiki](https://github.com/ethereum/go-ethereum/wiki/Private-network)

1. Creating The Genesis Block
You should change address (alloc) in genesis.json before run!

`geth --datadir /path_to_project/whissage/ init genesis.json`

2. Run geth

`bootnode --nodekey=boot.key`

`geth --rpc --shh --ws --wsapi web3,rpc,eth,net,shh --datadir . --networkid 2`

or copy url and:

`geth --rpc --shh --ws --wsapi web3,rpc,eth,net,shh --datadir . --networkid 2 --bootnodes enode://you_url`

#### Mainnet

`geth --rpc --shh --ws --wsapi web3,rpc,eth,net,shh` or `geth --syncmode "light" --rpc --shh --ws --wsapi web3,rpc,eth,net,shh`

#### Run app

`go run main.go` or only build `go build main.go`

## License
[MIT](https://github.com/semyon-dev/whissage/blob/master/LICENSE)
