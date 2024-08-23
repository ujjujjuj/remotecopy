package common

import "github.com/gorilla/websocket"

type ClientRole int
type MessageType int

const (
	Sender ClientRole = iota
	Receiver
)

const ()

type Client struct {
	Conn *websocket.Conn
	Role ClientRole
	Room string
}

type HelloMessage struct {
	Room string
}

type TextMessage struct {
	Text string
	Room string
}

func (c *Client) SendHello() error {
	message := HelloMessage{Room: c.Room}
	err := c.Conn.WriteJSON(message)

	return err
}

func (c *Client) SendText(text string) error {
	message := TextMessage{Text: text}
	err := c.Conn.WriteJSON(message)

	return err
}

func (c *Client) ReceiveHello() (*HelloMessage, error) {
	var msg HelloMessage
	err := c.Conn.ReadJSON(&msg)

	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (c *Client) ReceiveText() (*TextMessage, error) {
	var msg TextMessage
	err := c.Conn.ReadJSON(&msg)

	if err != nil {
		return nil, err
	}

	return &msg, nil
}
