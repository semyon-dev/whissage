package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/whisper/shhclient"
	"github.com/ethereum/go-ethereum/whisper/whisperv6"
	"github.com/gorilla/websocket"
	"github.com/semyon-dev/whissage/config"
	"log"
	"net/http"
	"os"
)

var keyID string
var url = "ws://127.0.0.1:8546"
var publicKey []byte
var client *shhclient.Client

func main() {

	log.SetOutput(os.Stdout)

	var err error
	client, err = shhclient.Dial(url)
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

	publicKey, err = client.PublicKey(context.Background(), config.TestKey)
	if err != nil {
		log.Print("PublicKey", err)
	}

	fmt.Println("publicKey:", hexutil.Encode(publicKey))

	go Subscribe() // Subscribe for messages

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", mainHandler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options

func mainHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		log.Printf("получили: %s", message)

		whisperMessage := whisperv6.NewMessage{
			Payload:   message,
			PublicKey: publicKey,
			TTL:       60,
			PowTime:   2,
			PowTarget: 2.5,
		}
		messageHash, err := client.Post(context.Background(), whisperMessage)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("messageHash: ", messageHash)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write error:", err)
			break
		} else {
			fmt.Println("успешно отправили")
		}
	}
}
