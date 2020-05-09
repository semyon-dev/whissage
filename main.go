package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/whisper/shhclient"
	"github.com/ethereum/go-ethereum/whisper/whisperv6"
	"github.com/semyon-dev/whissage/config"
	"github.com/semyon-dev/whissage/websockets"
	"log"
	"os"
	"time"
)

var keyID string
var url = "ws://127.0.0.1:8546"

func main() {

	log.SetOutput(os.Stdout)

	go websockets.Start()

	client, err := shhclient.Dial(url)
	if err != nil {
		log.Fatal("connection: ", err)
	}
	fmt.Println("we have a whisper connection")

	if len(config.TestKey) == 0 {
		keyID, err = client.NewKeyPair(context.Background())
		if err != nil {
			log.Fatal("NewKeyPair: ", err)
		}
		fmt.Println("keyID:", keyID)
		file, err := os.OpenFile("config/keys.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Fatal("Open file: ", err)
		}
		_, err = file.WriteString(keyID + "\n")
		if err != nil {
			log.Fatal("WriteString: ", err)
		}
		config.TestKey = keyID
	}

	publicKey, err := client.PublicKey(context.Background(), config.TestKey)
	if err != nil {
		log.Print("PublicKey", err)
	}

	fmt.Println("publicKey:", hexutil.Encode(publicKey))

	go Subscribe() // Subscribe for messages

	for {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("enter a message: ")
		body, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		message := whisperv6.NewMessage{
			Payload:   []byte(body),
			PublicKey: publicKey,
			TTL:       60,
			PowTime:   2,
			PowTarget: 2.5,
		}
		messageHash, err := client.Post(context.Background(), message)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("messageHash: ", messageHash)
	}

	// runtime.Goexit() // wait for goroutines to finish
}
