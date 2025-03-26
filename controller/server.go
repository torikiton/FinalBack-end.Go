package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StartServer(db *gorm.DB) {
	// Set Release Mode
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// load controller
	LoginShowPersonController(router, db)
	// PersonController(router)
	router.Run()
}
