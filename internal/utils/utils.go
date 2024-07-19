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

func ValidateLoginInput(username, password string) (string, string) {
	usernameErr, passwordErr := "", ""

	if username == "" {
		usernameErr = "Username Must be Filled"
	}

	if password == "" {
		passwordErr = "Password Must be Filled"
	}

	if len(username) != 0 && (len(username) < 5 || len(username) > 255) {
		usernameErr = "Username must be between 5 and 255 characters"
	}

	if len(password) != 0 && (len(password) < 5 || len(password) > 255) {
		passwordErr = "Password must be between 5 and 255 characters"
	}

	return usernameErr, passwordErr
}

func ValidateRegisterInput(username, password, password2, email string) (string, string, string, string) {
	usernameErr, passwordErr, password2Err,  emailErr := "", "", "", ""

	if username == "" {
		usernameErr = "Username must be filled"
	}

	if password == "" {
		passwordErr = "Password must be filled"
	}

	if password2 == "" {
		passwordErr = "Password Confirmation must be filled"
	}

	if password2 != "" && password != password2 {
		password2Err = "Password and Password Confirmation must be the same"
	}

	if email == "" {
		emailErr = "Email must be filled"
	}

	if len(username) != 0 && (len(username) < 5 || len(username) > 255) {
		usernameErr = "Username must be between 5 and 255 characters"
	}

	if len(password) != 0 && (len(password) < 5 || len(password) > 255) {
		passwordErr = "Password must be between 5 and 255 characters"
	}

	if len(password2) != 0 && (len(password2) < 5 || len(password2) > 255) {
		passwordErr = "Password must be between 5 and 255 characters"
	}

	if email != "" && !IsValidEmail(email) {
		emailErr = "Email must contain @ and end with .com or .co.id"
	}
	
	return usernameErr, passwordErr, password2Err, emailErr
}