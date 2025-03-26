package controller

import (
	"go-gorm/dto"
	"go-gorm/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddProductController(router *gin.Engine, db *gorm.DB) {
	routes := router.Group("/products")
	{
		routes.POST("/add-products", func(c *gin.Context) {
			addProductToCart(c, db)
		})
		routes.GET("/search-products", func(c *gin.Context) {
			searchProduct(c, db)
		})

	}
}

func addProductToCart(c *gin.Context, db *gorm.DB) {
	var req dto.AddProductToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
func searchProduct(c *gin.Context, db *gorm.DB) {
	var req dto.SearchProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var products []model.Product
	query := db.Model(&model.Product{})

	if req.ProductName != "" {
		query = query.Where("product_name LIKE ?", "%"+req.ProductName+"%")
	}
	if req.MinPrice > 0 {
		query = query.Where("CAST(price AS FLOAT) >= ?", req.MinPrice)
	}
	if req.MaxPrice > 0 {
		query = query.Where("CAST(price AS FLOAT) <= ?", req.MaxPrice)
	}

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}
