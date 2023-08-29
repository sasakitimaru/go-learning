package usecase_test

import (
	"errors"
	"fmt"
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

func TestSignUp(t *testing.T) {
	mockUserRepo := new(mockUserRepository)
	uc := usecase.NewUserUseCase(mockUserRepo)
	user := model.User{
		Email:    "test@example.com",
		Password: "password",
	}

	mockUserRepo.On("CreateUser", mock.AnythingOfType("*model.User")).Return(nil).Once()
	_, err := uc.SignUp(user)
	assert.NoError(t, err)

	mockError := errors.New("something went wrong")
	mockUserRepo.On("CreateUser", mock.AnythingOfType("*model.User")).Return(mockError)
	_, err = uc.SignUp(user)
	assert.Error(t, err)
	assert.Equal(t, mockError, err)
}

func TestLogIn(t *testing.T) {
	mockUserRepo := new(mockUserRepository)
	uc := usecase.NewUserUseCase(mockUserRepo)
	user := model.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "password",
	}
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), 10)
	storedUser := model.User{
		ID:       1,
		Email:    "test@example.com",
		Password: string(hash),
	}
	fmt.Println("storedUser:", storedUser)
	mockUserRepo.On("GetUserByEmail", mock.AnythingOfType("*model.User"), user.Email).Return(storedUser, nil).Once()
	tokenString, err := uc.LogIn(user)
	assert.NoError(t, err)
	parsedToken, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	claims := parsedToken.Claims.(jwt.MapClaims)
	fmt.Println("claims:", claims)
	assert.Equal(t, float64(storedUser.ID), claims["user_id"])

}
