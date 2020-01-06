package main

import (
	"encoding/json"
	"fmt"
	"log"
)

const ARENA_LOGIN = 101

type ArenaMessage struct {
	Bm struct {
		Payload   map[string]interface{} `json:"payload"`
		Pid       int                    `json:"pid"`
		Csq       int                    `json:"csq"`
		SessionID string                 `json:"sessionId"`
	} `json:"bm"`
}

func (msg *ArenaMessage) mkCmd(command int, params map[string]interface{}) []byte {
	csq++
	msg.Bm.Pid = command
	msg.Bm.Csq = csq
	msg.Bm.Payload = params
	msg.Bm.SessionID = sid
	data, err := json.Marshal(msg)
	if err != nil {
		log.Panic("Undefined behaivour")
	}
	header := []byte(fmt.Sprintf("%d:", len(data)))
	data = append(header, data...)
	log.Println(string(data))
	return data
}
