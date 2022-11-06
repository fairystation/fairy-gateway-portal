package stdlib

import (
	"fmt"
	"os"
	"log"
	"../"
	"strings"
	"time"
	"regexp"
	"golang.org/x/net/proxy"
	"io/ioutil"
)

var torModule = map[string]tengo.Object{
	"intercall":   &tengo.UserFunction{Name: "intercall", Value: torIntercall},
	"extercall":   &tengo.UserFunction{Name: "extercall", Value: torExtercall},
	"getstore":   &tengo.UserFunction{Name: "getstore", Value: getStore},
	"setstore":   &tengo.UserFunction{Name: "setstore", Value: setStore},
}

func getStore(args ...tengo.Object) (ret tengo.Object, err error) {

	fileName , _ := tengo.ToString(args[0])
	
	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	newStr := reg.ReplaceAllString(fileName, "")

	data, err := ioutil.ReadFile(newStr + ".log")

	if err != nil {
	  return &tengo.String{Value: ""}, nil
	}
	return &tengo.String{Value: string(data)}, nil
}

func setStore(args ...tengo.Object) (ret tengo.Object, err error) {

	fileName , _ := tengo.ToString(args[0])
	
	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	newStr := reg.ReplaceAllString(fileName, "")

	f, err := os.Create(newStr + ".log")

	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	valueString , _ := tengo.ToString(args[1])
	if _, err := f.WriteString(valueString); err != nil {
		log.Println(err)
	}

	return nil, nil
}

func torIntercall(args ...tengo.Object) (ret tengo.Object, err error) {

	// write to file
	f, err := os.OpenFile("intercall.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	for _, arg := range args {
		valueString , _ := tengo.ToString(arg)
		if _, err := f.WriteString(valueString + "\n"); err != nil {
			log.Println(err)
		}
	}

	return nil, nil
}

func torExtercall(args ...tengo.Object) (ret tengo.Object, err error) {
	serviceAddress , _ := tengo.ToString(args[0])
	commandString  , _ := tengo.ToString(args[1])
	port , _ := tengo.ToString(args[2])

	torAddress := "127.0.0.1:" + port

	go externalRequest(serviceAddress, torAddress, commandString)

	return nil, nil
}

func externalRequest(serviceAddress string, torAddress string , commandString string ) {
	dialer, err := proxy.SOCKS5("tcp", torAddress, nil, proxy.Direct)
	if err != nil {
		log.Println(err)
	}

	conn, err := dialer.Dial("tcp", serviceAddress)
	if err != nil {
		log.Println(err)
	} else {
		defer conn.Close()

		_, err = conn.Write([]byte("///\n\r" + commandString + "\n\r//within cells interlinkedd"))

		if err != nil {
			log.Println(err)
		}

		var fullText string
		var buffer [1]byte
		p := buffer[:]
		_, _ = conn.Read(p)

		for {
			n, err_interlinked := conn.Read(p)
			if n > 0 {
				fmt.Println( fullText )
				fullText = fullText + string(p[:n])
				fullText = strings.Replace(fullText, "  ", " ", -1)
			}

			if nil != err_interlinked {
				log.Println("Error:")
				log.Println(err_interlinked)
				_, err = conn.Write([]byte("//End on rocks!"))
				break
			}
		}
		time.Sleep(30 * time.Millisecond)
	}
}


func getTorArgs(args ...tengo.Object) ([]interface{}, error) {
	var printArgs []interface{}
	l := 0
	for _, arg := range args {
		s, _ := tengo.ToString(arg)
		slen := len(s)
		// make sure length does not exceed the limit
		if l+slen > tengo.MaxStringLen {
			return nil, tengo.ErrStringLimit
		}
		l += slen
		printArgs = append(printArgs, s)
	}
	return printArgs, nil
}
