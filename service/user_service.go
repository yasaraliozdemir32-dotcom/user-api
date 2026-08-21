package service

import (
	"os"
	"time"

	"user-api/models"
	"user-api/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return repository.CreateUser(user)
}

func GetAllUsers() ([]models.User, error) {
	return repository.GetAllUsers()
}

func GetUserByID(id string) (models.User, error) {
	return repository.GetUserByID(id)
}

func UpdateUser(user *models.User, passwordChanged bool) error {

	if passwordChanged {
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(user.Password),
			bcrypt.DefaultCost,
		)

		if err != nil {
			return err
		}

		user.Password = string(hashedPassword)
	}

	return repository.UpdateUser(user)
}

func DeleteUser(user *models.User) error {
	return repository.DeleteUser(user)
}

func Login(email string, password string) (string, error) {

	users, err := repository.GetAllUsers()

	if err != nil {
		return "", err
	}

	var user models.User

	for _, u := range users {
		if u.Email == email {
			user = u
			break
		}
	}

	if user.ID == 0 {
		return "", bcrypt.ErrMismatchedHashAndPassword
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if err != nil {
		return "", err
	}

	secret := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
