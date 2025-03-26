package controller

import (
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
	// รับ customer_id จาก query parameters
	customerID := c.DefaultQuery("customer_id", "")

	// ค้นหารถเข็นทั้งหมดของลูกค้า
	var carts []model.Cart
	if err := db.Where("customer_id = ?", customerID).Find(&carts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve carts"})
		return
	}

	// สร้างรายการสำหรับแต่ละรถเข็น
	var cartDetails []struct {
		CartName string `json:"cart_name"`
		Items    []struct {
			ProductName string  `json:"product_name"`
			Quantity    int     `json:"quantity"`
			Price       float64 `json:"price"`
			TotalPrice  float64 `json:"total_price"`
		} `json:"items"`
	}

	// ดึงข้อมูลสินค้าและคำนวณราคาของแต่ละรายการในแต่ละรถเข็น
	for _, cart := range carts {
		var cartItems []model.CartItem
		if err := db.Where("cart_id = ?", cart.CartID).Find(&cartItems).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart items"})
			return
		}

		var items []struct {
			ProductName string  `json:"product_name"`
			Quantity    int     `json:"quantity"`
			Price       float64 `json:"price"`
			TotalPrice  float64 `json:"total_price"`
		}

		for _, cartItem := range cartItems {
			var product model.Product
			if err := db.Where("product_id = ?", cartItem.ProductID).First(&product).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
				return
			}

			// แปลง product.Price จาก string เป็น float64
			productPrice, err := strconv.ParseFloat(product.Price, 64)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid product price"})
				return
			}

			totalPrice := float64(cartItem.Quantity) * productPrice
			items = append(items, struct {
				ProductName string  `json:"product_name"`
				Quantity    int     `json:"quantity"`
				Price       float64 `json:"price"`
				TotalPrice  float64 `json:"total_price"`
			}{
				ProductName: product.ProductName,
				Quantity:    cartItem.Quantity,
				Price:       productPrice,
				TotalPrice:  totalPrice,
			})
		}

		cartDetails = append(cartDetails, struct {
			CartName string `json:"cart_name"`
			Items    []struct {
				ProductName string  `json:"product_name"`
				Quantity    int     `json:"quantity"`
				Price       float64 `json:"price"`
				TotalPrice  float64 `json:"total_price"`
			} `json:"items"`
		}{
			CartName: cart.CartName,
			Items:    items,
		})
	}

	// ส่งข้อมูลรถเข็นและรายการสินค้ากลับไปยัง client
	c.JSON(http.StatusOK, gin.H{"carts": cartDetails})
}
