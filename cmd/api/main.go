package main

import (
	"context"
	"log"
	"os"

	"github.com/cloud-tech-develop/aura-back/cmd/server"
	"github.com/cloud-tech-develop/aura-back/infrastructure/messaging/memory"
	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/modules/admin/enterprise"
	"github.com/cloud-tech-develop/aura-back/modules/admin/third-parties"
	"github.com/cloud-tech-develop/aura-back/modules/admin/users"
	"github.com/cloud-tech-develop/aura-back/modules/catalog/brands"
	"github.com/cloud-tech-develop/aura-back/modules/catalog/categories"
	catalogProducts "github.com/cloud-tech-develop/aura-back/modules/catalog/products"
	"github.com/cloud-tech-develop/aura-back/modules/catalog/units"
	"github.com/cloud-tech-develop/aura-back/tenant"
	"github.com/joho/godotenv"
	"os/exec"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	driver := os.Getenv("DATABASE_DRIVER")
	if driver == "" {
		driver = "postgres" // Default
	}
	if driver == "sqlite" && dsn == "" {
		dsn = "aura_pos.db" // Default local db
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// ── Database ────────────────────────────────────────────────────────────
	database, err := db.New(driver, dsn)
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

	// ── Migrate existing tenants in background ─────────────────────────────
	go func() {
		log.Println("Migrating existing tenants...")
		if err := tenantMgr.MigrateAll(context.Background()); err != nil {
			log.Printf("MigrateAll: %v", err)
			return
		}
		log.Println("All tenants migrated successfully")
	}()

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

	// Catalog module
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

	productSvc := catalogProducts.NewService(database, eventBus)
	productHandler := catalogProducts.NewHandler(productSvc)

	// Third Parties module
	thirdPartiesHandler := thirdparties.NewHandler(database)

	// ── HTTP Server ──────────────────────────────────────────────────────────
	srv := server.NewServer(database.DB, tenantMgr)
	srv.RegisterModules(enterpriseHandler, usersHandler, categoryHandler, brandHandler, productHandler, thirdPartiesHandler, unitHandler)

	log.Println("servidor en :" + port)

	// Generate offline binary in background
	go buildOfflineBinary()

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
