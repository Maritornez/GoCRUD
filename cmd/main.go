package main

import (
	"fmt"
	"net/http"

	"github.com/Maritornez/GoCRUD/internal/config"
	"github.com/Maritornez/GoCRUD/internal/handlers"
	"github.com/Maritornez/GoCRUD/internal/storage"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

// Middleware для проверки подключения к базе данных
func dbCheckMiddleware(c *gin.Context) {
	var isDatabaseConnected bool
	if storage.DB.Status().Err != nil {
		isDatabaseConnected = false
	} else {
		isDatabaseConnected = true
	}

	if !isDatabaseConnected {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		c.Abort()
		return
	}
	c.Next()
}

func main() {
	// Загрузка переменных окружения из .env файла, если программа не в докер-контейнере
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Маршрутизаторр включает стандартные middleware: Logger, Recovery
	router := gin.Default()
	router.Use(dbCheckMiddleware)

	storage.ConnectDatabase()
	storage.InitCache()

	storage.InitializeDatabase()

	router.POST("/company", handlers.CreateCompany)
	router.GET("/companies", handlers.FindCompanies) //companies?limit=10&offset=0
	router.GET("/company/:id", handlers.FindCompany)
	router.PATCH("/company/:id", handlers.UpdateCompany)
	router.DELETE("/company/:id", handlers.DeleteCompany)

	router.POST("/man", handlers.CreateMan)
	router.GET("/men", handlers.FindMen) //men?limit=10&offset=0
	router.GET("/man/:id", handlers.FindMan)
	router.PATCH("/man/:id", handlers.UpdateMan)
	router.DELETE("/man/:id", handlers.DeleteMan)

	router.POST("/tip", handlers.CreateTip)
	router.GET("/tips", handlers.FindTips) //tips?limit=10&offset=0
	router.GET("/tip/:id", handlers.FindTip)
	router.PATCH("/tip/:id", handlers.UpdateTip)
	router.DELETE("/tip/:id", handlers.DeleteTip)

	config, err_config := config.NewConfig()
	if err_config != nil {
		panic(err_config)
	}

	router.Run(config.Server.IpAddress + ":" + config.Server.Port)
}
