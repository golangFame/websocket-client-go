package main

import (
	"fmt"
	websocketClient "websocketclient/client"
)

func main() {
	fmt.Println("Starting WebSocket Client Application")
	uri := "wss://websocket.goferhiro.repl.co/v1/ws"
	//uri = "wss://socketsbay.com/wss/v2/2/demo/"
	websocketClient.StartClient(uri)
}
