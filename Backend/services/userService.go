package services

import "github.com/arttkachev/X-Airlines/Backend/api/models"

// temp storage
var users []models.User = make([]models.User, 0)

func GetRepository() *[]models.User {
	return &users
}
