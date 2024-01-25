package main

import "github.com/gorilla/websocket"

type Client struct {
	Name string
	Room string
	Conn *websocket.Conn
}
