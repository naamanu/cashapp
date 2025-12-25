package api

import (
	"cashapp/core"
	"cashapp/internal/ledger/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RegisterPaymentRoutes registers payment-related routes
func RegisterPaymentRoutes(e *gin.Engine, s *service.PaymentService) {
	// SendMoney creates a new payment transaction
	// @Router /payments [post]
	e.POST("/payments", func(c *gin.Context) {
		var req core.CreatePaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		response := s.SendMoney(req)
		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)
	})

	// GetBalance retrieves wallet balance
	// @Router /wallets/:id/balance [get]
	e.GET("/wallets/:id/balance", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid wallet id",
			})
			return
		}

		response := s.GetBalance(id)
		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}
		c.JSON(response.Code, response.Meta)
	})
	// Create Payment Request
	// @Router /payments/requests [post]
	e.POST("/payments/requests", func(c *gin.Context) {
		var req core.CreateRequestDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		response := s.CreateRequest(req)
		if response.Error {
			c.JSON(response.Code, gin.H{"message": response.Meta.Message})
			return
		}
		c.JSON(response.Code, response.Meta)
	})

	// Pay a Payment Request
	// @Router /payments/requests/:id/pay [post]
	e.POST("/payments/requests/:id/pay", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request id"})
			return
		}

		// In real world, we get payer info from context.
		// For now, we assume the authorized user is the PayerID on the request.
		response := s.PayRequest(id, "mock-auth-key")
		if response.Error {
			c.JSON(response.Code, gin.H{"message": response.Meta.Message})
			return
		}
		c.JSON(response.Code, response.Meta)
	})

	// Get Feed (Social Activity)
	// @Router /feed [post]
	e.POST("/feed", func(c *gin.Context) {
		type FeedRequest struct {
			FriendIDs []int `json:"friend_ids"`
		}
		var req FeedRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		response := s.GetFeed(req.FriendIDs)
		if response.Error {
			c.JSON(response.Code, gin.H{"message": response.Meta.Message})
			return
		}
		c.JSON(response.Code, response.Meta)
	})

	// Split Bill
	// @Router /payments/split [post]
	e.POST("/payments/split", func(c *gin.Context) {
		var req core.SplitBillDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		response := s.SplitBill(req)
		if response.Error {
			c.JSON(response.Code, gin.H{"message": response.Meta.Message})
			return
		}
		c.JSON(response.Code, response.Meta)
	})
}
