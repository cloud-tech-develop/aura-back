package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/cloud-tech-develop/aura-back/application/enterprise"
	appHandlers "github.com/cloud-tech-develop/aura-back/application/handlers"
	domainEnterprise "github.com/cloud-tech-develop/aura-back/domain/enterprise"
	domainErrors "github.com/cloud-tech-develop/aura-back/domain/errors"
	"github.com/cloud-tech-develop/aura-back/infrastructure/messaging/memory"
	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/tenant"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type MigratorAdapter struct {
	manager *tenant.Manager
}

func (m *MigratorAdapter) RunMigrations(e *domainEnterprise.Enterprise, passwordHash string) error {
	te := &tenant.Enterprise{
		Name:           e.Name,
		CommercialName: e.CommercialName,
		Slug:           e.Slug,
		SubDomain:      e.SubDomain,
		Email:          e.Email,
		DV:             e.DV,
		Phone:          e.Phone,
		MunicipalityID: e.MunicipalityID,
		Municipality:   e.Municipality,
		Status:         e.Status,
	}
	return m.manager.CreateEnterprise(context.Background(), te, passwordHash)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")

	database, err := db.New(dsn)
	if err != nil {
		log.Fatal("DB:", err)
	}
	defer database.Close()

	eventBus := memory.NewMemoryEventBus(100, 5)
	if err := eventBus.Start(); err != nil {
		log.Fatalf("Failed to start event bus: %v", err)
	}
	defer eventBus.Stop()

	loggerHandler := appHandlers.NewLoggerHandler("logs")
	if err := eventBus.Subscribe(domainEnterprise.EventEnterpriseCreated, loggerHandler); err != nil {
		log.Printf("Failed to subscribe logger handler: %v", err)
	}
	if err := eventBus.Subscribe(domainEnterprise.EventEnterpriseUpdated, loggerHandler); err != nil {
		log.Printf("Failed to subscribe logger handler: %v", err)
	}
	if err := eventBus.Subscribe(domainEnterprise.EventEnterpriseDeleted, loggerHandler); err != nil {
		log.Printf("Failed to subscribe logger handler: %v", err)
	}

	tenantManager := tenant.NewManager(database.DB)
	if err := tenantManager.MigratePublic(); err != nil {
		log.Fatal("MigratePublic:", err)
	}

	migrator := &MigratorAdapter{manager: tenantManager}
	svc := enterprise.NewServiceWithDB(database.DB, eventBus, migrator)

	r := gin.Default()

	r.POST("/login", tenant.Login(database.DB))

	r.POST("/enterprises", func(c *gin.Context) {
		var req struct {
			Password       string `json:"password" binding:"required"`
			Name           string `json:"name" binding:"required"`
			CommercialName string `json:"commercial_name"`
			Slug           string `json:"slug" binding:"required"`
			SubDomain      string `json:"sub_domain"`
			Email          string `json:"email" binding:"required"`
			DV             string `json:"dv"`
			Phone          string `json:"phone"`
			MunicipalityID string `json:"municipality_id"`
			Municipality   string `json:"municipality"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": domainErrors.ErrFieldsRequired.Error()})
			return
		}

		hash, err := tenant.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error al encriptar password"})
			return
		}

		e := &domainEnterprise.Enterprise{
			Name:           req.Name,
			CommercialName: req.CommercialName,
			Slug:           req.Slug,
			SubDomain:      req.SubDomain,
			Email:          req.Email,
			DV:             req.DV,
			Phone:          req.Phone,
			MunicipalityID: req.MunicipalityID,
			Municipality:   req.Municipality,
		}

		if e.Name == "" || e.Slug == "" || e.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": domainErrors.ErrFieldsRequired.Error()})
			return
		}
		if err := svc.Create(c.Request.Context(), e, hash); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "enterprise " + e.Slug + " creado"})
	})

	r.Use(tenant.AuthMiddleware())

	r.GET("/enterprises", func(c *gin.Context) {
		enterprises, err := svc.List(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, enterprises)
	})

	r.GET("/enterprises/:slug", func(c *gin.Context) {
		slug := c.Param("slug")
		ent, err := svc.GetBySlug(c.Request.Context(), slug)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "enterprise no encontrado"})
			return
		}
		c.JSON(http.StatusOK, ent)
	})

	r.GET("/me", func(c *gin.Context) {
		slug, _ := c.Get(tenant.TenantKey)
		enterpriseID, _ := c.Get("enterprise_id")
		email, _ := c.Get("email")
		c.JSON(http.StatusOK, gin.H{
			"enterprise_id": enterpriseID,
			"slug":          slug,
			"email":         email,
		})
	})

	log.Println("servidor en :" + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
