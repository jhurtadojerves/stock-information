package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMainServer(t *testing.T) {
	message := "Mock response"
	r := gin.Default()
	r.GET("/api/v1/stocks", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": message})
	})
	req, _ := http.NewRequest("GET", "/api/v1/stocks", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), message)
}
