package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"prac/handlers"
	"prac/utils"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = ":9376"
	} else {
		PORT = ":" + PORT
	}

	ctx, cancel := context.WithCancel(context.Background())

	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		l.Close()
		cancel()
	}()

	fmt.Println("Server running on Port " + PORT)

	go handleSkipListExpiry(ctx)

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

	id, _ := utils.GenerateRandomId(6)

	connObj := handlers.Connection{IP: c.RemoteAddr().String(), Id: id}
	handlers.ConnectionMap[c.RemoteAddr().String()] = &connObj

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

		handlers.SwitchCases(command, args, &connObj, c)
	}
}

func handleSkipListExpiry(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			deletedKeys := handlers.PlainCache.SkipList.DeleteExpiredKeys()

			handlers.PlainCache.Mutex.Lock()

			for _, key := range deletedKeys {
				delete(handlers.PlainCache.Data, key)
			}

			handlers.PlainCache.Mutex.Unlock()

		case <-ctx.Done():
			return
		}
	}
}
