package controller

import (
	"go-gorm/dto"
	"go-gorm/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ShowCartController(router *gin.Engine, db *gorm.DB) {
	routes := router.Group("/show-cart")
	{
		routes.GET("/", func(c *gin.Context) {
			showCart(c, db)
		})
	}
}

func showCart(c *gin.Context, db *gorm.DB) {
	customerID := c.DefaultQuery("customer_id", "")

	var carts []model.Cart
	if err := db.Where("customer_id = ?", customerID).Find(&carts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carts"})
		return
	}

	var cartDetails []dto.CartDTO

	for _, cart := range carts {
		var cartItems []model.CartItem
		if err := db.Where("cart_id = ?", cart.CartID).Find(&cartItems).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart items"})
			return
		}

		var items []dto.CartItemDTO

		for _, cartItem := range cartItems {
			var product model.Product
			if err := db.Where("product_id = ?", cartItem.ProductID).First(&product).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
				return
			}

			productPrice, err := strconv.ParseFloat(product.Price, 64)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid product price"})
				return
			}

			totalPrice := float64(cartItem.Quantity) * productPrice
			items = append(items, dto.CartItemDTO{
				ProductName: product.ProductName,
				Quantity:    cartItem.Quantity,
				Price:       productPrice,
				TotalPrice:  totalPrice,
			})
		}

		cartDetails = append(cartDetails, dto.CartDTO{
			CartName: cart.CartName,
			Items:    items,
		})
	}

	c.JSON(http.StatusOK, gin.H{"carts": cartDetails})
}
