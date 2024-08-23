package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ujjujjuj/remotecopy/internal/client"
	"github.com/ujjujjuj/remotecopy/internal/common"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s <server_address> <room_name> <sender | receiver>\n", os.Args[0])
		os.Exit(1)
	}

	serverAddr := os.Args[1]
	roomName := os.Args[2]
	roleStr := os.Args[3]
	if roleStr != "sender" && roleStr != "receiver" {
		fmt.Println("Role should either be \"sender\" or \"receiver\"")
		os.Exit(1)
	}
	role := common.Sender
	if roleStr == "receiver" {
		role = common.Receiver
	}

	if role == common.Sender {
		content, err := client.Copy()
		if err != nil {
			fmt.Println("Error getting clipboard:", err)
			os.Exit(1)
		}

		msg := common.TextMessage{Text: content, Room: roomName}
		body, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("Error creating request body:", err)
			os.Exit(1)
		}

		r, err := http.NewRequest("POST", fmt.Sprintf("http://%s/send", serverAddr), bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Request error:", err)
			os.Exit(1)
		}

		r.Header.Add("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(r)
		if err != nil {
			fmt.Println("Send error:", err)
			os.Exit(1)
		}
		defer res.Body.Close()

		buf := new(strings.Builder)
		_, err = io.Copy(buf, res.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			os.Exit(1)
		}

		fmt.Println(buf.String())

	} else {
		client.Test()
		
		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws", serverAddr), nil)
		if err != nil {
			fmt.Println("Dial error:", err)
			os.Exit(1)
		}
		defer conn.Close()

		serverClient := common.Client{Conn: conn, Role: role, Room: roomName}

		err = serverClient.SendHello()
		if err != nil {
			fmt.Println("Hello error:", err)
			os.Exit(1)
		}

		for {
			text, err := serverClient.ReceiveText()
			if err != nil {
				fmt.Println("Receive error: ", err)
				break
			}

			fmt.Println("Received:", text.Text)
			time.Sleep(time.Second * 5)
			client.Paste(text.Text)
		}
	}

}
