package usecase_test

import (
	"errors"
	"go-rest-api/model"
	"os"

	// "go-rest-api/repository"
	"go-rest-api/usecase"
	"testing"

	// "github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	mock.Mock
}

type mockUserValidator struct {
	mock.Mock
}

type mockPasswordHasher struct {
	mock.Mock
}

func (m *mockPasswordHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	args := m.Called(password, cost)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockUserRepository) CreateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockUserRepository) GetUserByEmail(user *model.User, email string) error {
	args := m.Called(user, email)
	if args.Get(0) != nil {
		*user = args.Get(0).(model.User)
	}
	return args.Error(1)
}

func (m *mockUserValidator) UserValidate(user model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestSignUpSuccess(t *testing.T) {
	mockUserRepository := new(mockUserRepository)
	mockUserValidator := new(mockUserValidator)
	mockPasswordHasher := new(mockPasswordHasher)
	uc := usecase.NewUserUseCase(mockUserRepository, mockUserValidator, mockPasswordHasher)
	user := model.User{
		Email:    "test@example.com",
		Password: "password",
	}

	mockPasswordHasher.On("GenerateFromPassword", []byte(user.Password), 10).Return([]byte("HashedPasswordShouldBeHere"), nil)
	mockUserValidator.On("UserValidate", user).Return(nil)
	mockUserRepository.On("CreateUser", mock.AnythingOfType("*model.User")).Return(nil)
	_, err := uc.SignUp(user)
	assert.NoError(t, err)
}

func TestSignUpPasswordHasherFaild(t *testing.T) {
	mockUserRepository := new(mockUserRepository)
	mockUserValidator := new(mockUserValidator)
	mockPasswordHasher := new(mockPasswordHasher)
	uc := usecase.NewUserUseCase(mockUserRepository, mockUserValidator, mockPasswordHasher)
	user := model.User{
		Email:    "test@example.com",
		Password: "password",
	}

	mockError := errors.New("PasswordHasher failed")
	mockPasswordHasher.On("GenerateFromPassword", []byte(user.Password), 10).Return([]byte(""), mockError)
	mockUserValidator.On("UserValidate", user).Return(nil)
	mockUserRepository.On("CreateUser", mock.AnythingOfType("*model.User")).Return(nil)
	_, err := uc.SignUp(user)
	assert.Error(t, err)
	assert.Equal(t, mockError, err)
}

func TestSignUpCreateUserFaild(t *testing.T) {
	mockUserRepository := new(mockUserRepository)
	mockUserValidator := new(mockUserValidator)
	mockPasswordHasher := new(mockPasswordHasher)
	uc := usecase.NewUserUseCase(mockUserRepository, mockUserValidator, mockPasswordHasher)
	user := model.User{
		Email:    "test@example.com",
		Password: "password",
	}

	mockError := errors.New("CreateUser failed")
	mockPasswordHasher.On("GenerateFromPassword", []byte(user.Password), 10).Return([]byte("HashedPasswordShouldBeHere"), nil)
	mockUserValidator.On("UserValidate", user).Return(nil)
	mockUserRepository.On("CreateUser", mock.AnythingOfType("*model.User")).Return(mockError)
	_, err := uc.SignUp(user)
	assert.Error(t, err)
	assert.Equal(t, mockError, err)
}

func TestSignUpValidationFaild(t *testing.T) {
	mockUserRepository := new(mockUserRepository)
	mockUserValidator := new(mockUserValidator)
	mockPasswordHasher := new(mockPasswordHasher)
	uc := usecase.NewUserUseCase(mockUserRepository, mockUserValidator, mockPasswordHasher)
	user := model.User{
		Email:    "test@example.com",
		Password: "password",
	}

	mockError := errors.New("Validation failed")
	mockPasswordHasher.On("GenerateFromPassword", []byte(user.Password), 10).Return([]byte("HashedPasswordShouldBeHere"), nil)
	mockUserValidator.On("UserValidate", user).Return(mockError)
	mockUserRepository.On("CreateUser", mock.AnythingOfType("*model.User")).Return(nil)
	_, err := uc.SignUp(user)
	assert.Error(t, err)
	assert.Equal(t, mockError, err)
}

func TestLogInSuccess(t *testing.T) {
	mockPasswordHasher := new(mockPasswordHasher)
	mockUserRepo := new(mockUserRepository)
	mockUserValid := new(mockUserValidator)
	uc := usecase.NewUserUseCase(mockUserRepo, mockUserValid, mockPasswordHasher)
	user := model.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "password",
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	storedUser := model.User{
		ID:       1,
		Email:    "test@example.com",
		Password: string(hash),
	}

	mockUserValid.On("UserValidate", user).Return(nil)
	mockUserRepo.On("GetUserByEmail", mock.AnythingOfType("*model.User"), user.Email).Return(storedUser, nil)
	tokenString, err := uc.LogIn(user)
	assert.NoError(t, err)
	parsedToken, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	claims := parsedToken.Claims.(jwt.MapClaims)
	assert.Equal(t, float64(storedUser.ID), claims["user_id"])

}

func TestLogInGetUserByEmailFaild(t *testing.T) {
	mockPasswordHasher := new(mockPasswordHasher)
	mockUserRepo := new(mockUserRepository)
	mockUserValid := new(mockUserValidator)
	uc := usecase.NewUserUseCase(mockUserRepo, mockUserValid, mockPasswordHasher)
	user := model.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "password",
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	storedUser := model.User{
		ID:       1,
		Email:    "test@example.com",
		Password: string(hash),
	}

	mockError := errors.New("GetUserByEmail failed")
	mockUserValid.On("UserValidate", user).Return(nil)
	mockUserRepo.On("GetUserByEmail", mock.AnythingOfType("*model.User"), user.Email).Return(storedUser, mockError)
	_, err := uc.LogIn(user)
	assert.Error(t, err)
	assert.Equal(t, mockError, err)
}
