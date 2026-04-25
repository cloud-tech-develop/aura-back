package server

import (
	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/modules/admin/enterprise"
	thirdparties "github.com/cloud-tech-develop/aura-back/modules/admin/third-parties"
	"github.com/cloud-tech-develop/aura-back/modules/admin/users"
	"github.com/cloud-tech-develop/aura-back/modules/catalog/brands"
	"github.com/cloud-tech-develop/aura-back/modules/catalog/categories"
	"github.com/cloud-tech-develop/aura-back/modules/catalog/presentations"
	catalogproducts "github.com/cloud-tech-develop/aura-back/modules/catalog/products"
	"github.com/cloud-tech-develop/aura-back/modules/catalog/units"
	"github.com/cloud-tech-develop/aura-back/modules/offline"
	"github.com/cloud-tech-develop/aura-back/shared/response"
	"github.com/cloud-tech-develop/aura-back/tenant"
	"github.com/gin-gonic/gin"
)

// Server holds the router and shared dependencies.
type Server struct {
	router    *gin.Engine
	db        *db.DB
	tenantMgr *tenant.Manager
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// NewServer creates and configures the Gin router with global middleware.
func NewServer(database *db.DB, tenantMgr *tenant.Manager) *Server {
	r := gin.Default()
	r.Use(corsMiddleware())

	return &Server{
		router:    r,
		db:        database,
		tenantMgr: tenantMgr,
	}
}

// RegisterModules mounts all module routes onto the router.
func (s *Server) RegisterModules(
	enterpriseH *enterprise.Handler,
	userH *users.Handler,
	categoryH *categories.Handler,
	brandH *brands.Handler,
	productH *catalogproducts.Handler,
	presentationH *presentations.Handler,
	thirdPartiesH *thirdparties.Handler,
	unitH *units.Handler,
	offlineH *offline.Handler,
) {
	// Health Check
	s.router.GET("/", func(c *gin.Context) {
		response.OK(c, "Hello, Aura!")
	})

	// Auth
	s.router.POST("/login", tenant.Login(s.db.Wrap(s.db.DB)))

	// Static Files
	s.router.Static("/static", "./static")
	s.router.GET("/download/offline-pos", func(c *gin.Context) {
		c.FileAttachment("./static/bin/aura-pos-offline.exe", "aura-pos-offline.exe")
	})

	// Public routes (no middleware)
	public := s.router.Group("/")

	// Protected routes (JWT required)
	protected := s.router.Group("/")
	protected.Use(tenant.AuthMiddleware())

	// Well-known protected routes
	protected.GET("/me", s.meHandler)

	// Feature Modules
	enterprise.Register(public, protected, enterpriseH)
	users.Register(public, protected, userH)
	categories.Register(public, protected, categoryH)
	brands.Register(public, protected, brandH)
	catalogproducts.Register(public, protected, productH)
	catalogproducts.RegisterProductPresentations(public, protected, productH, presentationH)
	presentations.Register(protected, presentationH)
	thirdparties.Register(public, protected, thirdPartiesH)
	units.Register(public, protected, unitH)

	// Offline sync (only available in offline mode)
	if offlineH != nil {
		offline.Register(public, protected, offlineH)
	}
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
