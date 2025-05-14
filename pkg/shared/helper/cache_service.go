package helper

import (
	"github.com/patrickmn/go-cache"
)

var suggestionCache *cache.Cache

func InitSuggestionCache() {
	suggestionCache = cache.New(cache.NoExpiration, cache.NoExpiration)
}

func SetSuggestion(userID string, suggestionIDs []string) {
	suggestionCache.Set(userID, suggestionIDs, cache.NoExpiration)
}

func GetSuggestion(userID string) ([]string, bool) {
	if x, found := suggestionCache.Get(userID); found {
		return x.([]string), true
	}
	return nil, false
}
