package main

import (
	//"errors"
	"github.com/patrickmn/go-cache"
	"time"
)

var localCache *cache.Cache

func InitCache() {
	localCache = cache.New(5*time.Minute, 10*time.Minute)
}

func SetCache(rank string, univ University) bool {
	localCache.Set(rank, univ, 5*time.Minute)
	return true
}

func GetCache(rank string) (interface{}, bool, string) {
	var source string
	data, found := localCache.Get(rank)
	if found {
		source = "Cache"
	}
	return data, found, source
}
