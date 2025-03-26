package controller

import (
	"go-gorm/dto"
	"go-gorm/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
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
		routes.PUT("/change-password", func(c *gin.Context) {
			changePassword(c, db)
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

	err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(req.Password))
	if err != nil {
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
func changePassword(c *gin.Context, db *gorm.DB) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var customer model.Customer
	if err := db.Where("customer_id = ?", req.CustomerID).First(&customer).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Customer not found"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(req.OldPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Old password is incorrect"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	if err := db.Model(&model.Customer{}).Where("customer_id = ?", req.CustomerID).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
