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
	"os"
	//"flag"
)

var TengoHandler telnet.Handler = internalEchoHandler{}

type internalEchoHandler struct{}

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
