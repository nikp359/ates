package estream

import (
	"encoding/json"
)

//TODO: rewrite to code generate

type UserCreatedHandler func(Meta, UserCreatedPayload) error

func (h UserCreatedHandler) RawHandler() RawMessageHandler {
	return func(meta Meta, raw json.RawMessage) error {
		var d UserCreatedPayload
		if err := d.UnmarshalJSON(raw); err != nil {
			return err
		}
		return h(meta, d)
	}
}

type UserUpdatedHandler func(Meta, UserUpdatedPayload) error

func (h UserUpdatedHandler) RawHandler() RawMessageHandler {
	return func(meta Meta, raw json.RawMessage) error {
		var d UserUpdatedPayload
		if err := d.UnmarshalJSON(raw); err != nil {
			return err
		}
		return h(meta, d)
	}
}

type UserDeletedHandler func(Meta, UserDeletedPayload) error

func (h UserDeletedHandler) RawHandler() RawMessageHandler {
	return func(meta Meta, raw json.RawMessage) error {
		var d UserDeletedPayload
		if err := d.UnmarshalJSON(raw); err != nil {
			return err
		}
		return h(meta, d)
	}
}
