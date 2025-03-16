package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/getsentry/sentry-go"

	"k8s/metrics"
	"k8s/websocket"
)

type Message struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

func main() {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://035fe60ce0f9f9e6b0d9378abd67a269@o4508938459086848.ingest.de.sentry.io/4508938467475536",
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&Message{})
	if err != nil {
		panic(err)
	}

	ginEngine := gin.Default()

	// ginEngine.Use(sentrygin.New(sentrygin.Options{}))

	ginEngine.Use(metrics.PrometheusMiddleware())
	metrics.SetupPrometheusEndpoint(ginEngine)

	ginEngine.POST("/messages", func(c *gin.Context) {
		var message Message
		if err := c.ShouldBindJSON(&message); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&message).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, message)
	})

	ginEngine.GET("/messages", func(c *gin.Context) {
		var messages []Message
		if err := db.Find(&messages).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, messages)
	})

	ginEngine.GET("/error", func(c *gin.Context) {
		c.JSON(500, gin.H{"error": "This is a test error"})
	})

	websocket.SetupWSRoutes(ginEngine)
	ginEngine.Run(":8080")
}
