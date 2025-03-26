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
		routes.PUT("/update-address", func(c *gin.Context) {
			UpdateAddressUser(c, db)
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
func UpdateAddressUser(c *gin.Context, db *gorm.DB) {
	var req dto.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Model(&model.Customer{}).Where("customer_id = ?", req.CustomerID).Update("address", req.Address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address updated successfully"})
}
