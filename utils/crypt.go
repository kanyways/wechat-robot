package utils

import (
	"sort"
	"crypto/sha1"
	"io"
	"strings"
	"fmt"
	"github.com/gin-gonic/gin"
)

// 对请求的参数做SHA1校验
func makeSignature(token, timestamp, nonce string) string {
	sl := []string{token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}

// 检查url中的参数
func ValidateUrl(ctx *gin.Context, token string) bool {
	timestamp := ctx.Param("timestamp")
	nonce := ctx.Param("nonce")
	signatureGen := makeSignature(token, timestamp, nonce)
	signatureIn := ctx.Param("signature")
	if signatureGen != signatureIn {
		return false
	}
	return true
}

func ValidateGetUrl(ctx *gin.Context, token string) bool {
	timestamp := ctx.Query("timestamp")
	nonce := ctx.Query("nonce")
	signatureGen := makeSignature(token, timestamp, nonce)
	signatureIn := ctx.Query("signature")
	if signatureGen != signatureIn {
		return false
	}
	return true
}
