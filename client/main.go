package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func onMessage(conn net.Conn) {
	for {
		msgReader := bufio.NewReader(conn)
		msg, err := msgReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(msg)
	}
}
func main() {
	conn, er := net.Dial("tcp", ":3000")

	if er != nil {
		log.Fatal(er)
	}

	fmt.Println("Connect success!")

	nameReader := bufio.NewReader(os.Stdin)
	fmt.Print("Type your name: ")
	name, _ := nameReader.ReadString('\n')
	name = name[:len(name)-1]

	go onMessage(conn)

	for {
		msgReader := bufio.NewReader(os.Stdin)
		msg, _ := msgReader.ReadString('\n')
		msg = fmt.Sprintf("%s: %s\n", name, msg[:len(msg)-1])

		_, err := conn.Write([]byte(msg))
		if err != nil {
			log.Fatal(err)
		}
	}

	conn.Close()
}
