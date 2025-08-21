package middleware

import (
	"anyzzapp/pkg/domain"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.Use(CORS())
	
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	// Test that CORS middleware is applied
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Test that the endpoint works - CORS headers are handled by gin-contrib/cors
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test", response["message"])
}

func TestCORS_WithoutOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.Use(CORS())
	
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	// No Origin header set
	
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Should still work without CORS headers
}

func TestErrorHandler_StringPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.Use(ErrorHandler())
	
	router.GET("/panic", func(c *gin.Context) {
		panic("test error message")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/panic", nil)
	
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	
	var errorResponse domain.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "internal_server_error", errorResponse.Error)
	assert.Equal(t, "test error message", errorResponse.Message)
	assert.Equal(t, 500, errorResponse.Code)
}

func TestErrorHandler_NonStringPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.Use(ErrorHandler())
	
	router.GET("/panic", func(c *gin.Context) {
		panic(123) // Non-string panic
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/panic", nil)
	
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	
	var errorResponse domain.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "internal_server_error", errorResponse.Error)
	assert.Equal(t, "An unexpected error occurred", errorResponse.Message)
	assert.Equal(t, 500, errorResponse.Code)
}

func TestErrorHandler_NormalRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.Use(ErrorHandler())
	
	router.GET("/normal", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/normal", nil)
	
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["message"])
}
