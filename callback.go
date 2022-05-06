package ocSocket

import (
	"encoding/json"
	"fmt"
	"math"
)

type serverCallback func(client *Client, data string, cb func(data string))
type clientCallback func(data string)

type CallbackData struct {
	ID   int
	Name string
	Data string
}

func (socket *OCSocket) triggerCallbackHandler(event *Event) {
	eventData := &CallbackData{}
	event.ParseData(eventData)
	data := eventData.Data
	cb, ok := socket.ServerCallbacks[eventData.Name]
	if !ok {
		return
	}

	cb(event.Client, data, func(data2 string) {
		toSend, err := NewEvent("serverCallback", &CallbackData{
			ID:   eventData.ID,
			Name: eventData.Name,
			Data: data2,
		})
		if err != nil {
			fmt.Println(err)
			return
		}

		err = event.Client.SendEvent(toSend)
		if err != nil {
			fmt.Println(err)
			return
		}
	})
}

func (socket *OCSocket) RegisterCallback(name string, cb serverCallback) {
	socket.ServerCallbacks[name] = cb
}

func (client *Client) TriggerClientCallback(name string, data interface{}, callback clientCallback) error {
	client.clientCallbacks[client.callbackID] = callback
	dBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	event, err := NewEvent("triggerClientCallback", &CallbackData{
		ID:   client.callbackID,
		Name: name,
		Data: string(dBytes),
	})
	if err != nil {
		return err
	}

	err = client.SendEvent(event)
	if err != nil {
		return err
	}

	client.callbackID++
	if client.callbackID == math.MaxInt {
		client.callbackID = 0
	}
	return nil
}

func (socket *OCSocket) clientCallbackHandler(event *Event) {
	data := &CallbackData{}
	err := event.ParseData(data)
	if err != nil {
		fmt.Println(err)
	}
	event.Client.clientCallbacks[data.ID](data.Data)
	event.Client.clientCallbacks[data.ID] = nil
}

func pong(client *Client, data string, cb func(string)) {
	cb(data)
}

func getUUID(client *Client, data string, cb func(string)) {
	returnData := map[string]string{
		"UUID": client.UUID,
	}

	toReturn, err := json.Marshal(returnData)
	if err != nil {
		fmt.Println(err)
		return
	}
	cb(string(toReturn))
}
