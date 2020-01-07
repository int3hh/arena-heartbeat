package main

import (
	"encoding/json"
	"fmt"
	"log"
)

const ARENA_LOGIN = 101
const ARENA_HEARTBEAT = 100
const ARENA_SUBSCRIBE = 108

type ArenaMessage struct {
	Bm struct {
		Payload   map[string]interface{} `json:"payload,omitempty"`
		SessionID string                 `json:"sessionId,omitempty"`
		User      int32                  `json:"user,omitempty"`
		Pid       int                    `json:"pid"`
		Csq       int                    `json:"csq"`
	} `json:"bm"`
}

func (msg *ArenaMessage) mkCmd(command int, params map[string]interface{}) []byte {
	csq++
	msg.Bm.Pid = command
	msg.Bm.Csq = csq
	msg.Bm.Payload = params
	msg.Bm.SessionID = sid
	msg.Bm.User = uid
	data, err := json.Marshal(msg)
	if err != nil {
		log.Panic("Undefined behaivour")
	}
	header := []byte(fmt.Sprintf("%d:", len(data)))
	data = append(header, data...)
	log.Println(string(data))
	return data
}

func processMessage(rawMessage []byte) [][]byte {
	var reply [][]byte
	var msg ArenaMessage

	err := json.Unmarshal(rawMessage, &msg)
	if err == nil {
		if msg.Bm.Pid == 101 {
			userInfo := msg.Bm.Payload["user"].(map[string]interface{})
			sid = userInfo["sid"].(string)
			uid = msg.Bm.User
			for _, symbol := range symbols {
				reply = append(reply, (&ArenaMessage{}).mkCmd(ARENA_SUBSCRIBE, map[string]interface{}{"@class": "p.Subscription",
					"exchange": "BVB", "symbol": symbol}))
			}
			return reply
		}
	} else {
		log.Panic("Unable to unmarshall", err)
	}

	return nil
}
