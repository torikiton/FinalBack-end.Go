package controller

import (
	"go-gorm/dto"
	"go-gorm/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func LoginShowPersonController(router *gin.Engine, db *gorm.DB) {
	routes := router.Group("/auth")
	{
		routes.POST("/login", func(c *gin.Context) {
			loginUser(c, db)
		})
	}
}

func loginUser(c *gin.Context, db *gorm.DB) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var customer model.Customer
	if err := db.Where("email = ?", req.Email).First(&customer).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
		return
	}

	if req.Password != customer.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	response := dto.LoginResponse{}
	copier.Copy(&response, &customer)
	c.JSON(http.StatusOK, response)
}
