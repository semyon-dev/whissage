# whissage
The backend of blockchain-based messenger built on ethereum whisper.

[Сurrently in active development]

## Run & build
You need go v1.14 minimum & geth 

#### Private network
1. Creating The Genesis Block

`geth --datadir /path_to_project/whissage/ init genesis.json`

2. Run geth

`geth --datadir path/to/custom/data/folder --networkid 2`

#### Mainnet

`geth --rpc --shh --ws` or `geth --syncmode "light" --rpc --shh --ws`

#### Run app

`go run main.go` or `go build main.go`

## License
[MIT](https://github.com/semyon-dev/whissage/blob/master/LICENSE)
