package main

import (
	"fmt"
	"net/http"
	"os"

	server "github.com/ujjujjuj/remotecopy/internal/server"
)

func main() {
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	rm := server.CreateRoomManager()

	http.HandleFunc("/send", rm.NewRoomMessage)
	http.HandleFunc("/ws", rm.HandleWebsocket)

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println("Failed to run server:", err)
	}
}
