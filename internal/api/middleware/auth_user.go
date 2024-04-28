package middleware

import "github.com/gin-gonic/gin"

func AuthUser(c *gin.Context) {
	// check user token is valid
	// if not valid, return 401
	// if valid, set user id to context
	// and call next

	c.Next()
}
