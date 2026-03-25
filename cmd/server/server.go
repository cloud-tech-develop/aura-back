package server

import (
	"database/sql"

	"github.com/cloud-tech-develop/aura-back/modules/cart"
	"github.com/cloud-tech-develop/aura-back/modules/enterprise"
	"github.com/cloud-tech-develop/aura-back/modules/inventory"
	"github.com/cloud-tech-develop/aura-back/modules/invoices"
	"github.com/cloud-tech-develop/aura-back/modules/payments"
	"github.com/cloud-tech-develop/aura-back/modules/products"
	"github.com/cloud-tech-develop/aura-back/modules/reports"
	"github.com/cloud-tech-develop/aura-back/modules/sales"
	"github.com/cloud-tech-develop/aura-back/modules/sync"
	"github.com/cloud-tech-develop/aura-back/modules/third-parties"
	"github.com/cloud-tech-develop/aura-back/modules/users"
	"github.com/cloud-tech-develop/aura-back/shared/response"
	"github.com/cloud-tech-develop/aura-back/tenant"
	"github.com/gin-gonic/gin"
)

// Server holds the router and shared dependencies.
type Server struct {
	router    *gin.Engine
	db        *sql.DB
	tenantMgr *tenant.Manager
}

// NewServer creates and configures the Gin router with global middleware.
func NewServer(db *sql.DB, tenantMgr *tenant.Manager) *Server {
	r := gin.Default()

	return &Server{
		router:    r,
		db:        db,
		tenantMgr: tenantMgr,
	}
}

// RegisterModules mounts all module routes onto the router.
// All modules receive:
//   - a public group (no auth)
//   - a protected group (behind AuthMiddleware)
//
// To add a new module, call its Register() here.
func (s *Server) RegisterModules(
	enterpriseH *enterprise.Handler,
	userH *users.Handler,
	productsH *products.Handler,
	cartH *cart.Handler,
	salesH *sales.Handler,
	paymentsH *payments.Handler,
	invoicesH *invoices.Handler,
	reportsH *reports.Handler,
	thirdPartiesH *thirdparties.Handler,
	inventoryH *inventory.Handler,
	syncH *sync.Handler,
) {
	// ── Auth ─────────────────────────────────────────────────────────────────
	s.router.POST("/login", tenant.Login(s.db))

	// ── Static Files ──────────────────────────────────────────────────────────
	s.router.Static("/static", "./static")
	s.router.GET("/download/offline-pos", func(c *gin.Context) {
		c.FileAttachment("./static/bin/aura-pos-offline.exe", "aura-pos-offline.exe")
	})

	// ── Public routes (no middleware) ─────────────────────────────────────────
	public := s.router.Group("/")

	// ── Protected routes (JWT required) ───────────────────────────────────────
	protected := s.router.Group("/")
	protected.Use(tenant.AuthMiddleware())

	// ── Well-known protected routes ───────────────────────────────────────────
	protected.GET("/me", s.meHandler)

	// ── Feature Modules ───────────────────────────────────────────────────────
	enterprise.Register(public, protected, enterpriseH)
	users.Register(public, protected, userH)
	products.Register(public, protected, productsH)
	cart.Register(public, protected, cartH)
	sales.Register(public, protected, salesH)
	payments.Register(public, protected, paymentsH)
	invoices.Register(public, protected, invoicesH)
	reports.Register(public, protected, reportsH)
	thirdparties.Register(public, protected, thirdPartiesH)
	inventory.Register(public, protected, inventoryH)
	sync.Register(public, protected, syncH)
}

// Run starts the HTTP server on the given address.
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

// meHandler returns the authenticated user's context info.
func (s *Server) meHandler(c *gin.Context) {
	response.OK(c, gin.H{
		"user_id":       c.GetInt64("user_id"),
		"enterprise_id": c.GetInt64("enterprise_id"),
		"slug":          c.GetString(string(tenant.TenantKey)),
		"email":         c.GetString("email"),
		"roles":         c.GetStringSlice("roles"),
	})
}
