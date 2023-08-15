package main

// Import package
import (
	"InvoltaTask/internal/config"
	"InvoltaTask/internal/handlers"
	"InvoltaTask/internal/models"

	"github.com/gin-gonic/gin" // Framework "Gin"
)

func main() {
	// router with default middleware installed
	router := gin.Default()

	models.ConnectDatabase()

	// index route
	router.POST("/mans", handlers.CreateMan)
	router.GET("/mans", handlers.FindMans)
	router.GET("/mans/:id", handlers.FindMan)
	router.PATCH("/mans/:id", handlers.UpdateMan)
	router.DELETE("/mans/:id", handlers.DeleteMan)

	// Read server port from config
	config, err_config := config.NewConfig(config.YamlPath)
	if err_config != nil {
		panic(err_config)
	}

	// run the server on port
	router.Run( /*"localhost:"*/ "0.0.0.0:" + config.Server.Port)
}
