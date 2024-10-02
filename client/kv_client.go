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

var CommandsWithRequiredArgs []string = []string{"SET", "DEL", "GET", "NUM"}

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

		input, err := SerializeInput(scanner.Text())

		if err != nil {
			fmt.Println(err)
			continue
		}

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
		output = "- " + output
	}

	return output
}

func SerializeInput(input string) (string, error) {

	var output string
	arr := strings.Fields(strings.TrimSpace(input))

	if len(arr) == 0 {
		return "", fmt.Errorf(">> Nothing Entered !!!")
	}

	firstArg := strings.ToUpper(arr[0])

	if len(arr) == 1 {
		for _, c := range CommandsWithRequiredArgs {
			if c == firstArg {
				return "", fmt.Errorf("- %v : Missing Arguments !!!", firstArg)
			}
		}
	}

	output += firstArg + "\r\n"

	for i := 1; i < len(arr); i++ {
		output += arr[i] + "\r\n"
	}

	return output, nil
}
