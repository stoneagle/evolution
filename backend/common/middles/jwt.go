package middles

import (
	"errors"
	"evolution/backend/common/config"
	"evolution/backend/common/database"
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

type CustomClaims struct {
	Id    int    `json:"Id"`
	Name  string `json:"Name"`
	Email string `json:"Email"`
	jwt.StandardClaims
}

func JWTAuthLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.MustGet(gin.AuthUserKey).(string)
		condition := systemModel.User{
			Name: username,
		}
		user, err := services.NewUser(database.GetXorm(config.Get().System.System.Name), nil).OneByCondition(&condition)
		if err != nil {
			resp.ErrorBusiness(c, resp.ErrorLogin, fmt.Sprintf("get user %s error", username), err)
			return
		}
		j := NewJWT()
		claim := CustomClaims{}
		claim.Id = user.Id
		claim.Name = user.Name
		claim.Email = user.Email
		token, err := j.CreateToken(claim)
		if err != nil {
			resp.ErrorBusiness(c, resp.ErrorLogin, fmt.Sprintf("create token %s error", username), err)
			return
		}
		c.Header(JwtTokenHeader, "Bear "+token)
		resp.Success(c, claim)
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
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				if token, err = j.RefreshToken(token); err == nil {
					c.Header(JwtTokenHeader, "Bear "+token)
					resp.Success(c, claims)
					return
				}
			}
			resp.ErrorBusiness(c, resp.ErrorLogin, "token check error", err)
			return
		}
		c.Set(UserKey, claims)
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
		claims, err := j.ParseToken(token)
		if err != nil {
			resp.ErrorBusiness(c, resp.ErrorLogin, "token check error", err)
			return
		}
		resp.Success(c, claims)
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

func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
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
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(ExpireHours * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
