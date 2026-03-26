package httpserver

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ilyas/flower/services/auth/internal/httpserver/handlers"
	"github.com/ilyas/flower/services/auth/internal/httpserver/middleware"
	authusecase "github.com/ilyas/flower/services/auth/internal/usecase/auth"
	userusecase "github.com/ilyas/flower/services/auth/internal/usecase/user"
)

// newRouter настраивает маршруты HTTP-сервера.
func newRouter(authUC authusecase.UsecaseAuth, userUC userusecase.UsecaseUser, jwtSecret string) *mux.Router {
	router := mux.NewRouter()

	healthHandler := handlers.NewHealthHandler()
	healthRouter := router.PathPrefix("/health").Subrouter()
	healthRouter.HandleFunc("/live", healthHandler.Live).Methods("GET")
	healthRouter.HandleFunc("/ready", healthHandler.Ready).Methods("GET")

	authHandler := handlers.NewAuthHandler(authUC)
	userHandler := handlers.NewUserHandler(userUC)

	// Маршруты для аутентификации
	authRouter := router.PathPrefix("/v1/auth").Subrouter()
	authRouter.HandleFunc("/registration", authHandler.Register).Methods("POST")
	authRouter.HandleFunc("/set_password", authHandler.SetPassword).Methods("POST")
	authRouter.HandleFunc("/update_password/request_code", authHandler.SendConfirmUpdatePassword).Methods("POST")
	authRouter.HandleFunc("/update_password/confirm", authHandler.ConfirmUpdatePassword).Methods("PATCH")
	authRouter.HandleFunc("/confirm_code", authHandler.VerifyAccount).Methods("POST")
	authRouter.HandleFunc("/login", authHandler.Login).Methods("POST")
	authRouter.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")
	protectedAuth := authRouter.NewRoute().Subrouter()
	protectedAuth.Use(middleware.AuthMiddleware(jwtSecret, authUC))
	protectedAuth.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	protectedAuth.HandleFunc("/logout_all", authHandler.LogoutAll).Methods("POST")

	userRouter := router.PathPrefix("/v1/user").Subrouter()
	userRouter.Use(middleware.AuthMiddleware(jwtSecret, authUC))
	userRouter.HandleFunc("/me", userHandler.GetUserInfo).Methods("GET")
	userRouter.HandleFunc("/me/update", userHandler.UpdateUserInfo).Methods("PATCH")
	userRouter.HandleFunc("/me/avatar", userHandler.UploadAvatar).Methods("POST")
	userRouter.HandleFunc("/me/delete", userHandler.DeleteUser).Methods("DELETE")
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	return router
}
