package main

type room struct {
  //forwardは他のクライアントに転送するためのメッセージを保持するチャンネルです。
  forward chan []byte

  //joinはチャットルームに参加しようとしているクライアントのためのチャネルです。
  join chan *client

  //leaveはチャットルームから退室しようとしているクライアントのためのチャネルです。
  leave chan *client

  //clientsには在室している全てのクライアントが保持されます。
  clients map[*client]bool
}
