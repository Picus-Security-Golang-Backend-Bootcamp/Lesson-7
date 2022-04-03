package api

import (
	"example.com/with_gin/pkg/middleware"
	"log"
	"os"

	authApi "example.com/with_gin/internal/api/auth"
	cityApi "example.com/with_gin/internal/api/city"
	"example.com/with_gin/internal/city"
	"example.com/with_gin/internal/config"
	"example.com/with_gin/pkg/database_handler"
	"github.com/gin-gonic/gin"
)

var AppConfig = &config.Configuration{}

func RegisterHandlers(r *gin.Engine) {
	cfgFile := "./config/location." + os.Getenv("ENV") + ".yaml"
	AppConfig, err := config.GetAllConfigValues(cfgFile)
	if err != nil {
		log.Fatalf("Failed to read config file. %v", err.Error())
	}

	db := database_handler.NewMySQLDB(AppConfig.DatabaseSettings.DatabaseURI)
	cityRepository := city.NewCityRepository(db)
	cityService := city.NewCityService(*cityRepository)

	cityController := cityApi.NewCityController(cityService)
	authController := authApi.NewAuthController(AppConfig)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	cityGroup := r.Group("/city")
	cityGroup.GET("/:name/:code", func(c *gin.Context) {
		name := c.Param("name")
		code := c.Param("code")
		c.JSON(200, gin.H{"name": name, "code": code})
	})
	cityGroup.POST("", cityController.CreateCity)
	cityGroup.GET("", cityController.GetAllCities)
	cityGroup.GET("/get", cityController.GetQueryString)

	authGroup := r.Group("/auth")

	authGroup.POST("/login", authController.Login)
	authGroup.GET("/decode", middleware.AuthMiddleware(AppConfig.JwtSettings.SecretKey), middleware.AuthMiddleware(AppConfig.JwtSettings.SecretKey), authController.VerifyToken)

}
