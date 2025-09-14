package middleware

import (
	"net/http"
)

// CORSConfig represents CORS configuration
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
}

// DefaultCORSConfig returns default CORS configuration
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Cookie"},
		AllowCredentials: true,
	}
}

// CORSMiddleware provides CORS functionality
type CORSMiddleware struct {
	config *CORSConfig
}

// NewCORSMiddleware creates a new CORS middleware instance
func NewCORSMiddleware(config *CORSConfig) *CORSMiddleware {
	if config == nil {
		config = DefaultCORSConfig()
	}
	return &CORSMiddleware{
		config: config,
	}
}

// Handler applies CORS headers to HTTP responses
func (c *CORSMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		if len(c.config.AllowOrigins) > 0 {
			// For simplicity, setting the first origin. In production, you might want to check the request origin
			w.Header().Set("Access-Control-Allow-Origin", c.config.AllowOrigins[0])
		}

		if len(c.config.AllowMethods) > 0 {
			methods := ""
			for i, method := range c.config.AllowMethods {
				if i > 0 {
					methods += ", "
				}
				methods += method
			}
			w.Header().Set("Access-Control-Allow-Methods", methods)
		}

		if len(c.config.AllowHeaders) > 0 {
			headers := ""
			for i, header := range c.config.AllowHeaders {
				if i > 0 {
					headers += ", "
				}
				headers += header
			}
			w.Header().Set("Access-Control-Allow-Headers", headers)
		}

		if c.config.AllowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
