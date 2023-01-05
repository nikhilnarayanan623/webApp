package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func AdminAuth(ctx *gin.Context) {
	fmt.Println("Admin Auth")
}
