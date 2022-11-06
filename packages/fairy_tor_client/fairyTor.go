package main

import (
	"./go-telnet"
	"./tengo"
	"./tengo/stdlib"
	"time"
	"github.com/reiver/go-oi"
	"strings"
	"fmt"
	"context"
	//"encoding/json"
	"github.com/gorilla/websocket"
	"os"
	//"flag"
	"net/http"
	"html/template"
)

var TengoHandler telnet.Handler = internalEchoHandler{}

type internalEchoHandler struct{}

var upgrader = websocket.Upgrader{} // use default options

func (handler internalEchoHandler) ServeTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {
	var fullText string
	var checkPingText string
	var buffer [1]byte
	p := buffer[:]
	_, _ = r.Read(p)
	for {
		n, err := r.Read(p)
		if n > 0 {
			fullText = fullText + string(p[:n])
			fmt.Println( fullText )
			fullText = strings.Replace(fullText, "  ", " ", -1)
			checkPingText = checkPingText + string(p[:n])

			i := strings.Index(checkPingText, "fairy client ping")
			if i > -1 {
				fmt.Println( "!!!! Fairy Client !!!!" )
				time.Sleep(30 * time.Millisecond)
				checkPingText = ""

				fmt.Println(" !!!! RUN" + fullText )
				script := tengo.NewScript([]byte(fullText))
				script.SetImports(stdlib.GetModuleMap("fmt", "tor"))

				// run the script
				compiled, err := script.RunContext(context.Background())
				if err != nil {
					fmt.Println(err)
					panic(err)
				}

				fmt.Println( "!!!! Post !!!!" )
				// retrieve values
				mul := compiled.Get("fteftefte").String()

				if mul == "" {
					oi.LongWrite(w,[]byte("///ftefteftecnf"+ "\n\r"))
				} else {
					oi.LongWrite(w,[]byte(mul + "\n\r"))
				}

				fullText = ""
			}
		}

		if nil != err {
			oi.LongWrite(w,[]byte("End on rocks"))
			break
		}
	}
	time.Sleep(30 * time.Millisecond)
}

func serveSimple() {
	var handler telnet.Handler = TengoHandler
	fmt.Println( "---START SIMPLE ---")
	err := telnet.ListenAndServeSimple(":5066", handler)
	if nil != err {
		//@TODO: Handle this error better.
		panic(err)
	}
}

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

func main() {

    err_dir := os.RemoveAll("./data-dir-*")
    if err_dir != nil {
        fmt.Println( "---REMOVE DIR FAILED ---")
    }

	var handler telnet.Handler = TengoHandler
	fmt.Println( "---START TOR ---") 
	go serveSimple();
	go localWS();
	err := telnet.ListenAndServe("5566", handler)
	if nil != err {
		//@TODO: Handle this error better.
		panic(err)
	}
}
