package routes

import (
	"net/http"
	"strings"

	"github.com/vocal-tracker/broker-service/config"
	"github.com/vocal-tracker/broker-service/middleware"
)

func NewRouter(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// Create handlers
	authHandler := NewAuthHandler(cfg)
	vocabHandler := NewVocabHandler(cfg)
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	// Register auth routes with /auth prefix to match frontend expectations
	mux.HandleFunc("POST /auth/register", enableCORS(authHandler.Register))
	mux.HandleFunc("POST /auth/login", enableCORS(authHandler.Login))
	mux.HandleFunc("OPTIONS /auth/register", handleOptions)
	mux.HandleFunc("OPTIONS /auth/login", handleOptions)

	// Register vocabulary routes with auth middleware
	mux.Handle("GET /vocab", authMiddleware.RequireAuth(http.HandlerFunc(vocabHandler.GetVocabularies)))
	mux.Handle("POST /vocab", authMiddleware.RequireAuth(http.HandlerFunc(vocabHandler.CreateVocabulary)))
	mux.Handle("PUT /vocab/", authMiddleware.RequireAuth(http.HandlerFunc(vocabUpdateHandler(vocabHandler))))
	mux.Handle("DELETE /vocab/", authMiddleware.RequireAuth(http.HandlerFunc(vocabDeleteHandler(vocabHandler))))

	// OPTIONS for vocab routes
	mux.HandleFunc("OPTIONS /vocab", handleOptions)
	mux.HandleFunc("OPTIONS /vocab/", handleOptions)

	return mux
}

// vocabUpdateHandler handles PUT requests with ID in path
func vocabUpdateHandler(handler *VocabHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/vocab/") || len(strings.TrimPrefix(r.URL.Path, "/vocab/")) == 0 {
			http.Error(w, "Invalid vocabulary ID", http.StatusBadRequest)
			return
		}
		handler.UpdateVocabulary(w, r)
	}
}

// vocabDeleteHandler handles DELETE requests with ID in path
func vocabDeleteHandler(handler *VocabHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/vocab/") || len(strings.TrimPrefix(r.URL.Path, "/vocab/")) == 0 {
			http.Error(w, "Invalid vocabulary ID", http.StatusBadRequest)
			return
		}
		handler.DeleteVocabulary(w, r)
	}
}

// enableCORS adds CORS headers to the response
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// handleOptions handles preflight OPTIONS requests
func handleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusOK)
}
