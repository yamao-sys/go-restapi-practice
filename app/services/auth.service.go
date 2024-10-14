package services

import (
	"app/dto"
	"app/models"
	"app/repositories"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignUp(ctx *gin.Context) (models.User, error)
}

type authService struct {
	userRepository repositories.UserRepository
}

func NewAuthService(userRepository repositories.UserRepository) AuthService {
	return &authService{userRepository}
}

func (as *authService) SignUp(ctx *gin.Context) (models.User, error) {
	// NOTE: リクエストデータを構造体に変換
	requestParams := dto.SignUpRequest{}
	if err := ctx.ShouldBind(&requestParams); err != nil {
		log.Fatalln(err)
	}

	// TODO: バリデーションチェック
	// NOTE: パスワードのハッシュ化
	hashedPassword, err := as.encryptPassword(requestParams.Password)
	if err != nil {
		log.Fatalln(err)
	}

	user := models.User{}
	user.Name = requestParams.Name
	user.Email = requestParams.Email
	user.Password = hashedPassword
	// NOTE: Create処理
	as.userRepository.CreateUser(&user)

	return user, nil
}

// NOTE: パスワードの文字列をハッシュ化する
func (as *authService) encryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
