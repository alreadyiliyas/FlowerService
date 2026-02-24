package httpserver

import (
	"github.com/gorilla/mux"
	"github.com/ilyas/flower/services/auth/internal/httpserver/handlers"
	authusecase "github.com/ilyas/flower/services/auth/internal/usecase/auth"
)

// newRouter настраивает маршруты HTTP-сервера.
func newRouter(authUC authusecase.Usecase) *mux.Router {
	router := mux.NewRouter()

	healthHandler := handlers.NewHealthHandler()
	healthRouter := router.PathPrefix("/health").Subrouter()
	healthRouter.HandleFunc("/live", healthHandler.Live).Methods("GET")
	healthRouter.HandleFunc("/ready", healthHandler.Ready).Methods("GET")

	// Прописать Middleware
	// router.Use(loggingMiddleware)
	// router.Use(corsMiddleware)
	// authRouter.Use(authMiddleware)

	authHandler := handlers.NewAuthHandler(authUC)

	// Маршруты для аутентификации
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", authHandler.Register).Methods("POST")
	authRouter.HandleFunc("/set_password", authHandler.SetPassword).Methods("POST")
	authRouter.HandleFunc("/update_password", authHandler.UpdatePassword).Methods("PATCH")
	authRouter.HandleFunc("/verify", authHandler.VerifyAccount).Methods("POST")
	authRouter.HandleFunc("/login", authHandler.Login).Methods("POST")
	authRouter.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")
	authRouter.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	return router
}
