package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strings"
	"time"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
	"github.com/AndreanDjabbar/ElectiVote/internal/repositories"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var secretKey []byte = []byte(os.Getenv("SECRET_KEY"))
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateSecureToken(username string) (string, error) {
	timestamp := time.Now().Unix()
	data := fmt.Sprintf("%s:%d", username, timestamp)
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	signature := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s:%s", data, signature), nil
}

func ExtractUsername(token string) (string, error) {
	parts := strings.Split(token, ":")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid token format")
	}

	data := fmt.Sprintf("%s:%s", parts[0], parts[1])
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return "", fmt.Errorf("invalid token signature")
	}

	return parts[0], nil
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

func IsValidPhoneNumber(phone string) bool {
	isValid := regexp.MustCompile(`^[0-9]+$`).MatchString(phone)
	return isValid
}

func ValidateProfileInput(firstName, lastname, phone string, age uint) (string, string, string, string) {
	firstNameErr, lastNameErr, phoneErr, ageErr :=  "", "", "", ""

	if firstName != "" && (len(firstName) < 5 || len(firstName) > 255) {
		firstNameErr = "First Name must be between 5 and 255 characters"
	}

	if lastname != "" && (len(lastname) < 5 || len(lastname) > 255) {
		lastNameErr = "Last Name must be between 5 and 255 characters"
	}

	if age != 0 && (age <= 5 || age > 100) {
		ageErr = "Age must be between 5 and 100"
	}

	if phone != "" && !IsValidPhoneNumber(phone) {
		phoneErr = "Phone number must be a number"
	}

	return firstNameErr, lastNameErr, phoneErr, ageErr
}

func RenderError(c *gin.Context, statusCode int, errMsg string, source string) {
	context := gin.H{
		"title":  "Error",
		"error":  errMsg,
		"source": source,
	}
	c.HTML(
		statusCode,
		"error.html",
		context,
	)
}

func FormattedDob(nt models.NullTime) string {
	if nt.Valid {
		return nt.Time.Format("2006-01-02")
	}
	return ""
}

func voteCodeMaker() string {
    result := make([]byte, 6)
    for i := range result {
        num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
        if err != nil {
            return ""
        }
        result[i] = charset[num.Int64()]
    }
    return string(result)
}

func GenerateVoteCode() string {
	voteCode := voteCodeMaker()
	for !repositories.IsUniqueCode(voteCode) {
		voteCode = voteCodeMaker()
	}
	return voteCode
}