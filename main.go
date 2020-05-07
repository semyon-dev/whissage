package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/whisper/whisperv6"
	"log"
	"os"
	"runtime"

	"github.com/ethereum/go-ethereum/whisper/shhclient"
)

var keyID string

func main() {

	go Subscribe()

	client, err := shhclient.Dial("ws://127.0.0.1:8546")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("we have a whisper connection")

	keyID, err = client.NewKeyPair(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Open("keys.txt")
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.WriteString(keyID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("keyID:", keyID)

	publicKey, err := client.PublicKey(context.Background(), keyID)
	if err != nil {
		log.Print(err)
	}

	fmt.Println("publicKey:", hexutil.Encode(publicKey))

	message := whisperv6.NewMessage{
		Payload:   []byte("Hello from Semyon!"),
		PublicKey: publicKey,
		TTL:       60,
		PowTime:   2,
		PowTarget: 2.5,
	}

	messageHash, err := client.Post(context.Background(), message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("messageHash:", messageHash)
	runtime.Goexit() // wait for goroutines to finish
}
