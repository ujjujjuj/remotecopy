package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	common "github.com/ujjujjuj/remotecopy/internal/common"
)

type RoomManager struct {
	Receivers map[string]*common.Client
}

func CreateRoomManager() *RoomManager {
	return &RoomManager{Receivers: make(map[string]*common.Client)}
}

func (rm *RoomManager) joinRoom(roomName string, client *common.Client) {
	rm.Receivers[roomName] = client
	log.Printf("new receiver in room %s\n", roomName)
}

func (rm *RoomManager) deleteRoom(roomName string) {
	delete(rm.Receivers, roomName)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (rm *RoomManager) HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := common.Client{Conn: conn}

	hello, err := client.ReceiveHello()
	if err != nil {
		return
	}
	client.Room = hello.Room

	rm.joinRoom(hello.Room, &client)
}

func (rm *RoomManager) NewRoomMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var textMessage common.TextMessage
	err := json.NewDecoder(r.Body).Decode(&textMessage)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	client := rm.Receivers[textMessage.Room]
	if client == nil {
		json.NewEncoder(w).Encode(map[string]string{"message": "No client in the room"})
		return
	}

	err = client.SendText(textMessage.Text)
	if err != nil {
		rm.disconnectClient(client)
		json.NewEncoder(w).Encode(map[string]string{"message": "Error. Client might have disconnected"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Sent message"})
}

func (rm *RoomManager) disconnectClient(client *common.Client) {
	defer client.Conn.Close()
	rm.deleteRoom(client.Room)
}
