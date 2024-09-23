package handlers

import (
	"fmt"
	"net"
	"strings"
	"sync"
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

type Cache struct {
	Mutex            sync.Mutex
	TransactionMutex sync.Mutex
	Data             map[string]string
}

var ConnectionMap = make(map[string]*Connection)
var PlainCache = Cache{Data: make(map[string]string)}

func SwitchCases(command string, args []string, connectionObj *Connection, conn net.Conn) {

	inTransaction := connectionObj.TransactionFlag

	if inTransaction && command != "COMMIT" && command != "DISCARD" && command != "BEGIN" {
		connectionObj.TransactionQueue = append(connectionObj.TransactionQueue, Statement{command, args})
		conn.Write([]byte(">> QUEUED"))
		return
	}

	var err error
	var successMsg string

	if command == "BEGIN" || command == "COMMIT" || command == "DISCARD" {
		successMsg, err = TransactionHandler(command, args, connectionObj)
	} else {
		successMsg, err = CommandHandler(command, args)
	}

	if err != nil {
		conn.Write([]byte(err.Error()))
	} else {
		conn.Write([]byte(successMsg))
	}
}

func CommandHandler(command string, args []string) (string, error) {
	switch command {
	case "SET":

		if err := SetHandler(args); err != nil {
			return "", err
		}

		return ">> SUCCESS", nil

	case "GET":
		val, err := GetHandler(args)

		if err != nil {
			return "", err
		}

		return ">> " + val, nil

	case "DEL":
		err := DelHandler(args)
		if err != nil {
			return "", err
		}

		return ">> SUCCESS", nil
	}

	return "", fmt.Errorf("Unknown command !!!")
}

func DelHandler(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("DEL : Missing Key")
	}

	PlainCache.Mutex.Lock()
	defer PlainCache.Mutex.Unlock()

	if _, ok := PlainCache.Data[args[0]]; !ok {
		return fmt.Errorf("DEL %s : Key doesn't exist !!!", args[0])
	}

	delete(PlainCache.Data, args[0])

	return nil
}

func SetHandler(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("SET : Missing Key and Value")
	}

	if len(args) == 1 {
		return fmt.Errorf("SET %s: Add value as well !!!", args[0])
	}

	var value string = args[1]

	if args[1][0] == '"' {
		for i := 2; i < len(args); i++ {
			value += " " + args[i]
			if strings.ContainsAny(args[i], "\"") {
				break
			}
		}
	}
	PlainCache.Mutex.Lock()
	PlainCache.Data[args[0]] = value
	PlainCache.Mutex.Unlock()

	return nil
}

func GetHandler(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("GET : Missing Key")
	}

	PlainCache.Mutex.Lock()
	val, ok := PlainCache.Data[args[0]]
	PlainCache.Mutex.Unlock()

	if !ok {
		return "", fmt.Errorf("GET %v: Key doesn't exist!!!", args[0])
	}

	return val, nil
}
