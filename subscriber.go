package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/whisper/shhclient"
	"github.com/ethereum/go-ethereum/whisper/whisperv6"
	"log"
	"os"
	"runtime"
)

func Subscribe() {

	client, err := shhclient.Dial("ws://127.0.0.1:8546")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("we have a whisper connection from Subscriber")

	messages := make(chan *whisperv6.Message)
	criteria := whisperv6.Criteria{
		PrivateKeyID: keyID,
	}
	sub, err := client.SubscribeMessages(context.Background(), criteria, messages)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Fatal(err)
			case message := <-messages:
				fmt.Printf(string(message.Payload))
				os.Exit(0)
			}
		}
	}()
	runtime.Goexit() // wait for goroutines to finish
}
