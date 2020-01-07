package main // import "github.com/int3hh/arena-heartbeat"

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var csq int
var c *websocket.Conn
var sid string
var pairs []string
var done chan struct{}

func main() {
	banner := figure.NewFigure("Arena-HB", "", true)
	banner.Print()
	err := godotenv.Load()
	if err != nil {
		log.Panic("Unable to load .env file")
	}

	arenaUser := os.Getenv("ARENA_USER")
	arenaPass := os.Getenv("ARENA_PASS")
	arenaHost := os.Getenv("ARENA_HOST")
	rawPairs := os.Getenv("PAIRS")
	if len(rawPairs) > 0 {
		pairs = strings.Split(rawPairs, ",")
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	running := false
	sid = ""
	for {
		if !running {
			done = make(chan struct{})
			u := url.URL{Scheme: "wss", Host: arenaHost, Path: "/ws/channel"}
			log.Println("Connecting to ", u.String())
			headers := http.Header{
				"Origin":                 {"https://testg.arenaxt.ro"},
				"Sec-WebSocket-Protocol": {"TEXT"},
			}
			c, _, err = websocket.DefaultDialer.Dial(u.String(), headers)
			if err != nil {
				log.Println("dial error, waiting 2 minutes :", err)
				time.Sleep(time.Minute * 2)
				continue
			}
			running = true
			defer c.Close()
			go func() {
				defer close(done)
				c.WriteMessage(websocket.TextMessage, (&ArenaMessage{}).mkCmd(ARENA_LOGIN, map[string]interface{}{"@class": "p.Login",
					"device": "WebT", "username": arenaUser, "password": arenaPass, "totp": 0}))
				for {
					log.Println("reading message...")
					_, message, err := c.ReadMessage()
					log.Println("read message ...")
					if err != nil {
						log.Println("read:", err)
						return
					}
					log.Println("recv: ", string(message))
					processMessage(message)
				}
			}()
		}

		select {
		case <-done:
			running = false
			log.Println("Socket disconnected ... ")

		case <-ticker.C:
			log.Println("Sending ping ... ")
			c.WriteMessage(websocket.TextMessage, (&ArenaMessage{}).mkCmd(ARENA_HEARTBEAT, nil))

		case <-interrupt:
			log.Println("Sigterm received quitting !")
			os.Exit(1)
		}
	}

}
