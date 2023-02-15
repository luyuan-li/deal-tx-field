package utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func BuildHex(bytes []byte) string {
	return strings.ToUpper(hex.EncodeToString(bytes))
}

func ConvertErr(height int64, txHash, errTag string, err error) error {
	return fmt.Errorf("%v-%v-%v-%v", err.Error(), errTag, height, txHash)
}

func CheckSkipErr(err error, tag string) bool {
	return strings.Contains(err.Error(), tag)
}
func GetErrTag(err error) string {
	slice := strings.Split(err.Error(), "-")
	if len(slice) == 4 {
		return slice[1]
	}
	return ""
}

func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func MarshalJsonIgnoreErr(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func UnMarshalJsonIgnoreErr(data string, v interface{}) {
	json.Unmarshal([]byte(data), &v)
}

func ConvStrToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func ConvStrToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

// Intn returns, as an int, a non-negative pseudo-random number in [0,n)
// from the default Source.
// It panics if n <= 0.
func RandInt(n int) int {
	rand.NewSource(time.Now().Unix())
	return rand.Intn(n)
}
