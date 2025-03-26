package controller

import "github.com/gin-gonic/gin"

func StartServer() {
	// Set Release Mode
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// load controller
	// DemoController(router)
	// PersonController(router)
	router.Run()
}
