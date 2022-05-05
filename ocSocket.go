package ocSocket

import (
	"errors"
	"net"
)

type OCSocket struct {
	Clients   map[string]*Client
	Address   string
	EventChan chan *Event
}

func NewSocket(address string) *OCSocket {
	socket := &OCSocket{
		Clients: make(map[string]*Client),
		Address: address,

		EventChan: make(chan *Event),
	}

	return socket
}

func (socket *OCSocket) Open() error {
	listen, err := net.Listen("tcp", socket.Address)
	if err != nil {
		return err
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}

		newClient := newClient(conn)

		socket.Clients[newClient.UUID] = newClient
		go newClient.listen(socket.EventChan)

		<-newClient.ready

		uuid := &UUIDEvent{
			UUID: newClient.UUID,
		}

		event, err := NewEvent("uuid", uuid)
		if err != nil {
			return err
		}

		err = newClient.SendEvent(event)
		if err != nil {
			return err
		}
	}
}

func (socket *OCSocket) SendClientEvent(UUID string, event *Event) error {
	if client, ok := socket.Clients[UUID]; ok && client.IsOpen() {
		return client.SendEvent(event)
	}
	return errors.New("invalid client UUID")
}
