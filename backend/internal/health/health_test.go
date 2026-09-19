package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthAndReadinessCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/health", HealthCheckHandler("votify-backend"))
	r.GET("/ready", ReadinessCheckHandler("votify-backend"))

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for /health, got %d", w.Code)
	}

	reqReady, _ := http.NewRequest("GET", "/ready", nil)
	wReady := httptest.NewRecorder()
	r.ServeHTTP(wReady, reqReady)
	if wReady.Code != http.StatusOK {
		t.Fatalf("expected status 200 for /ready, got %d", wReady.Code)
	}
}
