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

	// ค้นหาสินค้าจาก product_id
	var product model.Product
	if err := db.Where("product_id = ?", req.ProductID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// ค้นหารถเข็นของลูกค้าโดยใช้ customer_id และ cart_name
	var cart model.Cart
	err := db.Where("customer_id = ? AND cart_name = ?", req.CustomerID, req.CartName).First(&cart).Error
	if err != nil {
		// ถ้าไม่พบรถเข็น, สร้างรถเข็นใหม่
		cart = model.Cart{
			CustomerID: req.CustomerID,
			CartName:   req.CartName,
		}
		if err := db.Create(&cart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
			return
		}
	}

	// ตรวจสอบว่าในรถเข็นมีสินค้านั้นอยู่แล้วหรือไม่
	var cartItem model.CartItem
	err = db.Where("cart_id = ? AND product_id = ?", cart.CartID, req.ProductID).First(&cartItem).Error
	if err == nil {
		// หากพบสินค้าในรถเข็นแล้ว, เพิ่มจำนวนสินค้า
		cartItem.Quantity += req.Quantity
		if err := db.Save(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product quantity updated"})
		return
	}

	// ถ้าไม่พบสินค้าภายในรถเข็น, เพิ่มสินค้าใหม่
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
