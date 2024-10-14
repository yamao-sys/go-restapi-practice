package services

import (
	"app/dto"
	"app/models"
	"app/repositories"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignUp(ctx *gin.Context) (models.User, error)
	SignIn(ctx *gin.Context) (string, error)
	Getuser(id int) models.User
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

func (as *authService) SignIn(ctx *gin.Context) (string, error) {
	requestParams := dto.SignInRequest{}
	if err := ctx.ShouldBind(&requestParams); err != nil {
		log.Fatalln(err)
	}

	// NOTE: emailからユーザの取得
	user := models.User{}
	if err := as.userRepository.FindUserByEmail(&user, requestParams.Email); err != nil {
		return "", err
	}

	// NOTE: パスワードの照合
	if err := as.compareHashPassword(user.Password, requestParams.Password); err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	// TODO: JWT_SECRETを環境変数に切り出す
	tokenString, err := token.SignedString([]byte("abcdefghijklmn"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (as *authService) Getuser(id int) models.User {
	return as.userRepository.FindUserById(id)
}

// NOTE: パスワードの文字列をハッシュ化する
func (as *authService) encryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// NOTE: パスワードの照合
func (as *authService) compareHashPassword(hashedPassword, requestPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(requestPassword)); err != nil {
		return err
	}
	return nil
}
