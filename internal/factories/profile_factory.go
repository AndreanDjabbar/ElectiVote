package factories

import (
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
)

func CreateFirstProfile(userId int) (models.Profile) {
	newProfile := models.Profile{
		UserID: uint(userId),
		Picture: "default.png",
	}
	return newProfile
}

func CreateProfile(firstName, lastName, phone, picture string, age uint, dob models.NullTime) (models.Profile) {
	if picture == "" {
		picture = "default.png"
	}
	newProfile := models.Profile{
		FirstName: firstName,
		LastName: lastName,
		PhoneNumber: phone,
		Age: age,
		Birthday: dob,
		Picture: picture,
	}
	return newProfile
}