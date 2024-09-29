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
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
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

	var currentCacheNum uint8 = 0
	fmt.Println("CONNECTED TO KV SERVER...")
	for {
		fmt.Printf("\n[%v]>>> ", currentCacheNum)

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
