package router

import (
	"net/http"

	"github.com/dimasbaguspm/penster/internal/application/service"
	"github.com/dimasbaguspm/penster/internal/interface/handler"
	"github.com/dimasbaguspm/penster/internal/interface/middleware"
)

// Router holds all HTTP handlers and provides route registration
type Router struct {
	healthHandler     *handler.HealthHandler
	accountHandler    *handler.AccountHandler
	categoryHandler  *handler.CategoryHandler
}

// NewRouter creates a new Router with all handlers
func NewRouter(
	healthHandler *handler.HealthHandler,
	accountSvc *service.AccountService,
	categorySvc *service.CategoryService,
) *Router {
	return &Router{
		healthHandler:    healthHandler,
		accountHandler:   handler.NewAccountHandler(accountSvc),
		categoryHandler:  handler.NewCategoryHandler(categorySvc),
	}
}

// Routes returns an http.Handler with all routes registered
func (r *Router) Routes() http.Handler {
	mux := http.NewServeMux()

	// Health endpoints
	mux.HandleFunc("GET /health", r.healthHandler.Health)
	mux.HandleFunc("GET /ready", r.healthHandler.Ready)

	// Account endpoints
	mux.HandleFunc("GET /accounts", r.accountHandler.List)
	mux.HandleFunc("POST /accounts", r.accountHandler.Create)
	mux.HandleFunc("GET /accounts/{id}", r.accountHandler.Get)
	mux.HandleFunc("PUT /accounts/{id}", r.accountHandler.Update)
	mux.HandleFunc("DELETE /accounts/{id}", r.accountHandler.Delete)

	// Category endpoints
	mux.HandleFunc("GET /categories", r.categoryHandler.List)
	mux.HandleFunc("POST /categories", r.categoryHandler.Create)
	mux.HandleFunc("GET /categories/{id}", r.categoryHandler.Get)
	mux.HandleFunc("PUT /categories/{id}", r.categoryHandler.Update)
	mux.HandleFunc("DELETE /categories/{id}", r.categoryHandler.Delete)

	// Apply middleware chain
	handlerChain := middleware.Logging(mux)
	handlerChain = middleware.Recovery(handlerChain)

	return handlerChain
}
