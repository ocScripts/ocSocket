package ocSocket

import (
	"errors"
	"net"
)

type OCSocket struct {
	Clients         map[string]*Client
	Address         string
	EventChan       chan *Event
	ServerCallbacks map[string]serverCallback
}

func NewSocket(address string) *OCSocket {
	socket := &OCSocket{
		Clients:         make(map[string]*Client),
		Address:         address,
		ServerCallbacks: make(map[string]serverCallback),

		EventChan: make(chan *Event),
	}

	socket.RegisterCallback("ping", pong)
	socket.RegisterCallback("getUUID", getUUID)

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
		go newClient.listen(socket)

		<-newClient.ready
	}
}

func (socket *OCSocket) SendClientEvent(UUID string, event *Event) error {
	if client, ok := socket.Clients[UUID]; ok && client.IsOpen() {
		return client.SendEvent(event)
	}
	return errors.New("invalid client UUID")
}
