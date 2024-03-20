/*
package main

import (
	"bufio"

	"log"
	"os"

	socketio_client "github.com/zhouhui8915/go-socket.io-client"
)

func main() {

	opts := &socketio_client.Options{
		//Transport:"polling",
		Transport:"websocket",
		Query:     make(map[string]string),
	}
	opts.Query["user"] = "user"
	opts.Query["pwd"] = "pass"
	uri := "http://localhost:3030"

	client, err := socketio_client.NewClient(uri, opts)
	if err != nil {
		log.Printf("NewClient error:%v\n", err)
		return
	}

	client.On("error", func() {
		log.Printf("on error\n")
	})
	client.On("connection", func() {
		log.Printf("on connect\n")
	})
	client.On("message", func(msg string) {
		log.Printf("on message:%v\n", msg)
	})
	client.On("disconnection", func() {
		log.Printf("on disconnect\n")
	})

	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		command := string(data)
		client.Emit("message", command)
		log.Printf("send message:%v\n", command)
	}
}
*/

/*
package main

import (
	gosocketio "github.com/VerticalOps/golang-socketio"
	transport "github.com/VerticalOps/golang-socketio/transport"
	"log"
	"runtime"
	"time"
)

type Channel struct {
	Channel string `json:"channel"`
}

type Message struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func sendJoin(c *gosocketio.Client) {
	log.Println("Acking /join")
	result, err := c.Ack("/join", Channel{"main"}, time.Second*5)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Ack result to /join: ", result)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	c, err := gosocketio.Dial(
		gosocketio.GetUrl("localhost", 3030, false),
		transport.GetDefaultWebsocketTransport())
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("/chat", func(h *gosocketio.Channel, args Message) {
		log.Println("--- Got chat message: ", args)
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Fatal("Disconnected")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected")
	})
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	go sendJoin(c)
	go sendJoin(c)
	go sendJoin(c)
	go sendJoin(c)
	go sendJoin(c)

	time.Sleep(60 * time.Second)
	c.Close()

	log.Println(" [x] Complete")
}
*/

/*
package main

import (
	socketio "github.com/googollee/go-socket.io"
)

func main() {
	// Simple client to talk to default-http example
	// uri := "http://127.0.0.1:8000"

	uri := "http://127.0.0.1:8000"

	_, err := socketio.NewClient(uri, nil)

	// Handle an incoming event
	// client.OnEvent("reply", func(s socketio.Conn, msg string) {
	// 	log.Println("Receive Message /reply: ", "reply", msg)
	// })

	// err = client.Connect()
	// if err != nil {
	// 	panic(err)
	// }

	// client.Emit("notice", "hello")

	// time.Sleep(1 * time.Second)
	// err = client.Close()
	// if err != nil {
	// 	panic(err)
	// }
}
*/

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type EventData struct {
	Event string  	`json:"event"`
    Data string	 	`json:"data"`
}

var addr = flag.String("addr", "localhost:3031", "http service address")

func SendIt(c *websocket.Conn, eventName string, sendData string) error {	
	jsonData, err := json.Marshal(EventData{eventName, sendData})
	if err != nil {
		return err
	}

	return c.WriteMessage(websocket.TextMessage, jsonData)
}

func SendClose(c *websocket.Conn) error {
	return c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			log.Printf("msgType: %d, recv: %s", messageType, message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	SendIt(c, "ping", "Hello, World!")

	Count := 1

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:			
			if (Count <= 3) {
				err := SendIt(c, "ws-chat", t.String())
				if err != nil {
					log.Println("write:", err)
					return
				}
			}
			Count++

			// err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			// if err != nil {
			//	log.Println("write:", err)
			// 	return
			// }
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			// err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			err := SendClose(c)
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}


////////////////////////////////////////////////
// https://github.com/graarh/golang-socketio

/*
package main

import (
	gosocketio "github.com/graarh/golang-socketio"
	transport "github.com/graarh/golang-socketio/transport"
	"log"
	"runtime"
	"sync"
	"time"
)

func doSomethingWith(c *gosocketio.Client, wg *sync.WaitGroup) {
    if res, err := c.Ack("join", "This is a client", time.Second*3); err != nil {
            log.Printf("error: %v", err)
    } else {
            log.Printf("result %q", res)
    }
    wg.Done()
}

func main() {

    runtime.GOMAXPROCS(runtime.NumCPU())

    c, err := gosocketio.Dial(
        gosocketio.GetUrl("localhost", 3031, false),
        transport.GetDefaultWebsocketTransport())
    if err != nil {
        log.Fatal(err)
    }

    err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
        log.Fatal("Disconnected")
    })
    if err != nil {
        log.Fatal(err)
    }

    err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
        log.Println("Connected")
    })
    if err != nil {
        log.Fatal(err)
    }

    wg := &sync.WaitGroup{}
    wg.Add(1)

	go doSomethingWith(c, wg)

	wg.Wait()
    log.Printf("Done")
}
*/

/*
package main

import (
	"encoding/json"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"log"
	"runtime"
	"sync"
	"time"
)

type Channel struct {
	Channel string `json:"channel"`
}

type Message struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type EventData struct {
	Event string  	`json:"event"`
    Data string	 	`json:"data"`
}


func sendJoin(c *gosocketio.Client) {
	log.Println("Acking /join")
	result, err := c.Ack("ws-chat", Channel{"ws-chat"}, time.Second*5)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Ack result to /join: ", result)
	}
}

func doSomethingWith(c *gosocketio.Client, wg *sync.WaitGroup) {
    if res, err := c.Ack("join", "This is a client", time.Second*3); err != nil {
        log.Printf("error: %v", err)
    } else {
        log.Printf("result %q", res)
    }
    wg.Done()
}


func SendIt(c *gosocketio.Client, eventName string, sendData string) error {	
	jsonData, err := json.Marshal(EventData{eventName, sendData})
	if err != nil {
		return err
	}

	return c.Emit(eventName, jsonData)
}


func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	c, err := gosocketio.Dial(
		gosocketio.GetUrl("localhost", 3031, false),
		transport.GetDefaultWebsocketTransport())
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("ws-chat", func(h *gosocketio.Channel, args Message) {
		log.Println("--- Got chat message: ", args)
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Fatal("Disconnected")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected")
	})
	if err != nil {
		log.Fatal(err)
	}

	SendIt(c, "ws-chat", "Hello, World!")

	// go sendJoin(c)

	// c.Emit("/join", Channel{"ws-chat"});

	// 데이터 전송
	// c.Emit("ws-chat", Message{1, "ws-chat", "Hello, World!"})

	// c.Emit("/", Message{1, "ws-chat", "Hello, World!"})
	
    // c.Emit("data", MyEventData{"data", "my data"})	

	// c.On("ws-chat", func(msg string) {
        // fmt.Println("message:", msg)
    // })

	// c.Emit("ws-chat", "hello world")
    

    select {}

	// time.Sleep(10 * time.Second)

	// go sendJoin(c)
	// go sendJoin(c)
	// go sendJoin(c)
	// go sendJoin(c)
	// go sendJoin(c)

	// time.Sleep(60 * time.Second)
	// c.Close()

	// log.Println(" [x] Complete")
}
*/