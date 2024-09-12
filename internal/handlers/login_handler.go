package handlers

import (
	"fmt"
	"net/http"
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
		username := middlewares.GetUserData(c)
		logger.Warn(
			"ViewLoginPage - User already logged in",
			"username", username,
			"client IP", c.ClientIP(),
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	logger.Info(
		"ViewLoginPage - Login page acessed",
		"Client IP", c.ClientIP(),
	)

	siteKey := os.Getenv("RECAPTCHA_SITE_KEY")
	context := gin.H{
		"title":   "Login",
		"siteKey": siteKey,
	}
	c.HTML(
		http.StatusOK,
		"login.html",
		context,
	)
}

func LoginPage(c *gin.Context) {
	if middlewares.IsLogged(c) {
		logger.Warn(
			"LoginPage - User already logged in",
			"username", middlewares.GetUserData(c),
			"client IP", c.ClientIP(),
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")
	remember := c.PostForm("remember")
	siteKey := os.Getenv("RECAPTCHA_SITE_KEY")

	usernameErr, passwordErr := utils.ValidateLoginInput(username, password, c)
	var usernameCheckErr, passwordCheckErr error
	var wg sync.WaitGroup
	var mu sync.Mutex
	captchaErr := ""

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
		logger.Warn(
			"LoginPage - Username not found",
			"Username Inputted", username,
			"client IP", c.ClientIP(),
		)
		usernameErr = "Username not found"
	}

	if passwordCheckErr != nil {
		logger.Warn(
			"LoginPage - Password is incorrect",
			"client IP", c.ClientIP(),
		)
		passwordErr = "Password is incorrect"
	}

	if !utils.IsValidReCAPTCHA(c) {
		logger.Warn(
			"LoginPage - Invalid ReCAPTCHA",
			"client IP", c.ClientIP(),
		)
		captchaErr = "Invalid ReCAPTCHA"
	}

	if usernameErr == "" && passwordErr == "" && captchaErr == "" {
		if remember == "on" {
			middlewares.SetCookies(c, username)
		} else {
			middlewares.SetSession(c, username)
		}
		logger.Info(
			"LoginPage - User logged in",
			"username", username,
			"client IP", c.ClientIP(),
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}
	context := gin.H{
		"title":       "Login",
		"usernameErr": usernameErr,
		"passwordErr": passwordErr,
		"captchaErr":  captchaErr,
		"username":    username,
		"password":    password,
		"siteKey":     siteKey,
	}
	c.HTML(
		http.StatusOK,
		"login.html",
		context,
	)
}

func LogoutPage(c *gin.Context) {
	username := middlewares.GetUserData(c)
	logger.Info(
		"LogoutPage - User logged out",
		"username", username,
		"client IP", c.ClientIP(),
		"action", "redirecting to login page",
	)
	middlewares.DeleteSession(c)
	middlewares.DeleteCookie(c)
	c.Redirect(
		http.StatusFound,
		"/electivote/login-page/",
	)
}

func ViewForgotPasswordPage(c *gin.Context) {
	if middlewares.IsLogged(c) {
		username := middlewares.GetUserData(c)
		logger.Warn(
			"ViewForgotPasswordPage - User already logged in",
			"username", username,
			"client IP", c.ClientIP(),
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}
	logger.Info(
		"ViewForgotPasswordPage - Rendering Forgot Password Page",
		"Client IP", c.ClientIP(),
	)
	context := gin.H{
		"title": "Forgot Password",
	}
	c.HTML(
		http.StatusOK,
		"forgotPassword.html",
		context,
	)
}

func verifyResetToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return utils.SecretKey, nil
	})

	if err != nil || !token.Valid {
		logger.Warn("verifyResetToken - Invalid or expired token")
		return "", fmt.Errorf("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		logger.Warn("verifyResetToken - Invalid token claims")
		return "", fmt.Errorf("invalid token claims")
	}
	logger.Info("verifyResetToken - Token verified")
	userEmail := claims["email"].(string)
	return userEmail, nil
}

func ForgotPasswordPage(c *gin.Context) {
	if middlewares.IsLogged(c) {
		username := middlewares.GetUserData(c)
		logger.Warn(
			"ForgotPasswordPage - User already logged in",
			"username", username,
			"client IP", c.ClientIP(),
			"action", "redirecting to home page",
		)
		c.Redirect(http.StatusFound, "/electivote/home-page/")
		return
	}

	logger.Info(
		"ForgotPasswordPage - Forgot Password page accessed",
		"Client IP", c.ClientIP(),
	)

	email := c.PostForm("email")
	emailErr := ""
	if len(email) == 0 {
		logger.Warn(
			"ForgotPasswordPage - Email must be filled",
			"Client IP", c.ClientIP(),
		)
		emailErr = "Email must be filled"
	}

	if !utils.IsValidEmail(email) {
		logger.Warn(
			"ForgotPasswordPage - Email is not valid",
			"Email Inputted", email,
			"Client IP", c.ClientIP(),
		)
		emailErr = "Email is not valid"
	}

	_, err := repositories.GetUserByEmail(email)
	if err != nil {
		logger.Warn(
			"ForgotPasswordPage - Email not found",
			"Email Inputted", email,
			"Client IP", c.ClientIP(),
		)
		emailErr = "Email not found"
	}

	if emailErr != "" {
		logger.Info(
			"ForgotPasswordPage - Rendering Forgot Password Page",
			"Client IP", c.ClientIP(),
		)
		context := gin.H{
			"title":    "Forgot Password",
			"emailErr": emailErr,
			"email":    email,
		}
		c.HTML(http.StatusOK, "forgotPassword.html", context)
		return
	}

	tokenString, err := utils.GenerateResetToken(email)
	if err != nil {
		logger.Error(
			"ForgotPasswordPage - Internal Server Error",
			"error", err,
		)
		utils.RenderError(c, http.StatusInternalServerError, "Internal Server Error", "/electivote/home-page/")
		return
	}
	emailDomain := utils.GetEmailDomain(email)
	emailProvider := utils.GetEmailProvider(emailDomain)
	resetURL := fmt.Sprintf("http://localhost:8080/electivote/reset-password-page/%s", tokenString)
	body := fmt.Sprintf(`
    <html>
    <head>
        <style>
            .email-container {
                font-family: Arial, sans-serif;
                line-height: 1.6;
                color: #333;
            }
            .email-header {
                font-size: 20px;
                font-weight: bold;
                margin-bottom: 20px;
            }
            .email-content {
                font-size: 16px;
                margin-bottom: 30px;
            }
            .email-link {
                display: inline-block;
                padding: 10px 15px;
                background-color: #4CAF50;
                color: white;
                text-decoration: none;
                border-radius: 5px;
            }
        </style>
    </head>
    <body>
        <div class="email-container">
            <div class="email-header">Reset Password</div>
            <div class="email-content">
                <p>To reset your password, please click the following link:</p>
                <p>
                    <a href="%s" class="email-link">Reset Your Password</a>
                </p>
                <p>If you did not request a password reset, please ignore this email.</p>
            </div>
        </div>
    </body>
    </html>
`, resetURL)

	subject := "Reset your password"
	go func() {
		err = utils.SendEmail(email, emailProvider, body, subject)
		if err != nil {
			logger.Error(
				"ForgotPasswordPage - Failed to send email",
				"error", err,
			)
		}
	}()
	logger.Info(
		"ForgotPasswordPage - Email sent",
		"client IP", c.ClientIP(),
		"action", "redirecting to login page",
	)
	c.Redirect(http.StatusFound, "/electivote/login-page/")
}

func ViewResetPasswordPage(c *gin.Context) {
	if middlewares.IsLogged(c) {
		username := middlewares.GetUserData(c)
		logger.Warn(
			"ViewResetPasswordPage - User already logged in",
			"username", username,
			"client IP", c.ClientIP(),
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	logger.Info(
		"ViewResetPasswordPage - Rendering Reset Password Page",
		"Client IP", c.ClientIP(),
	)

	token := c.Param("token")
	email, err := verifyResetToken(token)
	if err != nil {
		logger.Error(
			"ViewResetPasswordPage - Invalid Token",
			"error", err,
			"client IP", c.ClientIP(),
		)
		utils.RenderError(
			c,
			http.StatusBadRequest,
			"Invalid Token",
			"/electivote/home-page/",
		)
		return
	}
	context := gin.H{
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
		username := middlewares.GetUserData(c)
		logger.Warn(
			"ResetPasswordPage - User already logged in",
			"username", username,
			"client IP", c.ClientIP(),
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}
	token := c.Param("token")
	email, err := verifyResetToken(token)
	if err != nil {
		logger.Error(
			"ResetPasswordPage - Invalid Token",
			"client IP", c.ClientIP(),
			"error", err,
		)
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
		logger.Warn(
			"ResetPasswordPage - Password must be filled",
			"client IP", c.ClientIP(),
		)
		passwordErr = "Password must be filled"
	}

	if len(password) < 5 || len(password) > 255 {
		logger.Warn(
			"ResetPasswordPage - Password must be between 5 and 255 characters",
			"client IP", c.ClientIP(),
		)
		passwordErr = "Password must be between 5 and 255 characters"
	}

	if passwordErr != "" {
		context := gin.H{
			"title":       "Reset Password",
			"passwordErr": passwordErr,
			"email":       email,
			"token":       token,
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
		logger.Error(
			"ResetPasswordPage - Internal Server Error",
			"client IP", c.ClientIP(),
			"error", err,
		)
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
		logger.Error(
			"ResetPasswordPage - Internal Server Error",
			"client IP", c.ClientIP(),
			"error", err,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			"Internal Server Error",
			"/electivote/home-page/",
		)
		return
	}
	logger.Info(
		"ResetPasswordPage - Password has been reset",
		"client IP", c.ClientIP(),
		"action", "redirecting to login page",
	)
	c.Redirect(
		http.StatusFound,
		"/electivote/login-page/",
	)
}
