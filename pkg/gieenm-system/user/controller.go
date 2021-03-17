package user

import (
	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/authentication"
	"github.com/gin-gonic/gin"
)

func GetByHeaderController(ctx *gin.Context) (*User, error) {
	token, err := authentication.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	u, err := Get(*token.Claims.UserID)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
