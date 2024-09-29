package handlers

import (
	"fmt"
	"math"
	"net"
	"prac/utils"
	"strconv"
	"strings"
	"time"
)

func SwitchCases(command string, args []string, connectionObj *Connection, conn net.Conn) {

	inTransaction := connectionObj.TransactionFlag

	if inTransaction && command != "COMMIT" && command != "DISCARD" && command != "BEGIN" {
		connectionObj.TransactionQueue = append(connectionObj.TransactionQueue, Statement{command, args})
		conn.Write([]byte(utils.SerializeInput("TR", ">> QUEUED")))
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
		conn.Write([]byte(utils.SerializeInput("ERR", err.Error())))
	} else {
		conn.Write([]byte(utils.SerializeInput(command, successMsg)))
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

		if err := DelHandler(args); err != nil {
			return "", err
		}

		return ">> SUCCESS", nil

	case "NUM":
		num, err := SetCurrentCacheHandler(args)
		if err != nil {
			return "", err
		}

		return strconv.Itoa(num), nil
	}

	return "", fmt.Errorf("Unknown command !!!")
}

func DelHandler(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("DEL : Missing Key")
	}

	CurrentCache.Mutex.Lock()
	defer CurrentCache.Mutex.Unlock()

	item, exist := CurrentCache.Data[args[0]]

	if !exist {
		return fmt.Errorf("DEL %s : Key doesn't exist !!!", args[0])
	}

	delete(CurrentCache.Data, args[0])

	CurrentCache.SkipList.Delete(args[0], item.TTL)

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
	var ttl uint32
	if len(args) > 2 {
		val, err := strconv.Atoi(args[2])

		if err != nil {
			ttl = 0
		} else {
			ttl = uint32(val)
		}
	}

	if args[1][0] == '"' {
		for i := 2; i < len(args); i++ {
			value += " " + args[i]
			if strings.ContainsAny(args[i], "\"") {
				break
			}
		}
	}

	canExpire := false
	var expiry uint32

	if ttl > 0 {
		canExpire = true

		now := uint32(time.Now().Unix())

		if ttl > math.MaxInt32-now-1 {
			expiry = math.MaxInt32 - 1
		} else {
			expiry = now + ttl
		}

	}

	CurrentCache.Mutex.Lock()
	item, exist := CurrentCache.Data[args[0]]
	if exist {
		// NOTE: Seperate command for changing ttl, so don't bother with it here
		CurrentCache.Data[args[0]] = CacheItem{Val: value, CanExpire: item.CanExpire, TTL: item.TTL}
	} else {
		CurrentCache.Data[args[0]] = CacheItem{Val: value, CanExpire: canExpire, TTL: expiry}

		if canExpire {
			CurrentCache.SkipList.Insert(args[0], expiry)
		}
	}

	CurrentCache.Mutex.Unlock()

	return nil
}

func GetHandler(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("GET : Missing Key")
	}

	CurrentCache.Mutex.Lock()
	item, exist := CurrentCache.Data[args[0]]
	CurrentCache.Mutex.Unlock()

	if !exist {
		return "", fmt.Errorf("GET %v: Key doesn't exist!!!", args[0])
	}

	return item.Val, nil
}

func SetCurrentCacheHandler(args []string) (int, error) {
	if len(args) == 0 {
		return -1, fmt.Errorf("No number provided !!!")
	}

	num, err := strconv.Atoi(args[0])

	if err != nil {
		return -1, err
	}

	if num < 0 || num >= int(DefaultCacheNum) {
		return -1, fmt.Errorf("Cache number must be in range of [0, %v].", DefaultCacheNum-1)
	}

	CurrentCache = &Caches[num]
	return num, nil
}
