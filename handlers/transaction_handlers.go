package handlers

import (
	"fmt"
	"prac/utils"
	"strconv"
)

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

	CurrentCache.TransactionMutex.Lock()
	defer CurrentCache.TransactionMutex.Unlock()

	for _, statement := range statements {

		var previousItem CacheItem
		var keyExists bool

		if (statement.Command == "DEL" && len(statement.Args) != 0) || (statement.Command == "SET" && len(statement.Args) == 2) {
			CurrentCache.Mutex.Lock()
			previousItem, keyExists = CurrentCache.Data[statement.Args[0]]
			CurrentCache.Mutex.Unlock()
		}

		successMsg, err := CommandHandler(statement.Command, statement.Args)

		if err != nil {
			//fmt.Println(rollBackLog)
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
				rollbackStatement = Statement{"SET", []string{statement.Args[0], previousItem.Val, strconv.Itoa(int(previousItem.TTL))}}
			} else {
				rollbackStatement = Statement{"DEL", statement.Args}
			}
		} else if statement.Command == "DEL" {
			rollbackStatement = Statement{"SET", []string{statement.Args[0], previousItem.Val, strconv.Itoa(int(previousItem.TTL))}}
		}

		rollBackLog = utils.Prepend(rollBackLog, rollbackStatement)
	}

	return successMsgLog, nil
}
