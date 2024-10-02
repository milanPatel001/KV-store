package handlers

import (
	"fmt"
	"prac/utils"
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

type CacheItem struct {
	Val       string
	CanExpire bool
	TTL       uint32
}

type Cache struct {
	Mutex            sync.Mutex
	TransactionMutex sync.Mutex
	Data             map[string]CacheItem
	SkipList         *utils.TTLSkipList
}

type CurrentSnapshot struct {
	DoneChannel chan int
	TimePeriod  uint32
}

var ConnectionMap = make(map[string]*Connection)
var SnapShotMap = make(map[uint8]CurrentSnapshot)

var Caches []Cache
var CurrentCache *Cache
var DefaultCacheNum uint8 // total number of caches
var DefaultSkipListMaxHeight uint8

func SetUpCaches(cacheNum uint8, skipListMaxHeight uint8) error {

	if cacheNum < 1 || skipListMaxHeight < 10 {
		return fmt.Errorf("Cache Number should be >1 and skipListMaxHeight should be >10")
	}

	DefaultCacheNum = cacheNum
	DefaultSkipListMaxHeight = skipListMaxHeight

	Caches = make([]Cache, DefaultCacheNum)

	for index := range Caches {
		Caches[index] = Cache{Data: make(map[string]CacheItem), SkipList: utils.CreateTTLSkipList(DefaultSkipListMaxHeight)}
	}

	CurrentCache = &Caches[0]

	return nil
}
