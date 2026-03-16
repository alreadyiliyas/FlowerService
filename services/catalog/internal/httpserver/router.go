package httpserver

import (
	"github.com/gorilla/mux"
	"github.com/ilyas/flower/services/catalog/internal/httpserver/handlers"
	"github.com/ilyas/flower/services/catalog/internal/httpserver/middleware"
	usecaseCateg "github.com/ilyas/flower/services/catalog/internal/usecase/categories"
	usecaseProd "github.com/ilyas/flower/services/catalog/internal/usecase/products"
)

// newRouter настраивает маршруты HTTP-сервера.
func newRouter(cu usecaseCateg.UsecaseCategories, pu usecaseProd.ProductUsecase, jwtSecret string) *mux.Router {
	router := mux.NewRouter()

	healthHandler := handlers.NewHealthHandler()
	healthRouter := router.PathPrefix("/health").Subrouter()
	healthRouter.HandleFunc("/live", healthHandler.Live).Methods("GET")
	healthRouter.HandleFunc("/ready", healthHandler.Ready).Methods("GET")

	categoriesHandler := handlers.NewCategoriesHandler(cu)
	productsHandler := handlers.NewProductsHandler(pu)

	catalogRouter := router.PathPrefix("/v1/catalog").Subrouter()
	catalogRouter.HandleFunc("/products", productsHandler.ListProducts).Methods("GET")
	catalogRouter.HandleFunc("/products/{id}", productsHandler.GetProduct).Methods("GET")
	catalogRouter.HandleFunc("/categories", categoriesHandler.ListCategories).Methods("GET")
	catalogRouter.HandleFunc("/categories/{id}", categoriesHandler.GetCategory).Methods("GET")

	protected := catalogRouter.NewRoute().Subrouter()
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	protected.Use(middleware.RequireRoles("seller", "moderator"))
	protected.HandleFunc("/products", productsHandler.CreateProduct).Methods("POST")
	protected.HandleFunc("/products/{id}", productsHandler.UpdateProduct).Methods("PATCH")
	protected.HandleFunc("/products/{id}", productsHandler.DeleteProduct).Methods("DELETE")

	catAdmin := catalogRouter.NewRoute().Subrouter()
	catAdmin.Use(middleware.AuthMiddleware(jwtSecret))
	catAdmin.Use(middleware.RequireRoles("seller", "moderator"))
	catAdmin.HandleFunc("/categories", categoriesHandler.CreateCategory).Methods("POST")
	catAdmin.HandleFunc("/categories/{id}", categoriesHandler.UpdateCategory).Methods("PATCH")
	catAdmin.HandleFunc("/categories/{id}", categoriesHandler.DeleteCategory).Methods("DELETE")

	return router
}
