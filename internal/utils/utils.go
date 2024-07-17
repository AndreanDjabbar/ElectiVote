package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
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

func IsValidEmail(email string) bool {
	const emailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.(com|co\.id)$`
	re := regexp.MustCompile(emailPattern)
	return re.MatchString(email)
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
} 