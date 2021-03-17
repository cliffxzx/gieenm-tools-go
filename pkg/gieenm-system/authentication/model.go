package authentication

import (
	"context"
	"strings"
	"time"

	"github.com/cliffxzx/gieenm-tools/pkg/gieenm-system/database"
	"github.com/cliffxzx/gieenm-tools/pkg/utils"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/twinj/uuid"
)

// Claims ...
type Claims struct {
	jwt.Claims
	UserID  *int    `json:"user_id,omitempty"`
	TokenID *string `json:"token_id,omitempty"`
}

// Token ...
type Token struct {
	Claims  *Claims
	Expires *int64
}

// Value ...
func (t *Token) Value() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t.Claims)

	accessSecret := utils.MustGetEnv("TOKEN_SECRET")

	return token.SignedString([]byte(accessSecret))
}

// CreateToken ...
func CreateToken(uid int) (*Token, error) {
	tokenID := strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	expires := time.Now().Add(time.Hour * 24 * 365).Unix()

	token := &Token{
		Claims: &Claims{
			UserID:  &uid,
			TokenID: &tokenID,
		},
		Expires: &expires,
	}

	return token, nil
}

var ctx context.Context = context.Background()

// CreateAuth ...
func CreateAuth(token *Token) error {
	expires := time.Unix(*token.Expires, 0)
	now := time.Now()

	err := database.GetRedis().Set(ctx, *token.Claims.TokenID, *token.Claims.TokenID, expires.Sub(now)).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetAuth ...
func GetAuth(token *Token) error {
	return database.GetRedis().Get(ctx, *token.Claims.TokenID).Err()
}

// DelAuth ...
func DelAuth(token *Token) error {
	return database.GetRedis().Del(ctx, *token.Claims.TokenID).Err()
}
