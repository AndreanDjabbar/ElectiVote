package handlers

import (
	"fmt"
	"log"
	"net/http"
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
		c.Redirect(http.StatusFound, "/electivote/home-page/")
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
			log.Printf("Error sending email: %v", err)
		}
	}()

	c.Redirect(http.StatusFound, "/electivote/login-page/")
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