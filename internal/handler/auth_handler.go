package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/your-github-name/hospital-api/internal/service"
)

type AuthHandler struct {
	Service   service.StaffService
	JWTSecret string
}

func (h AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username   string `json:"username" binding:"required"`
		Password   string `json:"password" binding:"required"`
		HospitalID int    `json:"hospital"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	staff, err := h.Service.Create(c.Request.Context(), req.Username, req.Password, req.HospitalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": staff.ID})
}

func (h AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	staff, err := h.Service.Authenticate(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sid": staff.ID,
		"hid": staff.HospitalID,
		"exp": time.Now().Add(72 * time.Hour).Unix(),
	})
	t, _ := token.SignedString([]byte(h.JWTSecret))
	c.JSON(http.StatusOK, gin.H{"token": t})
}
