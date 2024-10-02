package torpacket

import (
	"encoding/json"
	"errors"
	"reflect"
)

// MessageMap is a map that associates message types with the corresponding structs.
var MessageMap = map[int]interface{}{
	GetAES:   &GetAesMsg{},
	Redirect: &RedirectMsg{},
	Receive:  &ReceiveMsg{},
	Destroy:  nil, // No data for Destroy, so we set it to nil
	Ack:      &AckMsg{},
}

func DeserializeMessage(rawMsg RawMessage) (ParsedMessage, error) {
	msgStruct, ok := MessageMap[rawMsg.Code]

	if !ok {
		return ParsedMessage{}, errors.New("unknown message type")
	}

	if msgStruct == nil {
		return ParsedMessage{Code: rawMsg.Code, Data: nil}, nil
	}

	msgInstance := reflect.New(reflect.TypeOf(msgStruct).Elem()).Interface()

	if err := json.Unmarshal([]byte(rawMsg.JsonData), msgInstance); err != nil {
		return ParsedMessage{}, err
	}

	return ParsedMessage{
		Code: rawMsg.Code,
		Data: msgInstance,
	}, nil
}
