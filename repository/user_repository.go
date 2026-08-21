package repository

import (
	"user-api/config"
	"user-api/models"
)

func CreateUser(user *models.User) error {
	result := config.DB.Create(user)
	return result.Error
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User

	result := config.DB.Find(&users)

	return users, result.Error
}

func GetUserByID(id string) (models.User, error) {
	var user models.User

	result := config.DB.First(&user, id)

	return user, result.Error
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User

	result := config.DB.Where("email = ?", email).First(&user)

	return user, result.Error
}

func UpdateUser(user *models.User) error {
	result := config.DB.Save(user)
	return result.Error
}

func DeleteUser(user *models.User) error {
	result := config.DB.Delete(user)
	return result.Error
}