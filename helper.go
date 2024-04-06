package ginx

import (
	"github.com/gin-gonic/gin"
	"path"
)

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}
func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

func lastHandler(chain gin.HandlersChain) gin.HandlerFunc {
	if len(chain) == 0 {
		return nil
	}
	return chain[len(chain)-1]
}

const logPrefix = "[GinX]"
