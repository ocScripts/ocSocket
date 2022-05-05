package ocSocket

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"strings"

	"github.com/google/uuid"
)

type Client struct {
	UUID  string
	conn  net.Conn
	ready chan bool
	open  bool
}

func newClient(con net.Conn) *Client {
	client := &Client{
		conn:  con,
		ready: make(chan bool),
		UUID:  uuid.New().String(),
		open:  false,
	}

	return client
}

func (client *Client) listen(eventChan chan *Event) {
	client.ready <- true
	client.open = true
	for client.open {
		event := &Event{}
		str, err := bufio.NewReader(client.conn).ReadString('\n')
		if err != nil {
			client.Close()
			return
		}
		err = json.Unmarshal([]byte(str), event)
		if err != nil {
			log.Printf("Error invalid event from client id %s\n%s", client.UUID, err)
			client.Close()
			return
		}
		event.Client = client

		eventChan <- event
	}
}

func (client *Client) SendEvent(event *Event) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	formated := strings.ReplaceAll(string(bytes), "\n", "")

	_, err = client.conn.Write([]byte(formated))
	if err != nil {
		return err
	}
	// Send a newline to indicate the end of the packet
	_, err = client.conn.Write([]byte("\n"))
	return err
}

func (client *Client) Close() error {
	client.open = false
	return client.conn.Close()
}

func (client *Client) IsOpen() bool {
	return client.open
}
