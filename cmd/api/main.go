package main

import (
	"context"
	"log"
	"os"

	"os/exec"

	"github.com/cloud-tech-develop/aura-back/cmd/server"
	"github.com/cloud-tech-develop/aura-back/infrastructure/messaging/memory"
	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/modules/admin/enterprise"
	thirdparties "github.com/cloud-tech-develop/aura-back/modules/admin/third-parties"
	"github.com/cloud-tech-develop/aura-back/modules/admin/users"
	"github.com/cloud-tech-develop/aura-back/modules/catalog/brands"
	"github.com/cloud-tech-develop/aura-back/modules/catalog/categories"
	"github.com/cloud-tech-develop/aura-back/modules/catalog/presentations"
	catalogProducts "github.com/cloud-tech-develop/aura-back/modules/catalog/products"
	"github.com/cloud-tech-develop/aura-back/modules/catalog/units"
	"github.com/cloud-tech-develop/aura-back/modules/offline"
	"github.com/cloud-tech-develop/aura-back/tenant"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	driver := os.Getenv("DATABASE_DRIVER")
	port := os.Getenv("PORT")

	// Default to SQLite (offline mode) if no DATABASE_URL is set
	if dsn == "" {
		driver = "sqlite"
		dsn = "aura_pos.db"
		port = "8091"
		log.Println("Running in offline mode with SQLite")
	} else if driver == "" {
		driver = "postgres"
		if port == "" {
			port = "8081"
		}
	}

	// Database
	database, err := db.New(driver, dsn)
	if err != nil {
		log.Fatal("DB:", err)
	}

	// Set SQLite mode for tenant manager if offline
	if driver == "sqlite" {
		tenant.SetSQLiteMode(true)
	}

	// Event Bus
	eventBus := memory.NewMemoryEventBus(100, 5)
	if err := eventBus.Start(); err != nil {
		log.Fatalf("Failed to start event bus: %v", err)
	}
	defer eventBus.Stop()

	// Tenant Manager & Public Migrations
	tenantMgr := tenant.NewManager(database.DB)
	if err := tenantMgr.MigratePublic(); err != nil {
		log.Fatal("MigratePublic:", err)
	}

	// Migrate existing tenants in background (PostgreSQL only)
	if driver != "sqlite" {
		go func() {
			log.Println("Migrating existing tenants...")
			if err := tenantMgr.MigrateAll(context.Background()); err != nil {
				log.Printf("MigrateAll: %v\n", err)
				return
			}
			log.Println("All tenants migrated successfully")
		}()
	} else {
		log.Println("Skipping MigrateAll in SQLite mode")
	}

	// Modules
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

	// Catalog modules
	categorySvc := categories.NewService(database.Wrap(database.DB))
	categoryHandler := categories.NewHandler(categorySvc)

	brandSvc := brands.NewService(database.Wrap(database.DB))
	brandHandler := brands.NewHandler(brandSvc)

	unitSvc := units.NewService(database.Wrap(database.DB))
	unitHandler := units.NewHandler(unitSvc)

	productsLogger := catalogProducts.NewLoggerHandler("logs")
	_ = eventBus.Subscribe(catalogProducts.EventCreated, productsLogger)
	_ = eventBus.Subscribe(catalogProducts.EventUpdated, productsLogger)
	_ = eventBus.Subscribe(catalogProducts.EventDeleted, productsLogger)

	// Presentations module (must be initialized before products)
	presSvc := presentations.NewService(database, eventBus)
	presHandler := presentations.NewHandler(presSvc)

	productSvc := catalogProducts.NewService(database, eventBus, presSvc)
	productHandler := catalogProducts.NewHandler(productSvc)

	// Third Parties module
	thirdPartiesHandler := thirdparties.NewHandler(database)

	// Offline module (only in offline mode)
	var offlineHandler *offline.Handler
	if driver == "sqlite" {
		offlineSvc := offline.NewService(database.DB, eventBus)
		offlineHandler = offline.NewHandler(offlineSvc)
	}

	// HTTP Server
	srv := server.NewServer(database.DB, tenantMgr)
	srv.RegisterModules(enterpriseHandler, usersHandler, categoryHandler, brandHandler, productHandler, presHandler, thirdPartiesHandler, unitHandler, offlineHandler)

	log.Println("servidor en :" + port)

	if driver == "postgres" {
		log.Println("Running in production mode")
		go buildOfflineBinary()
	}

	if err := srv.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func buildOfflineBinary() {
	log.Println("Generating offline binary...")

	// Create static directory if it doesn't exist
	_ = os.MkdirAll("static/bin", 0755)

	cmd := exec.Command("go", "build", "-o", "static/bin/aura-pos-offline.exe", "cmd/api/main.go")
	cmd.Env = append(os.Environ(), "GOOS=windows", "GOARCH=amd64", "CGO_ENABLED=0")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error building offline binary: %v\nOutput: %s", err, string(output))
		return
	}
	log.Println("Offline binary generated successfully at static/bin/aura-pos-offline.exe")
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
