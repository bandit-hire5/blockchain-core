package core

import (
	"encoding/json"
)

type Message struct {
	Type string `json:"type"`
	Body string `json:"body"`
}

func (self *Message) String() string {
	dataString, err := json.Marshal(self)
	if err != nil {
		return ""
	}

	return string(dataString)
}
