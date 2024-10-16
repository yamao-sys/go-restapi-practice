package services

import (
	"app/dto"
	"app/models"
	"app/repositories"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignUp(ctx *gin.Context) *dto.SignUpResponse
	SignIn(ctx *gin.Context) *dto.SignInResponse
	Getuser(id int) models.User
}

type authService struct {
	userRepository repositories.UserRepository
}

func NewAuthService(userRepository repositories.UserRepository) AuthService {
	return &authService{userRepository}
}

func (as *authService) SignUp(ctx *gin.Context) *dto.SignUpResponse {
	// NOTE: リクエストデータを構造体に変換
	requestParams := dto.SignUpRequest{}
	if err := ctx.ShouldBind(&requestParams); err != nil {
		return &dto.SignUpResponse{User: models.User{}, Error: err, ErrorType: "internalServerError"}
	}

	user := models.User{}
	user.Name = requestParams.Name
	user.Email = requestParams.Email
	user.Password = requestParams.Password
	// NOTE: バリデーションチェック
	validate := validator.New()
	validationErrors := validate.Struct(user)
	if validationErrors != nil {
		return &dto.SignUpResponse{User: user, Error: validationErrors, ErrorType: "validationError"}
	}

	// NOTE: パスワードをハッシュ化の上、Create処理
	hashedPassword, err := as.encryptPassword(requestParams.Password)
	if err != nil {
		return &dto.SignUpResponse{User: user, Error: err, ErrorType: "internalServerError"}
	}
	user.Password = hashedPassword
	as.userRepository.CreateUser(&user)

	return &dto.SignUpResponse{User: user, Error: nil, ErrorType: ""}
}

func (as *authService) SignIn(ctx *gin.Context) *dto.SignInResponse {
	requestParams := dto.SignInRequest{}
	if err := ctx.ShouldBind(&requestParams); err != nil {
		return &dto.SignInResponse{TokenString: "", NotFoundMessage: "", Error: err}
	}

	// NOTE: emailからユーザの取得
	user := models.User{}
	if err := as.userRepository.FindUserByEmail(&user, requestParams.Email); err != nil {
		return &dto.SignInResponse{TokenString: "", NotFoundMessage: "メールアドレスまたはパスワードに該当するユーザが存在しません。", Error: nil}
	}

	// NOTE: パスワードの照合
	if err := as.compareHashPassword(user.Password, requestParams.Password); err != nil {
		return &dto.SignInResponse{TokenString: "", NotFoundMessage: "メールアドレスまたはパスワードに該当するユーザが存在しません。", Error: nil}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	// TODO: JWT_SECRETを環境変数に切り出す
	tokenString, err := token.SignedString([]byte("abcdefghijklmn"))
	if err != nil {
		return &dto.SignInResponse{TokenString: "", NotFoundMessage: "", Error: err}
	}
	return &dto.SignInResponse{TokenString: tokenString, NotFoundMessage: "", Error: nil}
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
