package middleware

import "github.com/gin-gonic/gin"

const userIDKey = "userID"

// SetUserID stores the authenticated user's ID in the Gin context.
func SetUserID(c *gin.Context, id uint) {
	c.Set(userIDKey, id)
}

// GetUserID retrieves the authenticated user's ID from the Gin context.
// Returns (0, false) if not set.
func GetUserID(c *gin.Context) (uint, bool) {
	val, exists := c.Get(userIDKey)
	if !exists {
		return 0, false
	}
	id, ok := val.(uint)
	return id, ok
}
