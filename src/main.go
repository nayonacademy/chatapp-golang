package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

var upgrade = websocket.Upgrader{}
type Message struct {
	Email string `json:"email"`
	Username string `json:"username"`
	Message string `json:"message"`
}

func main()  {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnection)

	go handleMessage()
	log.Println("Listening server 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnection(w http.ResponseWriter, r *http.Request){
	ws, err := upgrade.Upgrade(w, r, nil)
	if err != nil{
		log.Fatal(err)
	}
	defer ws.Close()
	clients[ws] = true
	for{
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil{
			log.Printf("Error #{err}")
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}

}

func handleMessage(){
	for {
		msg := <-broadcast
		for client := range clients{
			err := client.WriteJSON(msg)
			if err != nil{
				log.Printf("Error #{err}")
				client.Close()
				delete(clients, client)
			}
		}
	}
}