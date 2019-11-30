package main

import (
	"harmony-server/socket"
	"harmony-server/socket/event"
	"net/http"
)

func handleSocket(w http.ResponseWriter, r *http.Request) {
	ws := socket.NewSocket(w, r)
	ws.Bind("login", event.OnLogin)
	ws.Bind("register", event.OnRegister)
	ws.Bind("getservers", event.OnGetServers)
}