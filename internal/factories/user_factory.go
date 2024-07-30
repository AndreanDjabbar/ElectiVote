package factories

import "github.com/AndreanDjabbar/ElectiVote/internal/models"

func CreateUser(username, password, email, role string) (models.User) {
	newUser := models.User{
		Username: username,
		Password: password,
		Email:    email,
		Role:     role,
	}
	return newUser
}