package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"os/exec"

	"github.com/cloud-tech-develop/aura-back/cmd/server"
	"github.com/cloud-tech-develop/aura-back/infrastructure/messaging/memory"
	"github.com/cloud-tech-develop/aura-back/infrastructure/messaging/rabbit"
	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/modules/admin/enterprise"
	thirdparties "github.com/cloud-tech-develop/aura-back/modules/admin/third-parties"
	"github.com/cloud-tech-develop/aura-back/modules/admin/users"
	"github.com/cloud-tech-develop/aura-back/modules/catalog/brands"
	"github.com/cloud-tech-develop/aura-back/modules/catalog/categories"
	catalogPresentations "github.com/cloud-tech-develop/aura-back/modules/catalog/presentations"
	"github.com/cloud-tech-develop/aura-back/modules/catalog/presentations"
	catalogProducts "github.com/cloud-tech-develop/aura-back/modules/catalog/products"
	"github.com/cloud-tech-develop/aura-back/modules/catalog/units"
	"github.com/cloud-tech-develop/aura-back/modules/offline"
	"github.com/cloud-tech-develop/aura-back/shared/events"
	"github.com/cloud-tech-develop/aura-back/tenant"
	"github.com/joho/godotenv"
)

var (
	DefaultUseRabbitMQ = ""
	DefaultRabbitHost  = ""
	DefaultRabbitUser  = ""
	DefaultRabbitPass  = ""
	DefaultAppEnv      = ""
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	// Apply baked-in defaults for offline binaries if env vars are missing
	if os.Getenv("USE_RABBITMQ") == "" && DefaultUseRabbitMQ != "" {
		os.Setenv("USE_RABBITMQ", DefaultUseRabbitMQ)
	}
	if os.Getenv("RABBITMQ_HOST") == "" && DefaultRabbitHost != "" {
		os.Setenv("RABBITMQ_HOST", DefaultRabbitHost)
	}
	if os.Getenv("RABBITMQ_DEFAULT_USER") == "" && DefaultRabbitUser != "" {
		os.Setenv("RABBITMQ_DEFAULT_USER", DefaultRabbitUser)
	}
	if os.Getenv("RABBITMQ_DEFAULT_PASS") == "" && DefaultRabbitPass != "" {
		os.Setenv("RABBITMQ_DEFAULT_PASS", DefaultRabbitPass)
	}
	if os.Getenv("APP_ENV") == "" && DefaultAppEnv != "" {
		os.Setenv("APP_ENV", DefaultAppEnv)
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
	var eventBus events.EventBus

	if os.Getenv("USE_RABBITMQ") == "true" && driver != "sqlite" {
		// Production mode: use RabbitMQ immediately
		log.Println("[Main] Using RabbitMQ event bus")
		rb, err := rabbit.NewRabbitMQEventBus()
		if err != nil {
			log.Printf("[Main] RabbitMQ not available, falling back to memory: %v", err)
			eventBus = memory.NewMemoryEventBus(100, 5)
		} else {
			eventBus = rb
		}
	} else if os.Getenv("USE_RABBITMQ") == "true" && driver == "sqlite" {
		// Offline mode: we'll use RabbitMQ AFTER we get the slug from /offline/ping
		// Use memory event bus initially, will upgrade to RabbitMQ when slug is known
		log.Println("[Main] Offline mode: will use RabbitMQ after tenant sync")
		eventBus = memory.NewMemoryEventBus(100, 5)
	} else {
		log.Println("[Main] Using memory event bus")
		eventBus = memory.NewMemoryEventBus(100, 5)
	}

	if err := eventBus.Start(); err != nil {
		log.Fatalf("Failed to start event bus: %v", err)
	}
	defer eventBus.Stop()

	// Tenant Manager & Migrations
	tenantMgr := tenant.NewManager(database.DB)
	if driver == "sqlite" {
		// Offline mode: runs both public and tenant offline migrations from offline/ folder
		if err := tenantMgr.MigrateOffline(); err != nil {
			log.Fatal("MigrateOffline:", err)
		}
	} else {
		// Online mode: run only public (Postgres) migrations
		if err := tenantMgr.MigratePublic(); err != nil {
			log.Fatal("MigratePublic:", err)
		}
		// Migrate existing tenant schemas in background
		go func() {
			log.Println("Migrating existing tenants...")
			if err := tenantMgr.MigrateAll(context.Background()); err != nil {
				log.Printf("MigrateAll: %v\n", err)
				return
			}
			log.Println("All tenants migrated successfully")
		}()
	}

	// Modules
	// Enterprise module
	enterpriseMigrator := &enterpriseMigratorAdapter{manager: tenantMgr}
	enterpriseSvc := enterprise.NewService(database, eventBus, enterpriseMigrator)
	enterpriseHandler := enterprise.NewHandler(enterpriseSvc, tenantMgr)

	// Logging (Enterprise)
	enterpriseLogger := enterprise.NewLoggerHandler("logs")
	_ = eventBus.Subscribe(enterprise.EventCreated, enterpriseLogger)
	_ = eventBus.Subscribe(enterprise.EventUpdated, enterpriseLogger)
	_ = eventBus.Subscribe(enterprise.EventDeleted, enterpriseLogger)

	// Users module
	usersSvc := users.NewService(database, eventBus)
	usersHandler := users.NewHandler(usersSvc)

	// Logging (Users)
	usersLogger := users.NewLoggerHandler("logs")
	_ = eventBus.Subscribe(users.EventCreated, usersLogger)
	_ = eventBus.Subscribe(users.EventUpdated, usersLogger)

	// Catalog modules
	categorySvc := categories.NewService(database, eventBus)
	categoryHandler := categories.NewHandler(categorySvc)

	brandSvc := brands.NewService(database, eventBus)
	brandHandler := brands.NewHandler(brandSvc)

	unitSvc := units.NewService(database)
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
		offlineSvc := offline.NewService(
			database,
			eventBus,
			tenantMgr,
			productSvc,
			presSvc,
			categorySvc,
			brandSvc,
			unitSvc,
			catalogPresentations.NewRepository(database),
		)
		offlineHandler = offline.NewHandler(offlineSvc)
	}

	// HTTP Server
	srv := server.NewServer(database, tenantMgr)
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

	useRabbit := os.Getenv("USE_RABBITMQ")
	host := os.Getenv("RABBITMQ_HOST")
	user := os.Getenv("RABBITMQ_DEFAULT_USER")
	pass := os.Getenv("RABBITMQ_DEFAULT_PASS")
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "dev"
	}

	log.Printf("Building offline binary (USE_RABBITMQ=%s, APP_ENV=%s)", useRabbit, appEnv)

	ldflags := fmt.Sprintf(
		`-X main.DefaultUseRabbitMQ=%s -X main.DefaultRabbitHost=%s -X main.DefaultRabbitUser=%s -X main.DefaultRabbitPass=%s -X main.DefaultAppEnv=%s`,
		useRabbit, host, user, pass, appEnv,
	)

	cmd := exec.Command("go", "build", "-ldflags", ldflags, "-o", "static/bin/aura-pos-offline.exe", "cmd/api/main.go")
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
