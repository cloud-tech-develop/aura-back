package main

import (
	"context"
	"log"
	"os"

	"github.com/cloud-tech-develop/aura-back/cmd/server"
	"github.com/cloud-tech-develop/aura-back/infrastructure/messaging/memory"
	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/modules/cart"
	"github.com/cloud-tech-develop/aura-back/modules/enterprise"
	"github.com/cloud-tech-develop/aura-back/modules/inventory"
	"github.com/cloud-tech-develop/aura-back/modules/invoices"
	"github.com/cloud-tech-develop/aura-back/modules/payments"
	"github.com/cloud-tech-develop/aura-back/modules/products"
	"github.com/cloud-tech-develop/aura-back/modules/reports"
	"github.com/cloud-tech-develop/aura-back/modules/sales"
	"github.com/cloud-tech-develop/aura-back/modules/third-parties"
	"github.com/cloud-tech-develop/aura-back/modules/users"
	"github.com/cloud-tech-develop/aura-back/tenant"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// ── Database ────────────────────────────────────────────────────────────
	database, err := db.New(dsn)
	if err != nil {
		log.Fatal("DB:", err)
	}
	defer database.Close()

	// ── Event Bus ───────────────────────────────────────────────────────────
	eventBus := memory.NewMemoryEventBus(100, 5)
	if err := eventBus.Start(); err != nil {
		log.Fatalf("Failed to start event bus: %v", err)
	}
	defer eventBus.Stop()

	// ── Tenant Manager & Public Migrations ──────────────────────────────────
	tenantMgr := tenant.NewManager(database.DB)
	if err := tenantMgr.MigratePublic(); err != nil {
		log.Fatal("MigratePublic:", err)
	}

	// ── Modules ──────────────────────────────────────────────────────────────
	// Enterprise module
	enterpriseMigrator := &enterpriseMigratorAdapter{manager: tenantMgr}
	enterpriseSvc := enterprise.NewService(database.DB, eventBus, enterpriseMigrator)
	enterpriseHandler := enterprise.NewHandler(enterpriseSvc, tenantMgr)

	// Logging (Enterprise)
	enterpriseLogger := enterprise.NewLoggerHandler("logs")
	_ = eventBus.Subscribe(enterprise.EventCreated, enterpriseLogger)
	_ = eventBus.Subscribe(enterprise.EventUpdated, enterpriseLogger)
	_ = eventBus.Subscribe(enterprise.EventDeleted, enterpriseLogger)

	// Users module
	usersSvc := users.NewService(database.DB, eventBus)
	usersHandler := users.NewHandler(usersSvc)

	// Logging (Users)
	usersLogger := users.NewLoggerHandler("logs")
	_ = eventBus.Subscribe(users.EventCreated, usersLogger)
	_ = eventBus.Subscribe(users.EventUpdated, usersLogger)
	// _ = eventBus.Subscribe(users.EventDeleted, usersLogger) // If implemented

	// Products module
	productsHandler := products.NewHandler(database.DB)

	// Cart module
	cartHandler := cart.NewHandler(database.DB)

	// Sales module
	salesHandler := sales.NewHandler(database.DB)

	// Payments module
	paymentsHandler := payments.NewHandler(database.DB)

	// Invoices module
	invoicesHandler := invoices.NewHandler(database.DB)

	// Reports module
	reportsHandler := reports.NewHandler(database.DB)

	// Third Parties module
	thirdPartiesHandler := thirdparties.NewHandler(database.DB)

	// Inventory module
	inventoryHandler := inventory.NewHandler(database.DB)

	// ── HTTP Server ──────────────────────────────────────────────────────────
	srv := server.NewServer(database.DB, tenantMgr)
	srv.RegisterModules(enterpriseHandler, usersHandler, productsHandler, cartHandler, salesHandler, paymentsHandler, invoicesHandler, reportsHandler, thirdPartiesHandler, inventoryHandler)

	log.Println("servidor en :" + port)
	if err := srv.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

// enterpriseMigratorAdapter adapts tenant.Manager to enterprise.Migrator.
type enterpriseMigratorAdapter struct {
	manager *tenant.Manager
}

func (a *enterpriseMigratorAdapter) RunMigrations(ctx context.Context, e *enterprise.Enterprise, passwordHash string) error {
	te := &tenant.Enterprise{
		Name:           e.Name,
		CommercialName: e.CommercialName,
		Slug:           e.Slug,
		SubDomain:      e.SubDomain,
		Email:          e.Email,
		Document:       e.Document.String(),
		DV:             e.DV,
		Phone:          e.Phone,
		MunicipalityID: e.MunicipalityID,
		Municipality:   e.Municipality,
		Status:         e.Status,
	}
	return a.manager.CreateEnterprise(ctx, te, passwordHash)
}
