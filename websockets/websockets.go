package websockets

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var addr = flag.String("addr", "localhost:8081", "http service address")
var upgrader = websocket.Upgrader{} // use default options

func Start() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", main)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func main(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write error:", err)
			break
		}
	}
}
