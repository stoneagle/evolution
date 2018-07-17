package middles

import (
	"errors"
	"evolution/backend/common/resp"
	systemModel "evolution/backend/system/models"
	"evolution/backend/system/services"
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     error         = errors.New("Token is expired")
	TokenNotValidYet error         = errors.New("Token not active yet")
	TokenMalformed   error         = errors.New("That's not even a token")
	TokenInvalid     error         = errors.New("Couldn't handle this token:")
	SignKey          string        = "system"
	JwtTokenHeader   string        = "JwtToken"
	ExpireHours      time.Duration = 24
)

type UserInfo struct {
	Id    int    `json:"Id"`
	Name  string `json:"Name"`
	Email string `json:"Email"`
	jwt.StandardClaims
}

func JWTAuthLogin(userSvc *services.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.MustGet(gin.AuthUserKey).(string)
		condition := systemModel.User{
			Name: username,
		}
		user, err := userSvc.OneByCondition(&condition)
		if err != nil {
			resp.ErrorBusiness(c, resp.ErrorSign, fmt.Sprintf("get user %s error", username), err)
			return
		}
		j := NewJWT()
		userInfo := UserInfo{}
		userInfo.Id = user.Id
		userInfo.Name = user.Name
		userInfo.Email = user.Email
		token, err := j.CreateToken(userInfo)
		if err != nil {
			resp.ErrorBusiness(c, resp.ErrorSign, fmt.Sprintf("create token %s error", username), err)
			return
		}
		c.Header(JwtTokenHeader, "Bear "+token)
		resp.Success(c, userInfo)
	}
}

func JWTAuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.DefaultQuery("token", "")
		if token == "" {
			token = c.Request.Header.Get("Authorization")
			if s := strings.Split(token, " "); len(s) == 2 {
				token = s[1]
			}
		}
		j := NewJWT()
		userInfo, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				if token, err = j.RefreshToken(token); err == nil {
					c.Header(JwtTokenHeader, "Bear "+token)
					resp.Success(c, userInfo)
					return
				}
			}
			resp.ErrorBusiness(c, resp.ErrorSign, "token check error", err)
			return
		}
		c.Set(UserKey, *userInfo)
	}
}

func JWTAuthCurrent() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.DefaultQuery("token", "")
		if token == "" {
			token = c.Request.Header.Get("Authorization")
			if s := strings.Split(token, " "); len(s) == 2 {
				token = s[1]
			}
		}
		j := NewJWT()
		userInfo, err := j.ParseToken(token)
		if err != nil {
			resp.ErrorBusiness(c, resp.ErrorSign, "token check error", err)
			return
		}
		resp.Success(c, userInfo)
	}
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

func GetSignKey() string {
	return SignKey
}

func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

func (j *JWT) CreateToken(userInfo UserInfo) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userInfo)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenString string) (*UserInfo, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserInfo{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if userInfo, ok := token.Claims.(*UserInfo); ok && token.Valid {
		return userInfo, nil
	}
	return nil, TokenInvalid
}

func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &UserInfo{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if userInfo, ok := token.Claims.(*UserInfo); ok && token.Valid {
		jwt.TimeFunc = time.Now
		userInfo.StandardClaims.ExpiresAt = time.Now().Add(ExpireHours * time.Hour).Unix()
		return j.CreateToken(*userInfo)
	}
	return "", TokenInvalid
}
