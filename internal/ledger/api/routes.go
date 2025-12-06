package api

import (
	"cashapp/core"
	"cashapp/internal/ledger/service"
	"net/http"

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
}
