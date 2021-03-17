package authentication

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/cliffxzx/gieenm-tools/pkg/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

func verifyToken(c *gin.Context) (*jwt.Token, error) {
	headerToken := c.GetHeader("Authorization")
	if isMatch, err := regexp.MatchString("^Bearer .*$", headerToken); err != nil || !isMatch {
		return nil, errors.New("Can't parsing authentication from header")
	}

	tokenStr := &strings.Split(headerToken, " ")[1]

	token, err := jwt.Parse(*tokenStr, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		accessSecret := utils.MustGetEnv("TOKEN_SECRET")

		return []byte(accessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func extractToken(token *jwt.Token) (*Token, error) {
	if token.Valid {
		jwtClaims := token.Claims.(jwt.MapClaims)
		return &Token{
			Claims: &Claims{
				UserID:  funk.PtrOf(int(jwtClaims["user_id"].(float64))).(*int),
				TokenID: funk.PtrOf(jwtClaims["token_id"]).(*string),
			},
		}, nil
	}

	return nil, errors.New("Invalid Token")
}

// VerifyToken ...
func VerifyToken(ctx *gin.Context) (*Token, error) {
	jwtToken, err := verifyToken(ctx)
	if err != nil {
		return nil, errors.New("Required authentication token")
	}

	token, err := extractToken(jwtToken)
	if err != nil {
		return nil, errors.New("Invalid authentication token")
	}

	err = GetAuth(token)
	if err != nil {
		//Token does not exists in Redis (User logged out or expired)
		return nil, errors.New("Undefined authentication token, maybe logged out or expired")
	}

	return token, nil
}
