package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/AndreanDjabbar/ElectiVote/internal/factories"
	"github.com/AndreanDjabbar/ElectiVote/internal/middlewares"
	"github.com/AndreanDjabbar/ElectiVote/internal/repositories"
	"github.com/AndreanDjabbar/ElectiVote/internal/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ViewRegisterPage(c *gin.Context) {
	if middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}
	siteKey := os.Getenv("RECAPTCHA_SITE_KEY")
	context := gin.H {
		"title": "Register",
		"siteKey": siteKey,
	}
	c.HTML(
		http.StatusOK,
		"register.html",
		context,
	)
}

func RegisterPage(c *gin.Context) {
	if middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")
	password2 := c.PostForm("password2")
	email := c.PostForm("email")
	siteKey := os.Getenv("RECAPTCHA_SITE_KEY")

	usernameErr, passwordErr, password2Err, emailErr := utils.ValidateRegisterInput(username, password, password2, email)
	captchaErr := ""

	emailDomain := utils.GetEmailDomain(email)
	emailProvider := utils.GetEmailProvider(emailDomain)
	otp, err := utils.GenerateOTP()
	if err != nil {
		emailErr = "Failed to generate OTP"
	}

	if !utils.IsValidReCAPTCHA(c) {
		captchaErr = "Invalid ReCAPTCHA"
	}

	if usernameErr == "" && passwordErr == "" && password2Err == "" && emailErr == "" && captchaErr == "" {
		hashedPassword, hashErr := utils.HashPassword(password)
		if hashErr != nil {
			utils.RenderError(c, http.StatusInternalServerError, hashErr.Error(), "/electivote/register-page/")
			return
		}

		subject := "ElectiVote Email Verification"
		body := `
		<html>
		<head>
			<style>
				.container {
					font-family: Arial, sans-serif;
					max-width: 600px;
					margin: auto;
					padding: 20px;
					border: 1px solid #ddd;
					border-radius: 10px;
					box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
				}
				.otp-code {
					font-size: 24px;
					font-weight: bold;
					color: #000000;
				}
				.note {
					font-size: 14px;
					color: #555555;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<p>Hello,</p>
				<p>Your OTP Code is: <span class="otp-code">` + otp + `</span></p>
				<p class="note">Please use this code to verify your email address. <strong>Note:</strong> The OTP is valid for <strong>5 minutes</strong> from the time it was generated.</p>
				<p>If you did not request this verification, please ignore this email.</p>
				<p>Thank you!</p>
			</div>
		</body>
		</html>
		`
		go func() {
			err = utils.SendEmail(email, emailProvider, body, subject)
			if err != nil {
				fmt.Println("Failed to send email:", err)
			}
		}()

		middlewares.SetRegisterSession(c, username, email, hashedPassword, otp)
		c.Redirect(
			http.StatusFound,
			"/electivote/email-verification-page/",
		)
		return
	}

	context := gin.H{
		"title":       "Register",
		"usernameErr": usernameErr,
		"passwordErr": passwordErr,
		"password2Err": password2Err,
		"emailErr":    emailErr,
		"captchaErr":  captchaErr,
		"username":    username,
		"password":    password,
		"password2":   password2,
		"siteKey":     siteKey,
		"email":       email,	
	}
	c.HTML(
		http.StatusOK,
		"register.html",
		context,
	)
}

func registerUser(c *gin.Context, username string, password string, email string) {
	newUser := factories.CreateUser(username, password, email, "user")

	_, err := repositories.RegisterUser(newUser)
	if err != nil {
		utils.RenderError(c, http.StatusInternalServerError, err.Error(), "/electivote/register-page/")
		return
	}

	userID, err := repositories.GetUserIdByUsername(username)
	if err != nil {
		utils.RenderError(c, http.StatusInternalServerError, err.Error(), "/electivote/register-page/")
		return
	}

	newProfile := factories.CreateFirstProfile(userID)
	_, err = repositories.CreateProfile(newProfile)
	if err != nil {
		utils.RenderError(c, http.StatusInternalServerError, err.Error(), "/electivote/register-page/")
		return
	}
}

func ViewVerifyEmailPage(c *gin.Context) {
	if middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	session := sessions.Default(c)
	otp := session.Get("otp")

	if otp == nil {
		c.Redirect(
			http.StatusFound,
			"/electivote/register-page/",
		)
		return
	}

	context := gin.H {
		"title": "Verify Email",
	}
	c.HTML(
		http.StatusOK,
		"verifyEmail.html",
		context,
	)
}

func VerifyEmailPage(c *gin.Context) {
	if middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}
	session := sessions.Default(c)
	username := session.Get("username")
	email := session.Get("email")
	password := session.Get("password")
	otp := session.Get("otp")
	otpCreatedAt := session.Get("created_at")
	otpInput := c.PostForm("otp")
	expirationTime := 5 * 60

	if otp == nil {
		c.Redirect(
			http.StatusFound,
			"/electivote/register-page/",
		)
		return
	}

	otpErr := ""
	if otpInput == "" {
		otpErr = "OTP must be filled"
	}

	if otpInput != otp {
		otpErr = "Invalid OTP"
	}
	
	if otpInput != "" && (len(otpInput) != 6) {
		otpErr = "OTP must be 6 characters"
	}

	if time.Now().Unix() > otpCreatedAt.(int64) + int64(expirationTime) {
		middlewares.DeleteRegisterSession(c)
		c.HTML(
			http.StatusOK,
			"message.html",
			nil,
		)
		return
	}

	if otpErr != "" {
		context := gin.H{
			"title": "Verify Email",
			"otpErr": otpErr,
		}
		c.HTML(
			http.StatusOK,
			"verifyEmail.html",
			context,
		)
		return
	}
	registerUser(c, username.(string), password.(string), email.(string))
	middlewares.DeleteRegisterSession(c)
	c.Redirect(
		http.StatusFound,
		"/electivote/login-page/",
	)
}