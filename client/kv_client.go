package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

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

		output := DeserializeOutput(string(buffer[:n]), &currentCacheNum)

		fmt.Println(output)
	}

}

// +COMMAND\r\nOUTPUT\r\n
func DeserializeOutput(s string, cacheNum *uint8) string {

	str := strings.Split(s, "\r\n")

	command := str[0]
	output := str[1]

	if command == "NUM" {
		val, _ := strconv.Atoi(output)

		*cacheNum = uint8(val)
	} else if command == "ERR" {
		output = "-" + output
	}

	return output
}

// func SerializeInput(input string) string {
// 	// COMMAND\r\nARG1\r\nARG2\r\n
// 	string
// }
