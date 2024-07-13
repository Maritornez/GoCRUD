package main

import (
	"net/http"

	"github.com/Maritornez/GoCRUD/internal/config"
	"github.com/Maritornez/GoCRUD/internal/context"
	"github.com/Maritornez/GoCRUD/internal/handlers"

	"github.com/gin-gonic/gin"
)

// Middleware для проверки подключения к базе данных
func dbCheckMiddleware(c *gin.Context) {
	if !isDatabaseConnected() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		c.Abort()
		return
	}
	c.Next()
}

// Функция для проверки подключения к базе данных (пример)
func isDatabaseConnected() bool {
	if context.DB.Status().Err != nil {
		return false
	} else {
		return true
	}
}

func main() {
	// маршрутизаторр включает стандартные middleware: Logger, Recovery
	router := gin.Default()
	router.Use(dbCheckMiddleware)

	context.ConnectDatabase()

	router.POST("/man", handlers.CreateMan)
	router.GET("/men", handlers.FindMen) //men?limit=10&offset=2
	router.GET("/man/:id", handlers.FindMan)
	router.PATCH("/man/:id", handlers.UpdateMan)
	router.DELETE("/man/:id", handlers.DeleteMan)

	router.POST("/company", handlers.CreateCompany)
	router.GET("/companies", handlers.FindCompanies) //companies?limit=10&offset=2
	router.GET("/company/:id", handlers.FindCompany)
	router.PATCH("/company/:id", handlers.UpdateCompany)
	router.DELETE("/company/:id", handlers.DeleteCompany)

	config, err_config := config.NewConfig(config.YamlPath)
	if err_config != nil {
		panic(err_config)
	}

	router.Run(config.Server.IpAddress + ":" + config.Server.Port)
}
