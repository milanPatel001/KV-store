package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type Statement struct {
	Command string
	Args    []string
}

type Connection struct {
	Id               string
	IP               string
	TransactionQueue []Statement
	TransactionFlag  bool
}

type PlainCache struct {
	Mutex sync.Mutex
	Data  map[string]string
}

var connectionMap = make(map[string]*Connection)
var plainCache = PlainCache{Data: make(map[string]string)}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = ":9376"
	} else {
		PORT = ":" + PORT
	}

	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	fmt.Println("Server running on Port " + PORT)

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}

}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	defer c.Close()

	id, _ := GenerateRandomId(6)

	connObj := Connection{IP: c.RemoteAddr().String(), Id: id}
	connectionMap[c.RemoteAddr().String()] = &connObj
	buffer := make([]byte, 1024)

	for {
		n, err := c.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Println(connObj.IP + ": Client disconnected")
				return
			}
			log.Println("Error reading:", err)
			return
		}

		data := string(buffer[:n])
		input := strings.TrimSpace(string(data))

		if strings.ToUpper(input) == "EXIT" {
			log.Println(connObj.IP + ": Client disconnected")
			break
		}

		parts := strings.Fields(input)
		if len(parts) == 0 {
			c.Write([]byte("Nothing entered!!!"))
			continue
		}

		command := strings.ToUpper(parts[0])
		args := parts[1:]

		fmt.Print(command + " : ")
		fmt.Print(args)
		fmt.Print("\n")

		SwitchCases(command, args, &connObj, c)
	}

}
