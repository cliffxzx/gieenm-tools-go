package utils

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

// ContextKey ...
type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

// GinContextSelfKey ...
const GinContextSelfKey ContextKey = ContextKey("GinContextSelf")

// StoreGinContextToSelf ...
func StoreGinContextToSelf(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), GinContextSelfKey, c)
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

// GetGinContext ...
func GetGinContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(GinContextSelfKey)
	if ginContext == nil {
		err := errors.New("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := errors.New("gin.Context has wrong type")
		return nil, err
	}

	return gc, nil
}
