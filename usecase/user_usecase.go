package usecase

import (
	"go-rest-api/model"
	"go-rest-api/repository"
	"go-rest-api/validator"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type IUserUseCase interface {
	SignUp(user model.User) (model.UserResponse, error)
	LogIn(user model.User) (string, error)
}

type userUseCase struct {
	ur repository.IUserRepository
	uv validator.IUserValidator
	ph PasswordHasher
}

type PasswordHasher interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
}

type BycryptPasswordHasher struct{}

func (b *BycryptPasswordHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func NewUserUseCase(ur repository.IUserRepository, uv validator.IUserValidator, ph PasswordHasher) IUserUseCase {
	if ph == nil {
		ph = &BycryptPasswordHasher{}
	}
	return &userUseCase{ur, uv, ph}
}

func (uu *userUseCase) SignUp(user model.User) (model.UserResponse, error) {
	if err := uu.uv.UserValidate(user); err != nil {
		return model.UserResponse{}, err
	}
	hash, err := uu.ph.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return model.UserResponse{}, err
	}
	newUser := model.User{Email: user.Email, Password: string(hash)}
	if err := uu.ur.CreateUser(&newUser); err != nil {
		return model.UserResponse{}, err
	}
	resUser := model.UserResponse{ID: newUser.ID, Email: newUser.Email}
	return resUser, nil
}

func (uu *userUseCase) LogIn(user model.User) (string, error) {
	if err := uu.uv.UserValidate(user); err != nil {
		return "", err
	}
	storedUser := model.User{}
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		return "", err
	}
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
