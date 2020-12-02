// Copyright 2020 The Xela07ax WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hub

import (
	"fmt"
	"time"
)

type Hub struct {
	// Зарегистрированный клиент.
	clients map[*Client]bool
	Input chan []byte
	// Входящие сообщения от клиентов.
	broadcast chan []byte

	// Регистрируйте запросы от клиентов.
	register chan *Client

	// Отмените регистрацию запросов от клиентов.
	unregister chan *Client
	Debug bool
	WebSocketOutput chan []byte
}

func NewHub(debug bool) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		Input:      make(chan []byte,100), // Отправляем клиентам подключившихся по вэбсокетам
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		Debug: 		debug,
		WebSocketOutput: make(chan []byte,100),
	}
}


func (h *Hub) Run() {
		h.Logx("-hub.run->init")
	for {
		h.Logx("-hub.run->for[circle-start]")
		select {
		case client := <-h.register:
			h.Logx("-hub.run->select[h.register]")
			h.clients[client] = true
		case client := <-h.unregister:
			h.Logx("-hub.run->select[h.unregister]")
			if _, ok := h.clients[client]; ok {
				h.Logx("-hub.run->select[h.unregister]-ok")
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.Input:
			h.Logx("-hub.run->select[h.Input]new")
			for client := range h.clients {
				h.Logx("-hub.run->select[h.Input]for")
				select {
				case client.send <- message:
					h.Logx(fmt.Sprintf("-hub.run->select[h.Input]for-select[%s]",message))
				default:
					h.Logx("-hub.run->select[h.Input]for-select-default")
					close(client.send)
					delete(h.clients, client)
				}
			}
		case message := <-h.broadcast:
			// Входящие сообщения от клиентов.
			// Собственно работать будем тут
			h.Logx("-hub.run->select[h.broadcast]")
			for client := range h.clients {
				h.Logx("-hub.run->select[h.broadcast]for")
				select {
				case client.send <- message:
					h.Logx(fmt.Sprintf("-hub.run->select[h.broadcast]for-select[%s]",message))

				default:
					h.Logx("-hub.run->select[h.broadcast]for-select-default")
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) Logx(txt string)  {
	if h.Debug {
		fmt.Printf("%s|%v\n",txt,Getime())
		time.Sleep(1*time.Second)
	}
}
func Getime()string  {
	return time.Now().Format("2006-01-02 15:04:05")
}