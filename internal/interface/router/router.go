package router

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/dimasbaguspm/penster/internal/application/service"
	"github.com/dimasbaguspm/penster/internal/interface/handler"
	"github.com/dimasbaguspm/penster/internal/interface/middleware"
)

// Router holds all HTTP handlers and provides route registration
type Router struct {
	healthHandler      *handler.HealthHandler
	accountHandler     *handler.AccountHandler
	categoryHandler    *handler.CategoryHandler
	transactionHandler *handler.TransactionHandler
	draftHandler       *handler.DraftHandler
	reportHandler      *handler.ReportHandler
}

// NewRouter creates a new Router with all handlers
func NewRouter(
	healthHandler *handler.HealthHandler,
	accountSvc *service.AccountService,
	categorySvc *service.CategoryService,
	transactionSvc *service.TransactionService,
	draftSvc *service.DraftService,
	reportSvc *service.ReportService,
) *Router {
	return &Router{
		healthHandler:      healthHandler,
		accountHandler:     handler.NewAccountHandler(accountSvc),
		categoryHandler:   handler.NewCategoryHandler(categorySvc),
		transactionHandler: handler.NewTransactionHandler(transactionSvc),
		draftHandler:       handler.NewDraftHandler(draftSvc),
		reportHandler:      handler.NewReportHandler(reportSvc),
	}
}

// Routes returns an http.Handler with all routes registered
func (r *Router) Routes() http.Handler {
	mux := http.NewServeMux()

	// Health endpoints
	mux.HandleFunc("GET /health", r.healthHandler.Health)
	mux.HandleFunc("GET /ready", r.healthHandler.Ready)

	// Swagger endpoint
	mux.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))

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

	// Transaction endpoints
	mux.HandleFunc("GET /transactions", r.transactionHandler.List)
	mux.HandleFunc("POST /transactions", r.transactionHandler.Create)
	mux.HandleFunc("GET /transactions/{id}", r.transactionHandler.Get)
	mux.HandleFunc("PUT /transactions/{id}", r.transactionHandler.Update)
	mux.HandleFunc("DELETE /transactions/{id}", r.transactionHandler.Delete)

	// Draft endpoints
	mux.HandleFunc("GET /drafts", r.draftHandler.List)
	mux.HandleFunc("POST /drafts", r.draftHandler.Create)
	mux.HandleFunc("GET /drafts/{id}", r.draftHandler.Get)
	mux.HandleFunc("PATCH /drafts/{id}", r.draftHandler.Update)
	mux.HandleFunc("POST /drafts/{id}/confirm", r.draftHandler.Confirm)
	mux.HandleFunc("POST /drafts/{id}/reject", r.draftHandler.Reject)
	mux.HandleFunc("DELETE /drafts/{id}", r.draftHandler.Delete)

	// Report endpoints
	mux.HandleFunc("GET /reports/summary", r.reportHandler.Summary)
	mux.HandleFunc("GET /reports/by-account", r.reportHandler.ByAccount)
	mux.HandleFunc("GET /reports/by-category", r.reportHandler.ByCategory)
	mux.HandleFunc("GET /reports/trends", r.reportHandler.Trends)

	// Apply middleware chain
	handlerChain := middleware.Logging(mux)
	handlerChain = middleware.Recovery(handlerChain)

	return handlerChain
}
