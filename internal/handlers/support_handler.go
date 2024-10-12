package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AndreanDjabbar/ElectiVote/internal/factories"
	"github.com/AndreanDjabbar/ElectiVote/internal/middlewares"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
	"github.com/AndreanDjabbar/ElectiVote/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

var saweriaSecretKey = os.Getenv("SAWERIA_SECRET_KEY")

func initRedis() {
    redisAddr := os.Getenv("REDIS_URL")
    if strings.HasPrefix(redisAddr, "redis://") {
        redisAddr = strings.TrimPrefix(redisAddr, "redis://")
    }
    redisClient = redis.NewClient(&redis.Options{
        Addr: redisAddr,
        Password: "", 
        DB: 0,
    })
}

func ViewSaweriaPage(c *gin.Context) {
    initRedis()
    if !middlewares.IsLogged(c) {
        logger.Warn(
            "ViewSaweriaPage - User is not logged in",
            "Client IP", c.ClientIP(),
            "action", "redirecting to login page",
        )
        c.Redirect(http.StatusFound, "/electivote/login-page/")
        return
    }

    username := middlewares.GetUserData(c)

    status, transactionID := CheckDonationStatus()
    if status {
        logger.Info(
            "ViewSaweriaPage - Donation already processed",
            "Client IP", c.ClientIP(),
            "Username", username,
            "action", "redirecting to thanks page",
        )

        err := redisClient.Del(c.Request.Context(), "donationStatus", "transactionID").Err()
        if err != nil {
            logger.Error("ViewSaweriaPage - Error deleting donationStatus or transactionID from Redis", "error", err)
        }

        redirectURL := fmt.Sprintf("/electivote/thanks-page?transaction_id=%s", transactionID)
        c.Redirect(http.StatusFound, redirectURL)
        return
    }

    logger.Info(
        "ViewSaweriaPage - rendering saweria page",
        "Client IP", c.ClientIP(),
        "Username", username,
    )
    c.HTML(
        http.StatusOK,
        "saweriaPage.html",
        gin.H{
            "title": "Saweria",
        },
    )
}



func generateHMACSignature(secretKey string, data string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func validateSignature(c *gin.Context, secretKey string, data string) bool {
	sentSignature := c.GetHeader("saweria-callback-signature")
	expectedSignature := generateHMACSignature(secretKey, data)
	return hmac.Equal([]byte(expectedSignature), []byte(sentSignature))
}

func SaweriaWebhook(c *gin.Context) {
    initRedis()
	var payload struct {
		Version      string  `json:"version"`
		CreatedAt    string  `json:"created_at"`
		ID           string  `json:"id"`
		Type         string  `json:"type"`
		AmountRaw    float64 `json:"amount_raw"`
		DonatorName  string  `json:"donator_name"`
		DonatorEmail string  `json:"donator_email"`
        Message     string  `json:"message"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
    data := fmt.Sprintf("%s%s%.0f%s%s", payload.Version, payload.ID, payload.AmountRaw, payload.DonatorName, payload.DonatorEmail)

	if !validateSignature(c, saweriaSecretKey, data) {
        logger.Warn(
            "SaweriaWebhook - Invalid signature",
        )
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	if payload.Type == "donation" {
        logger.Info(
            "SaweriaWebhook - Donation received",
            "Client IP", c.ClientIP(),
            "transaction_id", payload.ID,
            "amount_raw", payload.AmountRaw,
            "donator", payload.DonatorName,
        )

        err := redisClient.Set(c.Request.Context(), "donationStatus", "completed", 0).Err()
        if err != nil {
            logger.Error("SaweriaWebhook - Failed to save donation", "error", err, "Client IP", c.ClientIP())
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process donation"})
            return
        }
        transactionID := payload.ID
        donatorName := payload.DonatorName
        donatorEmail := payload.DonatorEmail
        SaveTransaction(transactionID, 5*time.Minute)
        timeCreated := models.CustomTime{
            Time: time.Now(),
        }
        newSupport := factories.SupportFactory(donatorName, donatorEmail, transactionID, payload.AmountRaw, payload.Message, timeCreated)
        err = repositories.SaveSupport(newSupport)
        if err != nil {
            logger.Error(
                "SaweriaWebhook - failed to save support",
                "Client IP", c.ClientIP(),
                "error", err,
            )
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process donation"})
        }
		c.JSON(http.StatusOK, gin.H{"message": "Donation processed"})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported event type"})
	}
}

func CheckDonationStatus() (bool, string) {
    // Periksa status donasi
    status, err := redisClient.Get(context.Background(), "donationStatus").Result()
    if err == redis.Nil {
        logger.Warn("CheckDonationStatus - donationStatus key not found in Redis")
        return false, ""
    } else if err != nil {
        logger.Error("CheckDonationStatus - Failed to get donation status from Redis", "error", err)
        return false, ""
    }

    // Periksa transactionID
    transactionID, err := redisClient.Get(context.Background(), "transactionID").Result()
    if err == redis.Nil {
        logger.Warn("CheckDonationStatus - transactionID key not found in Redis")
        return false, ""
    } else if err != nil {
        logger.Error("CheckDonationStatus - Failed to get transactionID from Redis", "error", err)
        return false, ""
    }
    return status == "completed", transactionID
}


func ViewThanksPage(c *gin.Context) {
    if !middlewares.IsLogged(c) {
        logger.Warn(
            "ThanksPage - User is not logged in",
            "Client IP", c.ClientIP(),
            "action", "redirecting to login page",
        )
        c.Redirect(http.StatusFound, "/electivote/login-page/")
        return
    }
    transactionID := c.Query("transaction_id")
    if !IsValidTransaction(transactionID) {
        c.Redirect(
            http.StatusFound,
            "/electivote/about-us-page/",
        )
        return
    }

    context := gin.H{
        "title": "Thanks",
    }
    c.HTML(
        http.StatusOK,
        "thanksPage.html",
        context,
    )
}

func IsValidTransaction(transactionID string) bool {
    initRedis()
    val, err := redisClient.Get(ctx, transactionID).Result()
    if err == redis.Nil {
        return false
    } else if err != nil {
        return false
    }

    return val == "valid"
}

func MarkTransactionAsUsed(transactionID string) {
    redisClient.Del(ctx, transactionID)
}

func SaveTransaction(transactionID string, ttl time.Duration) {
    redisClient.Set(ctx, "transactionID", transactionID, ttl)
    redisClient.Set(ctx, transactionID, "valid", ttl)
}
