package factories

import "github.com/AndreanDjabbar/ElectiVote/internal/models"

func SupportFactory(donatorName string, donatorEmail string, transactionID string, amount float64, message string, createdAt models.CustomTime) models.Support {
	return models.Support{
		DonatorEmail: donatorEmail,
		DonatorName:  donatorName,
		SupportedTime: createdAt,
		TransactionID: transactionID,
		Amount:        amount,
		Message:       message,
	}
}