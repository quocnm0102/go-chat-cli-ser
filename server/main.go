package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

var (
	clients []net.Conn

	connCh = make(chan net.Conn)
	msgCh = make(chan string)
	closeCh = make(chan net.Conn)
)

func main()  {
	server, err := net.Listen("tcp", ":3000")

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Fatal(err)
			}

			clients = append(clients, conn)
			fmt.Println("clients len: ", len(clients))
			connCh <- conn
		}
	}()

	for {
		select {
			case conn := <- connCh :
				go onMessage(conn)
			case msg := <- msgCh :
				fmt.Print(msg)
			case conn := <- closeCh :
				fmt.Println("close client ", conn)
				removeConn(conn)
		}
	}

}

func publish(conn net.Conn, msg string) {
	for _, client := range clients {
		if client == conn {
			continue
		}
		_, err := client.Write([]byte(msg))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func removeConn(conn net.Conn)  {
	for i, client := range clients {
		if client == conn {
			clients = append(clients[:i], clients[i+1:]...)
			return
		}
	}
	fmt.Println(clients)
}

func onMessage(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("onMessage err: ", err)
			break
		}
		msgCh <- msg
		publish(conn, msg)
	}

	closeCh <- conn
}
