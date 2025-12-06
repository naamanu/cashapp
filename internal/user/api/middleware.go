package api

import (
	"cashapp/internal/user/models"
	"cashapp/internal/user/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireKYC ensures the user has a sufficient KYC level
func RequireKYC(minLevel int, s *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// In a real system, the user ID would come from the JWT auth middleware
		// For this example, we'll look for a "X-User-Tag" header for simulation
		userTag := c.GetHeader("X-User-Tag")
		if userTag == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: missing X-User-Tag header"})
			return
		}

		response := s.GetUser(userTag)
		if response.Error || response.Meta.Data == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid user"})
			return
		}

		// Parse user (unsafe map access for brevity, in real app use struct)
		data := response.Meta.Data.(*map[string]interface{})
		userObj := (*data)["user"].(*models.User)

		if userObj.KYCLevel < minLevel {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":      "insufficient KYC level",
				"required":   minLevel,
				"current":    userObj.KYCLevel,
				"kyc_status": userObj.KYCStatus,
			})
			return
		}

		c.Next()
	}
}
