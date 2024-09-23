package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ADDR := os.Getenv("ADDR")

	if ADDR == "" {
		ADDR = "localhost:9376"
	}

	conn, err := net.Dial("tcp", ADDR)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("CONNECTED TO KV SERVER...")
	for {
		fmt.Print("\n$ ")

		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		_, err = conn.Write([]byte(input))
		if err != nil {
			log.Fatal(err)
		}

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(buffer[:n]))
	}

}
