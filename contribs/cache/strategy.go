package cache

import (
	gincache "github.com/chenyahui/gin-cache"
	"github.com/gin-gonic/gin"
	"net/url"
	"sort"
	"strings"
)

// CacheByUri generates cache key by request uri, if ignoreOrder is true, all query parameters will be sorted.
func CacheByUri(ignoreOrder bool) gincache.GetCacheStrategyByRequest {
	return func(c *gin.Context) (bool, gincache.Strategy) {
		if !ignoreOrder {
			return true, gincache.Strategy{CacheKey: c.Request.RequestURI}
		}

		orderUri, err := getRequestUriIgnoreQueryOrder(c.Request.RequestURI)
		if err != nil {
			c.Error(err)
		}
		return true, gincache.Strategy{CacheKey: orderUri}
	}
}

func getRequestUriIgnoreQueryOrder(requestURI string) (string, error) {
	parsedUrl, err := url.ParseRequestURI(requestURI)
	if err != nil {
		return "", err
	}

	values := parsedUrl.Query()

	if len(values) == 0 {
		return requestURI, nil
	}

	queryKeys := make([]string, 0, len(values))
	for queryKey := range values {
		queryKeys = append(queryKeys, queryKey)
	}
	sort.Strings(queryKeys)

	queryVals := make([]string, 0, len(values))
	for _, queryKey := range queryKeys {
		sort.Strings(values[queryKey])
		for _, val := range values[queryKey] {
			queryVals = append(queryVals, queryKey+"="+val)
		}
	}

	return parsedUrl.Path + "?" + strings.Join(queryVals, "&"), nil
}

// CacheByPath generates cache key by request path
func CacheByPath() gincache.GetCacheStrategyByRequest {
	return func(c *gin.Context) (bool, gincache.Strategy) {
		return true, gincache.Strategy{
			CacheKey: c.Request.URL.Path,
		}
	}
}
