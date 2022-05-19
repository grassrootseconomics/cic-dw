package address

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/sha3"
)

func SarafuAddress(address string) (string, error) {
	l := len(address)

	if l < 40 || l > 42 {
		return "", fmt.Errorf("%s is not a valid eth address", address)
	}

	if len(address) == 42 {
		return strings.ToLower(address)[2:], nil
	}

	return strings.ToLower(address), nil
}

func Checksum(address string) string {
	address = strings.ToLower(address)
	address = strings.Replace(address, "0x", "", 1)

	sha := sha3.NewLegacyKeccak256()
	sha.Write([]byte(address))
	hash := sha.Sum(nil)
	hashstr := hex.EncodeToString(hash)
	result := []string{"0x"}
	for i, v := range address {
		res, _ := strconv.ParseInt(string(hashstr[i]), 16, 64)
		if res > 7 {
			result = append(result, strings.ToUpper(string(v)))
			continue
		}
		result = append(result, string(v))
	}

	return strings.Join(result, "")
}
