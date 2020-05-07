package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/whisper/shhclient"
)

func main() {
	client, err := shhclient.Dial("ws://127.0.0.1:8546")
	if err != nil {
		log.Fatal(err)
	}

	_ = client // we'll be using this in the next section
	fmt.Println("we have a whisper connection")
}
