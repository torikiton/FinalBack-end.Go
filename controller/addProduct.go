package controller

import (
	"go-gorm/dto"
	"go-gorm/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddProductController(router *gin.Engine, db *gorm.DB) {
	routes := router.Group("/add-products")
	{
		routes.POST("/", func(c *gin.Context) {
			addProductToCart(c, db)
		})
	}
}
func addProductToCart(c *gin.Context, db *gorm.DB) {

	var req dto.AddProductToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product model.Product
	if err := db.Where("product_id = ?", req.ProductID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var cart model.Cart
	err := db.Where("customer_id = ? AND cart_name = ?", req.CustomerID, req.CartName).First(&cart).Error
	if err != nil {
		cart = model.Cart{
			CustomerID: req.CustomerID,
			CartName:   req.CartName,
		}
		if err := db.Create(&cart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
			return
		}
	}

	var cartItem model.CartItem
	err = db.Where("cart_id = ? AND product_id = ?", cart.CartID, req.ProductID).First(&cartItem).Error
	if err == nil {
		cartItem.Quantity += req.Quantity
		if err := db.Save(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product quantity updated"})
		return
	}

	cartItem = model.CartItem{
		CartID:    cart.CartID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := db.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added to cart"})
}
