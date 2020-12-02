// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hub

import (
	"bytes"
	"fmt"
	"github.com/xela07ax/toolsXela/ttp"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

}
// serveWs handles websocket requests from the peer.
func (h *Hub) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	h.Logx("-serveWs->init")
	client := &Client{hub: h, conn: conn, send: make(chan []byte, 256)}

	client.hub.register <- client

	// Разрешить сбор памяти, на которую ссылается вызывающий абонент, выполнив всю работу в
	// новых goroutines.
	go client.writePump()
	go client.readPump()
	h.Logx("-serveWs->init-end")
}
// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	c.hub.Logx("-Client.readPump->init")
	defer func() {
		c.hub.Logx("-Client.readPump->defer func[X]")
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		c.hub.Logx("-Client.readPump->for[circle]")
		_, message, err := c.conn.ReadMessage()
		c.hub.Logx(fmt.Sprintf("-Client.readPump->for-msg[%s]",message))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.Logx(fmt.Sprintf("-Client.readPump->for[broadcast]"))
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	c.hub.Logx("-Client.writePump->init")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.hub.Logx("-Client.writePump->defer func[X]")
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		c.hub.Logx("-Client.writePump->for[circle]")
		select {
		case message, ok := <-c.send:
			c.hub.Logx(fmt.Sprintf("-Client.writePump->for-message[%s]",message))
			c.hub.WebSocketOutput <- message
			//log.Printf("msg:%s\n",message)
			//time.Sleep(1*time.Second)
			//continue

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
			c.hub.Logx(fmt.Sprintf("-Client.writePump->for-message[%s]-ok",message))
		case <-ticker.C:
			c.hub.Logx(fmt.Sprintf("-Client.writePump->for-ticker.C-SetWriteDeadline"))
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
func readBodySimple(w http.ResponseWriter, r *http.Request) []byte {
	fmt.Println("Пришла команда по HTTP 1.1")
	b, err := ioutil.ReadAll(r.Body) // Считывание тело, ожидаем адрес сервера "localhost:5578"
	if err != nil {
		ertx := fmt.Sprintf("COM:Ошибка чтения тела: %s | ERTX:can't read body", err)
		fmt.Println(ertx)
		http.Error(w, ertx, http.StatusConflict) // 409
		return []byte{}
	}
	return b
}
type Notify struct {
	FuncName string
	Text string
	Status int
	Show bool
	UpdNum int
}

func resp (w http.ResponseWriter, r *http.Request, funcName string,text string, status int, show bool) {
	Notify := Notify {
		FuncName: funcName,
		Text:     text,
		Status:   status,
		Show:     show,
	}
	if err := ttp.Httpjson(w, r, Notify); err != nil {
		log.Fatalf("Критическая ошибка, не удалось отправить сообщение в UI: %s| %v", err,Notify)
	}

}
// serveWs handles websocket requests from the peer.
func sendMsg(loger chan <- [4]string, w http.ResponseWriter, r *http.Request) {
	//c.hub.Logx("-sendMsg->init")
	msgRaw := readBodySimple(w,r)
	if len(msgRaw) == 0 {
		return
	}
	// Реализация Writer-а
	fmt.Println("--W-process> начнем отправку")
	loger <- [4]string{"sendMsg","nil",string(msgRaw)}
	log.Printf("--W-true> отправили:%s\n", msgRaw)
	// ---
	ertx := "-sendMsg->true"
	resp(w,r,"sendMsg", ertx, 0, true)
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	hub.Logx("-serveWs->init")
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}

	client.hub.register <- client

	// Разрешить сбор памяти, на которую ссылается вызывающий абонент, выполнив всю работу в
	// новых goroutines.
	go client.writePump()
	go client.readPump()
	hub.Logx("-serveWs->init-end")
}
