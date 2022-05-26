package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// token:use sha256 encrypt userID-timeStamp
func GenerateToken(id int64) string {
	ts := fmt.Sprintf("%08x-%08x", id, time.Now().Unix())
	h := sha256.New()
	h.Write([]byte(ts))
	return hex.EncodeToString(h.Sum(nil))
}
