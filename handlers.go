package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"strings"
)

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

func TransactionHandler(command string, args []string, connectionObj *Connection) (string, error) {

	switch command {
	case "BEGIN":

		if connectionObj.TransactionFlag == true {
			return "", fmt.Errorf("Transaction has already started !!!")
		}

		connectionObj.TransactionFlag = true
		return ">> TRANSACTION BEGINS", nil

	case "DISCARD":
		if connectionObj.TransactionFlag == false {
			return "", fmt.Errorf("Start the transaction first and queue some commands to discard !!!")
		}

		connectionObj.TransactionFlag = false
		connectionObj.TransactionQueue = connectionObj.TransactionQueue[:0]
		return ">> DISCARDED", nil

	case "COMMIT":

		if connectionObj.TransactionFlag == false {
			return "", fmt.Errorf("Start the Transaction first using : BEGIN !!!")
		}

		successMsgLog, err := CommitHandler(connectionObj.TransactionQueue)

		connectionObj.TransactionFlag = false
		connectionObj.TransactionQueue = connectionObj.TransactionQueue[:0]

		if err != nil {
			return "", err
		}

		fmt.Println("{")
		for _, successMsg := range successMsgLog {
			fmt.Println(successMsg)
		}
		fmt.Println("}")

		return ">> COMMIT SUCCESS", nil
	}

	return "", fmt.Errorf("Unknown command !!!")
}

func CommitHandler(statements []Statement) ([]string, error) {
	rollBackLog := []Statement{}
	successMsgLog := []string{}

	for _, statement := range statements {

		var previousVal string
		var keyExists bool

		if (statement.Command == "DEL" && len(statement.Args) != 0) || (statement.Command == "SET" && len(statement.Args) == 2) {
			plainCache.Mutex.Lock()
			previousVal, keyExists = plainCache.Data[statement.Args[0]]
			plainCache.Mutex.Unlock()
		}

		successMsg, err := CommandHandler(statement.Command, statement.Args)

		if err != nil {
			for _, st := range rollBackLog {
				CommandHandler(st.Command, st.Args)
			}

			return nil, err
		}

		successMsgLog = append(successMsgLog, successMsg)

		var rollbackStatement Statement

		// TODO : expand this properly for other commands
		if statement.Command == "SET" {
			if keyExists {
				rollbackStatement = Statement{"SET", []string{statement.Args[0], previousVal}}
			} else {
				rollbackStatement = Statement{"DEL", statement.Args}
			}
		} else if statement.Command == "DEL" {
			rollbackStatement = Statement{"SET", []string{statement.Args[0], previousVal}}
		}

		rollBackLog = append(rollBackLog, rollbackStatement)
	}

	return successMsgLog, nil
}

func DelHandler(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("DEL : Missing Key")
	}

	plainCache.Mutex.Lock()
	defer plainCache.Mutex.Unlock()

	if _, ok := plainCache.Data[args[0]]; !ok {
		return fmt.Errorf("DEL %s : Key doesn't exist !!!", args[0])
	}

	delete(plainCache.Data, args[0])

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
	plainCache.Mutex.Lock()
	plainCache.Data[args[0]] = value
	plainCache.Mutex.Unlock()

	return nil

}

func GetHandler(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("GET : Missing Key")
	}

	plainCache.Mutex.Lock()
	val, ok := plainCache.Data[args[0]]
	plainCache.Mutex.Unlock()

	if !ok {
		return "", fmt.Errorf("GET %v: Key doesn't exist!!!", args[0])
	}

	return val, nil
}

func MakeMap[K any]() map[string]K {
	return make(map[string]K)
}

func GenerateRandomId(length int) (string, error) {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}
