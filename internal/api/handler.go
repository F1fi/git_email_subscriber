package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"git_email_subscriber/internal/subscriptions"
)

type Handler struct {
	service subscriptions.ISubscriptionService
}

func NewHandler(s subscriptions.ISubscriptionService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	api.POST("/subscribe", h.Subscribe)
	api.GET("/confirm/:token", h.Confirm)
	api.GET("/unsubscribe/:token", h.Unsubscribe)
	api.GET("/subscriptions", h.GetSubscriptions)
}

func (h *Handler) Subscribe(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
		Repo  string `json:"repo"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.service.Subscribe(c, req.Email, req.Repo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "confirmation email sent",
		"token":   sub.ConfirmToken,
	})
}

func (h *Handler) Confirm(c *gin.Context) {
	token := c.Param("token")

	err := h.service.Confirm(c, token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "confirmed"})
}

func (h *Handler) Unsubscribe(c *gin.Context) {
	token := c.Param("token")

	err := h.service.Unsubscribe(c, token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unsubscribed"})
}

func (h *Handler) GetSubscriptions(c *gin.Context) {
	email := c.Query("email")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email required"})
		return
	}

	subs, err := h.service.GetSubscriptions(c, email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed"})
		return
	}

	c.JSON(http.StatusOK, subs)
}