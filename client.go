package main

import (
  "github.com/gorilla/websocket"
)

//clientはチャットを行なっている一人のユーザーを表します。
type client struct {
  //socketはこのクライアントのためのWebSocketです。
  socket *websocket.Conn

  //sendはメッセージが送られるチャネルです。
  send chan []byte

  //roomはこのクライアントが参加しているチャットルームです。
  room *room
}
