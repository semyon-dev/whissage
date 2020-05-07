# whissage (whisper message)
The backend of blockchain-based messenger built on ethereum whisper.

**Сurrently in active development**

## Run & build
You need go v1.14 minimum & geth 

#### Private network
1. Creating The Genesis Block
You should change address (alloc) in genesis.json before run!

`geth --datadir /path_to_project/whissage/ init genesis.json`

2. Run geth

`bootnode --nodekey=boot.key`

than copy url and:

`geth --rpc --shh --ws --wsapi web3,rpc,eth,net,shh --datadir . --networkid 2 --bootnodes enode://you_url`

#### Mainnet

`geth --rpc --shh --ws` or `geth --syncmode "light" --rpc --shh --ws`

#### Run app

`go run main.go` or `go build main.go`

## License
[MIT](https://github.com/semyon-dev/whissage/blob/master/LICENSE)
