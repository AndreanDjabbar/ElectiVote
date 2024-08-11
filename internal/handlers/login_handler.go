package handlers

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"sync"

	"github.com/AndreanDjabbar/ElectiVote/internal/middlewares"
	"github.com/AndreanDjabbar/ElectiVote/internal/repositories"
	"github.com/AndreanDjabbar/ElectiVote/internal/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func ViewLoginPage(c *gin.Context) {
	if middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	context := gin.H {
		"title": "Login",
	}
	c.HTML(
		http.StatusOK,
		"login.html",
		context,
	)
}

func LoginPage(c *gin.Context) {
	if middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")
	remember := c.PostForm("remember")

	usernameErr, passwordErr := utils.ValidateLoginInput(username, password)
	var usernameCheckErr, passwordCheckErr error
	var wg sync.WaitGroup
	var mu sync.Mutex

	if usernameErr == "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := repositories.GetUserByUsername(username); err != nil {
				mu.Lock()
				usernameCheckErr = err
				mu.Unlock()
			}
		}() 
	}

	if passwordErr == "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := repositories.CheckPasswordByUSername(username, password); err != nil {
				mu.Lock()
				passwordCheckErr = err
				mu.Unlock()
			}
		}()
	
	}
	wg.Wait()

	if usernameCheckErr != nil {
		usernameErr = "Username not found"
	}

	if passwordCheckErr != nil {
		passwordErr = "Password is incorrect"
	}

	if usernameErr == "" && passwordErr == "" {
		if remember == "on" {
			middlewares.SetCookies(c, username)
		} else {
			middlewares.SetSession(c, username)
		}
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	context := gin.H {
		"title": "Login",
		"usernameErr": usernameErr,
		"passwordErr": passwordErr,
		"username": username,
		"password": password,
	}
	c.HTML(
		http.StatusOK,
		"login.html",
		context,
	)
}

func LogoutPage(c *gin.Context) {
	middlewares.DeleteSession(c)
	middlewares.DeleteCookie(c)
	c.Redirect(
		http.StatusFound,
		"/electivote/login-page/",
	)
}

func ViewForgotPasswordPage(c *gin.Context) {
	if middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	context := gin.H {
		"title": "Forgot Password",
	}
	c.HTML(
		http.StatusOK,
		"forgotPassword.html",
		context,
	)
}

func sendResetPasswordEmail(email, token string) error {
	from := os.Getenv("EMAIL_SENDER")
	password := os.Getenv("EMAIL_PASSWORD")
	fmt.Println(from, password)
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	resetURL := fmt.Sprintf("http://localhost:8080/electivote/reset-password-page/%s", token)

	subject := "Subject: Reset your password\n"
	body := fmt.Sprintf("To reset your password, please click the following link:\n\n%s\n", resetURL)
	message := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", smtpHost, smtpPort),
		auth,
		from,
		[]string{email},
		message,
	)
	if err != nil {
		return err
	}
	return nil
}

func verifyResetToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return utils.SecretKey, nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("Invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("Invalid token claims")
	}

	userEmail := claims["email"].(string)
	return userEmail, nil
}

func ForgotPasswordPage(c *gin.Context) {
	if middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}
	email := c.PostForm("email")
	emailErr := ""
	if len(email) == 0 {
		emailErr = "Email must be filled"
	}

	if !utils.IsValidEmail(email) {
		emailErr = "Email is not valid"
	}

	_, err := repositories.GetUserByEmail(email)
	if err != nil {
		emailErr = "Email not found"
	}

	if emailErr != "" {
		context := gin.H {
			"title": "Forgot Password",
			"emailErr": emailErr,
			"email": email,
		}
		c.HTML(
			http.StatusOK,
			"forgotPassword.html",
			context,
		)
		return
	}

	tokenString, err := utils.GenerateResetToken(email)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			"Internal Server Error",
			"/electivote/home-page/",
		)
		return
	}
	
	err = sendResetPasswordEmail(email, tokenString)
	fmt.Println(err)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			"Internal Server Error",
			"electivote/home-page/",
		)
		return
	}

	c.Redirect(
		http.StatusFound,
		"/electivote/login-page/",
	)
}

func ViewResetPasswordPage(c *gin.Context) {
	if middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}
	token := c.Param("token")
	email, err := verifyResetToken(token)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusBadRequest,
			"Invalid Token",
			"/electivote/home-page/",
		)
		return
	}

	context := gin.H {
		"title": "Reset Password",
		"email": email,
		"token": token,
	}
	c.HTML(
		http.StatusOK,
		"resetPassword.html",
		context,
	)
}

func ResetPasswordPage(c *gin.Context) {
	if middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}
	token := c.Param("token")
	email, err := verifyResetToken(token)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusBadRequest,
			"Invalid Token",
			"/electivote/home-page/",
		)
		return
	}

	password := c.PostForm("password")
	passwordErr := ""

	if len(password) == 0 {
		passwordErr = "Password must be filled"
	}

	if len(password) < 5 || len(password) > 255 {
		passwordErr = "Password must be between 5 and 255 characters"
	}

	if passwordErr != "" {
		context := gin.H {
			"title": "Reset Password",
			"passwordErr": passwordErr,
			"email": email,
			"token": token,
		}
		c.HTML(
			http.StatusOK,
			"resetPassword.html",
			context,
		)
		return
	}
	passwordHashed, err := utils.HashPassword(password)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			"Internal Server Error",
			"/electivote/home-page/",
		)
		return
	}
	_, err = repositories.UpdatePasswordByEmail(email, passwordHashed)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			"Internal Server Error",
			"/electivote/home-page/",
		)
		return
	}
	c.Redirect(
		http.StatusFound,
		"/electivote/login-page/",
	)
}