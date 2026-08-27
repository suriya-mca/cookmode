package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSON(c *gin.Context, code int, data any) {
	c.JSON(code, data)
}

func Err(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
}

func ErrBadRequest(c *gin.Context, msg string) {
	Err(c, http.StatusBadRequest, msg)
}

func ErrUnauthorized(c *gin.Context, msg string) {
	Err(c, http.StatusUnauthorized, msg)
}

func ErrForbidden(c *gin.Context, msg string) {
	Err(c, http.StatusForbidden, msg)
}

func ErrNotFound(c *gin.Context, msg string) {
	Err(c, http.StatusNotFound, msg)
}

func ErrInternal(c *gin.Context, msg string) {
	Err(c, http.StatusInternalServerError, msg)
}
