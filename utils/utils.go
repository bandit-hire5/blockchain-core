package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"time"
)

func CalculateHash(index int64, prevBlockHash string, timestamp time.Time, data []byte) (string, error) {
	timeString := timestamp.String()
	str := []byte(string(index) + prevBlockHash + timeString + string(data))

	hasher := sha1.New()
	hasher.Write(str)

	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return sha, nil
}
