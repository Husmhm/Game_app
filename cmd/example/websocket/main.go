package main

import (
	"fmt"
	"gameApp/entity"
	"gameApp/pkg/protobufencoder"
	"log"
	"net"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func printDecodeNotification() {
	d := protobufencoder.EncodeNotification(entity.Notification{
		EventType: "ping",
		Payload:   "victory!",
	})
	fmt.Println(d)
}

func main() {
	printDecodeNotification()
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			panic(err)
		}
		done := make(chan bool)

		go readMessage(conn, done)

		<-done
		go func() {
			defer conn.Close()

			for {
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					panic(err)
				}
				err = wsutil.WriteServerMessage(conn, op, msg)
				if err != nil {
					panic(err)
				}
			}
		}()
	}))
}

func readMessage(conn net.Conn, done chan<- bool) {
	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			log.Println(err)
			done <- true
		}

		notif := protobufencoder.DecodeNotification(string(msg))

		fmt.Println("notif", notif)
		fmt.Println("OpCode:", op)

	}

}
