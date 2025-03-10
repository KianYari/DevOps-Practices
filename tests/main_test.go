package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Message struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

func setupTestRouter() (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	// Use in-memory SQLite for testing
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&Message{})

	r := gin.Default()

	r.POST("/messages", func(c *gin.Context) {
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

	r.GET("/messages", func(c *gin.Context) {
		var messages []Message
		if err := db.Find(&messages).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, messages)
	})

	return r, db
}

func TestPostMessage(t *testing.T) {
	router, _ := setupTestRouter()

	// Create a test message
	message := Message{Content: "Test message"}
	jsonValue, _ := json.Marshal(message)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/messages", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert status code
	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse response body
	var response Message
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Verify response
	assert.NotEqual(t, 0, response.ID)
	assert.Equal(t, "Test message", response.Content)
}

func TestGetMessages(t *testing.T) {
	router, db := setupTestRouter()

	// Add test data
	testMessages := []Message{
		{Content: "First message"},
		{Content: "Second message"},
	}
	for _, msg := range testMessages {
		db.Create(&msg)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/messages", nil)
	router.ServeHTTP(w, req)

	// Assert status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response body
	var response []Message
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Verify response
	assert.Equal(t, 2, len(response))
	assert.Equal(t, "First message", response[0].Content)
	assert.Equal(t, "Second message", response[1].Content)
}

func TestPostMessageInvalidJSON(t *testing.T) {
	router, _ := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/messages", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert status code for bad request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
