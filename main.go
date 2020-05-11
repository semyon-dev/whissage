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
var connections map[*websocket.Conn]*websocket.Conn

func main() {

	// логирование в консоль
	log.SetOutput(os.Stdout)

	var err error
	client, err = shhclient.Dial(url)
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	fmt.Println("We have a whisper connection")

	if len(config.TestKey) == 0 {
		keyID, err = client.NewKeyPair(context.Background())
		if err != nil {
			log.Fatal("NewKeyPair error: ", err)
		}
		fmt.Println("keyID:", keyID)
		file, err := os.OpenFile("config/keys.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Fatal("Open file error: ", err)
		}
		_, err = file.WriteString(keyID + "\n")
		if err != nil {
			log.Fatal("WriteString error: ", err)
		}
		config.TestKey = keyID
	}

	publicKey, err = client.PublicKey(context.Background(), config.TestKey)
	if err != nil {
		log.Fatal("PublicKey fatal error", err)
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
		log.Println("upgrade ws error:", err)
		return
	}
	fmt.Println("new conn: ", c)
	connections[c] = c
	defer c.Close()
	defer delete(connections, c)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		fmt.Println("Получили от клиента: ", message)

		whisperMessage := whisperv6.NewMessage{
			Payload:   message,
			PublicKey: publicKey,
			TTL:       60,
			PowTime:   2,
			PowTarget: 2.5,
		}
		messageHash, err := client.Post(context.Background(), whisperMessage)
		if err != nil {
			fmt.Println("Не удалось отправить сообщение: ", err)
		}
		fmt.Println("message hash: ", messageHash)
	}
}

// слушаем сообщения от ethereum whisper
func Subscribe() {

	messages := make(chan *whisperv6.Message)
	criteria := whisperv6.Criteria{
		PrivateKeyID: config.TestKey,
	}

	sub, err := client.SubscribeMessages(context.Background(), criteria, messages)
	if err != nil {
		log.Fatal("Не удалось подписаться на сообщения: ", err)
	}

	for {
		select {
		case err := <-sub.Err():
			fmt.Println("ошибка в Subscribe: ", err)
		case message := <-messages:
			fmt.Println("Получили сообщение через Subscribe: ", string(message.Payload))
			for _, c := range connections {
				err = c.WriteMessage(websocket.TextMessage, message.Payload)
				if err != nil {
					fmt.Println("Ошибка при отправке сообщения: ", err)
					break
				} else {
					fmt.Println("Успешно отправили сообщение клиенту: ", c.RemoteAddr())
				}
			}
		}
	}
}
