package common

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
)

func IntToHex(i interface{}) string {
	hex := fmt.Sprintf("%x", i)
	if strings.HasPrefix(hex, "0x") {
		return hex
	}
	return "0x" + hex
}

func HexToUnit64(hex string) (number uint64, err error) {
	if len(hex) < 2 {
		return 0, nil
	}
	return strconv.ParseUint(hex[2:], 16, 64)
}

func HexToInt64(hex string) (number int64, err error) {
	return strconv.ParseInt(hex[2:], 16, 64)
}

func HexToBigInt(hex string) *big.Int {
	if strings.HasPrefix(hex, "0x") {
		hex = hex[2:]
	}
	n := new(big.Int)
	n, _ = n.SetString(hex[:], 16)
	return n
}

// GetRandString 随机生成N位字符串
func GetRandString(n int) string {
	mainBuff := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, mainBuff)
	if err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
	return hex.EncodeToString(mainBuff)[:n]
}
