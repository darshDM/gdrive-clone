package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/darshDM/gdrive-clone-api/internal/store"
	"github.com/darshDM/gdrive-clone-api/types"
	"github.com/darshDM/gdrive-clone-api/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey string

type UserService struct {
	store store.Store
}

func NewUserService(store store.Store) *UserService {
	secretKey = utils.GetStringEnv("JWT_SECRET_KEY", "secret")
	return &UserService{store: store}
}

var (
	ErrUserAlreadyExist = errors.New("user already exists, please use another username")
)

func GenerateHashPassword(password string) (string, error) {
	fmt.Println(password)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (userService *UserService) CreateNewUser(ctx context.Context, user *types.CreateUserRequest) error {
	_, err := userService.store.GetUser(ctx, user.Username)

	if err != nil {
		passwordHash, err := GenerateHashPassword(user.Password)
		fmt.Println(passwordHash)
		if err != nil {
			return err
		}
		newUser := &store.User{
			Username:     user.Username,
			Password:     passwordHash,
			TotalStorage: 10000000,
			UsedStorage:  0,
		}

		if err := userService.store.CreateNewUser(ctx, newUser); err != nil {
			return err
		}
		return nil
	}
	return ErrUserAlreadyExist
}

func (userService *UserService) LoginUser(ctx context.Context, userRequest *types.LoginUserRequest) (string, error) {
	user, err := userService.store.GetUser(ctx, userRequest.Username)
	if err != nil {
		log.Println("Error getting user: ", err)
		return "", err
	}
	if !CheckHashPassword(userRequest.Password, user.Password) {
		return "", errors.New("invalid password")
	}
	token, err := createToken(user.Username)
	if err != nil {
		log.Println("Error creating token: ", err)
		return "", err
	}
	return token, nil

}

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (userService *UserService) Authenticate(tokenString string) (*store.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	user, err := userService.store.GetUser(context.Background(), username)
	if err != nil {
		return nil, fmt.Errorf("user doesn't exist.")
	}
	return user, nil
}
