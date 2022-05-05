package ocSocket

import (
	"encoding/base64"
	"encoding/json"
)

type Event struct {
	Name   string
	Data   string
	Client *Client `json:"-"`
}

type UUIDEvent struct {
	UUID string
}

func NewEvent(name string, data interface{}) (*Event, error) {
	event := &Event{
		Name: name,
	}
	var err error
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	event.Data = base64.StdEncoding.EncodeToString(bytes)
	return event, nil
}

func (event *Event) ParseData(target interface{}) error {
	data, err := base64.StdEncoding.DecodeString(event.Data)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}
