package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/whisper/shhclient"
	"github.com/ethereum/go-ethereum/whisper/whisperv6"
	"github.com/semyon-dev/whissage/config"
	"log"
)

func Subscribe() {

	client, err := shhclient.Dial(url)
	if err != nil {
		log.Fatal("connection: ", err)
	}
	fmt.Println("we have a whisper connection from subscriber")

	messages := make(chan *whisperv6.Message)
	criteria := whisperv6.Criteria{
		PrivateKeyID: config.TestKey,
	}

	sub, err := client.SubscribeMessages(context.Background(), criteria, messages)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Println("ошибка: ", err)
		case message := <-messages:
			fmt.Print("we get a message: ", string(message.Payload))
		}
	}
}
