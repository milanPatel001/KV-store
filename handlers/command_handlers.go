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
		conn.Write([]byte(utils.SerializeOutput("TR", ">> QUEUED")))
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
		conn.Write([]byte(utils.SerializeOutput("ERR", err.Error())))
	} else {
		conn.Write([]byte(utils.SerializeOutput(command, successMsg)))
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

		return fmt.Sprintf(">> %v", val), nil

	case "DEL":
		if err := DelHandler(args); err != nil {
			return "", err
		}

		return ">> SUCCESS", nil

	case "BF_CREATE":
		if err := BloomFilterCreationHandler(args); err != nil {
			return "", err
		}

		return ">> SUCCESS", nil

	case "BF_ADD":
		if err := BloomFilterAddHandler(args); err != nil {
			return "", err
		}

		return ">> SUCCESS", nil

	case "BF_EXISTS":
		val, err := BloomFilterExistsHandler(args)

		if err != nil {
			return "", err
		}

		return fmt.Sprintf(">> %v", val), nil

	case "NUM":
		num, err := SetCurrentCacheHandler(args)
		if err != nil {
			return "", err
		}

		return strconv.Itoa(num), nil

	case "SAVE":
		if err := SaveCacheHandler(args); err != nil {
			return "", err
		}

		return ">> SUCCESS", nil

	case "RETAIN":
		if err := RetainCacheHandler(args); err != nil {
			return "", err
		}

		return ">> SUCCESS", nil

	case "HALT":
		if err := StopSnapshot(args); err != nil {
			return "", err
		}

		return ">> SUCCESS", nil
	}

	return "", fmt.Errorf("Unknown command !!!")
}

func SaveCacheHandler(args []string) error {
	// SAVE [cacheIndex] [time]

	// NOTE : serialize input from client and put default value of current Cache for SAVE if only "SAVE" is entered by client.

	// CASE no time -> save cache[cacheIndex] in dump.gob file
	// CASE time -> save cache[cacheIndex] in "snapshot+cachIndex".gob file periodically (time in seconds)

	if len(args) == 0 {
		return utils.StoreCacheGobEncoded("dump", CurrentCache.Data)
	}

	num, err := strconv.Atoi(args[0])

	if err != nil {
		return err
	}

	if num < 0 || num >= int(DefaultCacheNum) {
		return fmt.Errorf("Cache index should lie in the range of [0, %v]", DefaultCacheNum)
	}

	var period int

	if len(args) > 1 {
		period, err = strconv.Atoi(args[1])

		if err != nil {
			return err
		}
	}

	var fileName string

	// Time of atleast 60 seconds is required to be considered for periodic snapshots
	if period <= 60 {
		fileName = fmt.Sprintf("dump_%v", num)
		return utils.StoreCacheGobEncoded(fileName, Caches[num].Data)
	}

	currentTime := strconv.Itoa(int(time.Now().Unix()))
	fileName = fmt.Sprintf("snapshot_%v_%v", currentTime, num)

	return SetSnapshots(fileName, uint8(num), uint32(period))
}

/*
RETAIN [fileName] (fileName default : dump.gob)
Overwrites current cache and skiplist
*/
func RetainCacheHandler(args []string) error {

	var fileName string = "dump"

	if len(args) > 0 {
		fileName = args[0]
	}

	m, err := utils.DecodeGobFile[string, CacheItem](fileName)

	if err != nil {
		fmt.Println(err)
		return err
	}

	CurrentCache = &Cache{Data: m, SkipList: utils.CreateTTLSkipList(48)}

	CurrentCache.Mutex.Lock()
	defer CurrentCache.Mutex.Unlock()

	t := uint32(time.Now().Unix())

	for k, v := range m {
		if v.TTL > t {
			CurrentCache.SkipList.Insert(k, v.TTL)
		}
	}

	return nil

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

func SetSnapshots(fileName string, cacheIndex uint8, t uint32) error {
	_, exists := SnapShotMap[cacheIndex]

	if exists {
		return fmt.Errorf("Snapshot for current index already running. Use HALT [cacheIndex] to stop snapshotting and then create new one!!!")
	}

	SnapShotMap[cacheIndex] = CurrentSnapshot{DoneChannel: make(chan int), TimePeriod: t}

	go runSnapShot(fileName, cacheIndex, t, SnapShotMap[cacheIndex].DoneChannel)

	return nil
}

func runSnapShot(fileName string, cacheIndex uint8, t uint32, doneChannel <-chan int) {
	for {
		select {
		case <-time.Tick(time.Second * time.Duration(t)):
			fmt.Printf("Snapshotted cacheIndex: %v!!!", cacheIndex)
			utils.StoreCacheGobEncoded(fileName, Caches[cacheIndex].Data)
		case <-doneChannel:
			fmt.Printf("Snapshotting stopped for cacheIndex: %v!!!", cacheIndex)
			return
		}
	}
}

func StopSnapshot(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("HALT: Missing cacheIndex")
	}

	cacheIndex, err := strconv.Atoi(args[0])

	if err != nil {
		return err
	}

	if cacheIndex < 0 || cacheIndex > int(DefaultCacheNum) {
		return fmt.Errorf("Cache index should lie in the range of [0, %v]", DefaultCacheNum)
	}

	snap, exists := SnapShotMap[uint8(cacheIndex)]

	if !exists {
		return fmt.Errorf("Snapshot is not set for cacheIndex: %v", cacheIndex)

	}

	snap.DoneChannel <- 1
	delete(SnapShotMap, uint8(cacheIndex))

	return nil
}

// BF_CREATE name [error_rate] [capacity] [SCALABLE -> T/F]
func BloomFilterCreationHandler(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("BF_CREATE : Missing name of bloom filter")
	}

	errorRate := 0.01
	cap := 1000
	key := args[0]
	scalable := false

	var err error

	if len(args) > 1 {
		errorRate, err = strconv.ParseFloat(args[1], 32)
		if err != nil {
			return err
		}
	}

	if len(args) > 2 {
		cap, err = strconv.Atoi(args[2])
		if err != nil {
			return err
		}
	}

	if len(args) > 3 {
		if args[3] == "T" || args[3] == "TRUE" {
			scalable = true
		}
	}

	if scalable {
		BloomFilterMap[key] = utils.CreateAdaptiveBloomFilter(uint(cap), errorRate)
	} else {
		BloomFilterMap[key] = utils.CreateBloomFilter(uint(cap), errorRate)
	}

	return nil

}

// BF_ADD name key
func BloomFilterAddHandler(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("BF_ADD : Missing name of the bloom filter and Value")
	}

	if len(args) == 1 {
		return fmt.Errorf("BF_ADD : Missing Value")
	}

	val, exists := BloomFilterMap[args[0]]

	if !exists {
		return fmt.Errorf("BF_ADD : Wrong name of the bloom filter")
	}

	return val.Set(args[1])

}

// BF_EXISTS name key
func BloomFilterExistsHandler(args []string) (bool, error) {
	if len(args) == 0 {
		return false, fmt.Errorf("BF_EXISTS : Missing name of the bloom filter and Value")
	}

	if len(args) == 1 {
		return false, fmt.Errorf("BF_EXISTS : Missing Value")
	}

	val, exists := BloomFilterMap[args[0]]

	if !exists {
		return false, fmt.Errorf("BF_EXISTS : Wrong name of the bloom filter")
	}

	return val.DoesExist(args[1]), nil
}
