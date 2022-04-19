package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"net/url"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	wsconn    *websocket.Conn
	configStr string
}

func NewWebSocketClient(host, channel string) *WebSocketClient {
	conn := WebSocketClient{}

	u := url.URL{Scheme: "wss", Host: host, Path: channel} //scheme has to be wss not ws 
	conn.configStr = u.String()
	return &conn
}

func (conn *WebSocketClient) Connect() error {
	if conn.configStr == "" {
		err := fmt.Errorf("Connection Url is empty")
		conn.log("Connect", err, "Empty Url ")
		return err
	}

	ws, res, err := websocket.DefaultDialer.Dial(conn.configStr, nil)
	defer res.Body.Close()
	if err != nil {
		conn.log("Connect", err, res.Status)
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(body))
		return err
	}
	conn.wsconn = ws
	return nil
}

func (conn *WebSocketClient) Listen() {
	//ticker := time.NewTicker(time.Second) //better CPU usage
	for {

		for {
			_, byteMsg, err := conn.wsconn.ReadMessage()
			if err != nil {
				conn.log("Listen", err, "Stopping Connection")
				//	conn.Stop()
				continue
			}
			conn.log("Listen", nil, fmt.Sprintf("Recieved msg %s\n", byteMsg))
		}
	}
}

func (conn *WebSocketClient) Stop() {
	if conn.wsconn != nil {
		conn.wsconn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		conn.wsconn.Close()
		conn.wsconn = nil
	}
}

func StartClient(u string) {
	uri, _ := url.Parse(u)
	fmt.Println(uri.Host, uri.Path)
	client1 := NewWebSocketClient(uri.Host, uri.Path)
	for {
		if err := client1.Connect(); err == nil {
			client1.log("main", nil, "Connection successful")
			break
		}
	}
	go func() {
		client1.log("RandomWrite", nil, "Starting")
		for i := 0; ; i++ {
			payload := []byte("V" + strconv.Itoa(i))
			client1.log("RandomWrite", nil, string(payload))
			client1.Write(payload)
		}
	}()
	client1.log("main", nil, "Listening to ws server")
	client1.Listen()
}

func main() {
	client1 := NewWebSocketClient("https://websocket.goferhiro.repl.co", "/v1/ws")
	client1.log("main", nil, "Listening to ws server")
	go func() {
		client1.log("RandomWrite", nil, "Starting")
		for i := 0; ; i++ {
			payload := []byte("V" + strconv.Itoa(i))
			client1.log("RandomWrite", nil, string(payload))
			client1.Write(payload)
		}
	}()
	client1.Listen()
}

func (conn *WebSocketClient) log(f string, err error, msg string) {
	if err != nil {
		fmt.Printf("Error in func: %s, err: %v, msg: %s\n", f, err, msg)
	} else {
		fmt.Printf("Log in func: %s, %s\n", f, msg)
	}
}

func (conn *WebSocketClient) Write(payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		conn.log("Write", err, "Error while JSON Marshaling")
		return err
	}

	ws := conn.wsconn

	if ws == nil {
		err := fmt.Errorf("conn.ws is nil")
		conn.log("Write", err, "Invalid connection")
		return err
	}

	if err := ws.WriteMessage(
		websocket.TextMessage,
		data,
	); err != nil {
		conn.log("Write", nil, "WebSocket Write Error")
	}
	conn.log("Write", nil, fmt.Sprintf("sent: %s", data))
	return nil
}
  