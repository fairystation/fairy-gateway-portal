package main

import (
	"fmt"
	//"encoding/json"
	"github.com/gorilla/websocket"
	//"flag"
	"net/http"
	"html/template"
)

var upgrader = websocket.Upgrader{} // use default options

//var addr = flag.String("addr", "localhost:5866", "http service address")

func localEcho(w http.ResponseWriter, r *http.Request) {
	//quick fix for CORS
	// https://github.com/gorilla/websocket/issues/367
	upgrader.CheckOrigin = func(r *http.Request) bool { fmt.Println(r); return true }
	
	// c, err := upgrader.Upgrade(w, r, http.Header{
	// 	"Host": {"interhost"},
	// })

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(fmt.Sprint("upgrade:", err))
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println(fmt.Sprint("read:", err))
			break
		}
		fmt.Println(fmt.Sprint("recv: %s", message))
		err = c.WriteMessage(mt, message)
		if err != nil {
			fmt.Println(fmt.Sprint("write:", err))
			break
		}
	}
}

func localHome(w http.ResponseWriter, r *http.Request) {
 	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func localWS() {
	//flag.Parse()
	http.HandleFunc("/echo", localEcho)
	http.HandleFunc("/", localHome)
	fmt.Println( "---START WS ---")
	// it must be localhost , not to attract attention
	// fmt.Println(http.ListenAndServe("localhost:5866", nil))
	// but I am testing it with Docker , so I need it public
	fmt.Println(http.ListenAndServe("0.0.0.0:5866", nil))
}
var homeTemplate = template.Must(template.New("").Parse(`<html><head><title>Websocket Echo</title></head><body><script></script>{{.}}<body></html>`))
