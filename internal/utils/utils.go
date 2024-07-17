package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

func GenerateSecureToken(username string) (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	data := username + ":" + time.Now().Format(time.RFC3339) + ":" + base64.URLEncoding.EncodeToString(randomBytes)

	mac := hmac.New(sha256.New, randomBytes)
	mac.Write([]byte(data))
	signature := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	return base64.URLEncoding.EncodeToString([]byte(data + ":" + signature)), nil
}