package handlers

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/AndreanDjabbar/CaysAPIHub/internal/factories"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/middlewares"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/models"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/repositories"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/utils"
	"github.com/gin-gonic/gin"
)

func ViewProfilePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}

	username := middlewares.GetUserData(c)
	var userProfile models.Profile
	var userEmail string
	var err, profileErr, emailErr error
	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(2)

	go func() {
		defer wg.Done()
		if userProfile, profileErr = repositories.GetProfilesByUsername(username); profileErr != nil {
			mu.Lock()
			profileErr = err
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		if userEmail, emailErr = repositories.GetUserEmailByUsername(username); emailErr != nil {
			mu.Lock()
			emailErr = err
			mu.Unlock()
		}
	}()

	wg.Wait()

	if profileErr != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			profileErr.Error(),
			"/electivote/profile-page/",
		)
		return
	}

	if emailErr != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			emailErr.Error(),
			"/electivote/profile-page/",
		)
		return
	}

	formattedDob := utils.FormattedDob(userProfile.Birthday)
	context := gin.H{
		"title":       "Edit Profile",
		"username":    username,
		"userProfile": userProfile,
		"userEmail":   userEmail,
		"birthday":    formattedDob,
	}

	c.HTML(
		http.StatusOK,
		"profile.html",
		context,
	)
}

func ViewEditProfilePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}

	username := middlewares.GetUserData(c)
	var userProfile models.Profile
	var userEmail string
	var err, profileErr, emailErr error
	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(2)

	go func() {
		defer wg.Done()
		if userProfile, err = repositories.GetProfilesByUsername(username); profileErr != nil {
			mu.Lock()
			profileErr = err
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		if userEmail, err = repositories.GetUserEmailByUsername(username); emailErr != nil {
			mu.Lock()
			emailErr = err
			mu.Unlock()
		}
	}()

	wg.Wait()

	if profileErr != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			profileErr.Error(),
			"/electivote/edit-profile-page/",
		)
		return
	}

	if emailErr != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			emailErr.Error(),
			"/electivote/edit-profile-page/",
		)
		return
	}

	formattedDob := utils.FormattedDob(userProfile.Birthday)
	context := gin.H{
		"title":       "Edit Profile",
		"username":    username,
		"userProfile": userProfile,
		"userEmail":   userEmail,
		"birthday":    formattedDob,
	}

	c.HTML(
		http.StatusOK,
		"editProfile.html",
		context,
	)
}

func EditProfilePage(c *gin.Context) {
    if !middlewares.IsLogged(c) {
        c.Redirect(
            http.StatusFound,
            "/electivote/login-page/",
        )
        return
    }
    username := middlewares.GetUserData(c)
    var firstNameErr, lastNameErr, phoneErr, ageErr, dobErr string
    var finalDOB models.NullTime
	var userProfile models.Profile
	var convertedAge uint

	userProfile, err := repositories.GetProfilesByUsername(username)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/profile-page/",
		)
	}

    firstName := c.PostForm("firstName")
    lastName := c.PostForm("lastName")
    phoneNumber := c.PostForm("phone")
    dob := c.PostForm("dob")
    age := c.PostForm("age")
	file, fileErr := c.FormFile("picture")
    
	
    if age != "" {
        convAge, err := strconv.Atoi(age)
        if err != nil {
            ageErr = "Age must be a number"
        } else {
            convertedAge = uint(convAge)
        }
    }

	userEmail, err := repositories.GetUserEmailByUsername(username)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/profile-page/",
		)
	
	}

    if dob != "" {
        parsedDob, err := time.Parse("2006-01-02", dob)
        if err != nil {
            dobErr = "Date of Birth must be in format YYYY-MM-DD"
            finalDOB = models.NullTime{Valid: false}
        } else {
            finalDOB = models.NullTime{Time: parsedDob, Valid: true}
        }
    } else {
        finalDOB = models.NullTime{Valid: false}
    }

    firstNameErr, lastNameErr, phoneErr, ageErr = utils.ValidateProfileInput(firstName, lastName, phoneNumber, convertedAge)
	formattedDob := utils.FormattedDob(userProfile.Birthday)

    
    if firstNameErr == "" && lastNameErr == "" && phoneErr == "" && ageErr == "" && dobErr == "" {
        newProfile := factories.CreateProfile(firstName, lastName, phoneNumber, userProfile.Picture, convertedAge, finalDOB)

        if file != nil {
			if fileErr != nil {
				utils.RenderError(
					c,
					http.StatusInternalServerError,
					fileErr.Error(),
					"/electivote/edit-profile-page/",
				)
			}
			newProfile.Picture = file.Filename
			err = c.SaveUploadedFile(file, "../../internal/assets/images/"+file.Filename)
			if err != nil {
				utils.RenderError(
					c,
					http.StatusInternalServerError,
					err.Error(),
					"/electivote/edit-profile-page/",
				)
				return
			}
			_, err = repositories.UpdateProfileByUsername(username, newProfile)
			if err != nil {
				utils.RenderError(
					c,
					http.StatusInternalServerError,
					err.Error(),
					"/electivote/edit-profile-page/",
				)
				return
			}
			c.Redirect(
				http.StatusFound,
				"/electivote/profile-page/",
			)
			return
        }

        _, err := repositories.UpdateProfileByUsername(username, newProfile)
        if err != nil {
            utils.RenderError(
                c,
                http.StatusInternalServerError,
                err.Error(),
                "/electivote/edit-profile-page/",
            )
            return
        }
        c.Redirect(
            http.StatusFound,
            "/electivote/profile-page/",
        )
        return
    }

    context := gin.H{
        "firstNameErr": firstNameErr,
		"firstName":    firstName,
        "lastNameErr":  lastNameErr,
        "phoneErr":     phoneErr,
        "ageErr":       ageErr,
		"userEmail":	userEmail,
		"birthday":		formattedDob,
		"userProfile":	userProfile,
        "dobErr":       dobErr,
    }
    c.HTML(
        http.StatusOK,
        "editProfile.html",
        context,
    )
}