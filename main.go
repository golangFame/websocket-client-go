package main

import (
	"fmt"
  "os"
	websocketClient "websocketclient/client"
)

func main() {
	fmt.Println("Starting WebSocket Client Application")
	uri := os.Getenv("socket_url")
  if url==""{
    url="wss://socketsbay.com/wss/v2/2/demo/"
  }
	websocketClient.StartClient(uri)
}
