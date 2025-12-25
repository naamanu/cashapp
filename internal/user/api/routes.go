package api

import (
	"cashapp/core"
	"cashapp/internal/user/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(e *gin.Engine, s *service.UserService) {
	// CreateUser creates a new user account
	// @Router /users [post]
	e.POST("/users", func(c *gin.Context) {
		var req core.CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		response := s.CreateUser(req)
		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}
		c.JSON(response.Code, response.Meta)
	})

	// GetUser retrieves a user by tag
	// @Router /users/:tag [get]
	e.GET("/users/:tag", func(c *gin.Context) {
		tag := c.Param("tag")
		response := s.GetUser(tag)
		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}
		c.JSON(response.Code, response.Meta)
	})

	// InitVerification starts the identity verification process
	// @Router /verification/session [post]
	e.POST("/verification/session", func(c *gin.Context) {
		var req core.VerifyIdentityRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		response := s.InitVerification(req)
		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}
		c.JSON(response.Code, response.Meta)
	})

	// Webhook for identity provider
	// @Router /webhooks/identity [post]
	e.POST("/webhooks/identity", func(c *gin.Context) {
		var req core.IdentityWebhookRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		// In real world, verify signature here

		response := s.HandleIdentityWebhook(req)
		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}
		c.JSON(response.Code, response.Meta)
	})

	// Example protected route: Request High Limits
	// Requires KYC Level 2 (Full Verified)
	e.POST("/users/request-high-limits", RequireKYC(2, s), func(c *gin.Context) {
		// This user has passed middleware, so they are KYC Level 2
		c.JSON(http.StatusOK, gin.H{
			"message": "High limits application received (Demo: You are verified!)",
		})
	})

	// Link Funding Source
	// @Router /wallets/funding-sources [post]
	e.POST("/wallets/funding-sources", func(c *gin.Context) {
		var req core.LinkFundingSourceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		response := s.LinkFundingSource(req)
		if response.Error {
			c.JSON(response.Code, gin.H{"message": response.Meta.Message})
			return
		}
		c.JSON(response.Code, response.Meta)
	})

	// Deposit Funds
	// @Router /wallets/deposit [post]
	e.POST("/wallets/deposit", func(c *gin.Context) {
		var req core.DepositRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		response := s.Deposit(req)
		if response.Error {
			c.JSON(response.Code, gin.H{"message": response.Meta.Message})
			return
		}
		c.JSON(response.Code, response.Meta)
	})

	// Add Friend
	// @Router /users/friends [post]
	e.POST("/users/friends", func(c *gin.Context) {
		var req core.CreateFriendshipRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		response := s.AddFriend(req)
		if response.Error {
			c.JSON(response.Code, gin.H{"message": response.Meta.Message})
			return
		}
		c.JSON(response.Code, response.Meta)
	})
}
