package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// ReverseProxy forwards incoming HTTP requests from the API Gateway to downstream microservices.
func ReverseProxy(targetURL string, trimPrefix string) gin.HandlerFunc {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		panic("Invalid target URL for proxy: " + targetURL)
	}

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)

	// Custom Director to adjust request path
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = parsedURL.Host

		if trimPrefix != "" && strings.HasPrefix(req.URL.Path, trimPrefix) {
			req.URL.Path = strings.TrimPrefix(req.URL.Path, trimPrefix)
			if !strings.HasPrefix(req.URL.Path, "/") {
				req.URL.Path = "/" + req.URL.Path
			}
		}
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"success":false,"error":{"code":"SERVICE_UNAVAILABLE","message":"Downstream service unavailable"}}`))
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
