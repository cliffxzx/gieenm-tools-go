package validator

import (
	"context"

	"github.com/cliffxzx/gieenm-tools/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidatorKey ...
const ValidatorKey utils.ContextKey = utils.ContextKey("ValidatorKey")

// Validator ...
func Validator(c *gin.Context) {
	validtor := validator.New()
	ctx := context.WithValue(c.Request.Context(), ValidatorKey, validtor)
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

// Get ...
func Get(ctx context.Context) *validator.Validate {
	return ctx.Value(ValidatorKey).(*validator.Validate)
}
