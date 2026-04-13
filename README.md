# AURA BACKEND

Backend multi-tenant para sistema POS (Point of Sale) construido con Go y Gin.

## Stack Tecnológico

- **Lenguaje**: Go 1.26.1
- **Framework HTTP**: Gin (github.com/gin-gonic/gin)
- **Base de datos**: PostgreSQL con driver `lib/pq`
- **Migraciones**: golang-migrate/v4 con archivos SQL embebidos
- **Autenticación**: JWT (golang-jwt/jwt/v5) con validación de IP
- **Testing**: stretchr/testify + DATA-DOG/go-sqlmock

## Estructura del Proyecto

```
aura-back/
├── cmd/
│   ├── api/main.go              # Punto de entrada y wiring de dependencias
│   └── server/server.go         # Router, middleware y registro de módulos
├── internal/
│   └── db/db.go                 # Pool de conexiones PostgreSQL
├── shared/                      # Código compartido entre módulos
│   ├── domain/vo/               # Value objects (Email, Document)
│   ├── errors/                  # Errores sentinel del dominio
│   ├── events/                  # Interfaces del event bus
│   ├── logging/                 # Handlers de logging genéricos
│   └── response/                # Helpers para respuestas HTTP
├── tenant/                      # Lógica multi-tenant
│   ├── manager.go               # CRUD de tenants y migraciones
│   ├── auth.go                  # JWT, Login y AuthMiddleware
│   ├── middleware.go            # Middleware de resolución de tenant
│   └── migrations/              # Migraciones SQL
│       ├── public/              # Tablas compartidas (enterprises, users, roles)
│       └── tenant/              # Tablas por tenant (products, sales, etc.)
├── modules/                     # Módulos de funcionalidad (autocontenidos)
│   ├── enterprise/              # Gestión de empresas
│   ├── users/                   # Gestión de usuarios y roles
│   ├── catalog/                 # Módulo de catálogo (agrupado)
│   │   ├── products/            # Sub-módulo: productos
│   │   ├── brands/              # Sub-módulo: marcas
│   │   └── categories/          # Sub-módulo: categorías
│   ├── sales/                   # Órdenes de venta
│   ├── cart/                    # Carrito de compras
│   ├── payments/                # Pagos
│   ├── invoices/                # Facturas
│   ├── inventory/               # Inventario
│   ├── third-parties/           # Terceros (clientes/proveedores)
│   ├── reports/                 # Reportes
│   ├── cash/                    # Caja registradora
│   ├── purchases/               # Compras
│   ├── shrinkage/               # Mermas
│   ├── transfers/               # Transferencias
│   ├── commissions/             # Comisiones
│   └── payroll/                 # Nómina
├── infrastructure/
│   └── messaging/memory/        # Implementación del event bus en memoria
├── .env                         # Variables de entorno (no versionado)
├── go.mod
└── go.sum
```

## Inicio Rápido

### Requisitos previos

- Go 1.26.1+
- PostgreSQL 14+

### Configuración

1. Crear archivo `.env` desde el ejemplo:

```bash
cp .env.example .env
```

2. Configurar las variables de entorno:

```env
DATABASE_URL=postgres://user:pass@localhost:5432/aura?sslmode=disable
JWT_SECRET=tu-secreto-jwt-aqui
PORT=8081
```

3. Ejecutar la aplicación:

```bash
go run ./cmd/api/main.go
```

### Comandos Esenciales

```bash
# Ejecutar la aplicación
go run ./cmd/api/main.go

# Compilar
go build ./...

# Ejecutar todos los tests
go test ./...

# Ejecutar un test específico
go test -v -run TestService_Create_ValidSlugFormat ./modules/enterprise/...

# Tests con detección de race conditions
go test -race ./...

# Tests con cobertura
go test -cover ./...
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

# Formatear y verificar código
go fmt ./...
go vet ./...
go mod tidy
```

### Docker

```bash
# Build the image
docker build -t aura-back .

# Run the container
docker run -p 8081:8081 --env-file .env aura-back
```

### Docker Compose

```bash
docker-compose up --build
```

## Arquitectura Multi-Tenant

El sistema utiliza el patrón **schema-per-tenant** de PostgreSQL:

- **Esquema `public`**: Tablas compartidas entre todos los tenants
  - `enterprises` - Registro de empresas
  - `users` - Usuarios del sistema
  - `roles` - Roles y permisos
  - `user_roles` - Relación usuario-rol
  - `plans` - Planes de suscripción

- **Esquema por tenant** (ej: `empresa_uno`): Tablas específicas del tenant
  - `third_parties` - Clientes y proveedores
  - `products` - Catálogo de productos
  - `sales_orders` - Órdenes de venta
  - Y más...

### Flujo de identificación de tenant

1. El middleware de autenticación extrae el `slug` del JWT
2. El middleware de tenant valida que el slug existe y está activo
3. Las operaciones de base de datos se realizan en el esquema correspondiente

## Módulos

Cada módulo sigue esta estructura consistente:

| Archivo         | Responsabilidad                                                                 |
| --------------- | ------------------------------------------------------------------------------- |
| `domain.go`     | Entidades, interfaces de Repository y Service, eventos                          |
| `service.go`    | Lógica de negocio (struct no exportado, constructor retorna interface)          |
| `repository.go` | Implementación PostgreSQL con interface `querier` para soporte de transacciones |
| `handler.go`    | HTTP handlers de Gin con tipos de request/response                              |
| `routes.go`     | Función `Register(public, protected, handler)`                                  |

### Módulos agrupados (Group Modules)

Algunos módulos contienen entidades relacionadas organizadas como **sub-módulos independientes** dentro de una carpeta padre. Cada sub-módulo mantiene su propio `domain.go`, `service.go`, `repository.go`, `handler.go` y `routes.go`, y se registra de forma independiente en `cmd/server/server.go`.

**Ejemplo — módulo `catalog`:**
```
modules/catalog/
├── products/   → package products
├── brands/     → package brands
└── categories/ → package categories
```

En `cmd/server/server.go` cada sub-módulo se registra de forma independiente:
```go
categories.Register(public, protected, categoryH)
brands.Register(public, protected, brandH)
catalogproducts.Register(public, protected, productH)
```

### Agregar un nuevo módulo simple

1. Crear directorio `modules/<nombre>/` con los archivos base
2. Definir entidad, interfaces de repositorio y servicio en `domain.go`
3. Implementar servicio y repositorio
4. Crear handler con validación de request (`ShouldBindJSON`)
5. Definir rutas en `routes.go`
6. Crear migraciones SQL en `tenant/migrations/tenant/`
7. Registrar handler en `cmd/api/main.go` y `cmd/server/server.go`

### Agregar un sub-módulo dentro de un módulo agrupado

1. Crear el sub-directorio `modules/<grupo>/<nombre>/`
2. Seguir el mismo patrón de archivos (`domain.go`, `service.go`, etc.) con `package <nombre>`
3. Registrar el handler y el servicio **independientemente** en `cmd/api/main.go`
4. Añadir la llamada `<nombre>.Register(...)` en `RegisterModules` dentro de `cmd/server/server.go`

## Migraciones

### Crear una nueva migración

```bash
# Para tablas públicas
migrate create -ext sql -dir tenant/migrations/public -seq nombre_migracion

# Para tablas de tenant
migrate create -ext sql -dir tenant/migrations/tenant -seq nombre_migracion
```

Esto genera dos archivos:

```
000009_nombre_migracion.up.sql    ← Aplica los cambios
000009_nombre_migracion.down.sql  ← Revierte los cambios
```

Las migraciones se aplican automáticamente al iniciar el servidor.

## Autenticación

- **Login**: `POST /login` con email y password
- **JWT**: Incluye `user_id`, `enterprise_id`, `tenant_id`, `slug`, `email`, `roles`, `role_level`, `ip`
- **Validación de IP**: El token incluye la IP del cliente para prevenir robo de tokens
- **Passwords**: Hasheados con bcrypt

### Roles y permisos

| Nivel | Rol        | Descripción                        |
| ----- | ---------- | ---------------------------------- |
| 0     | SUPERADMIN | Acceso total al sistema            |
| 1     | ADMIN      | Administrador de empresa           |
| 2     | SUPERVISOR | Supervisor con permisos extendidos |
| 3+    | USER       | Usuario estándar                   |

## Testing

Los tests utilizan:

- **testify/assert**: Para aserciones legibles
- **testify/mock**: Para mocks de interfaces
- **go-sqlmock**: Para simular operaciones de base de datos

### Convención de nombres

```go
// Test<Componente>_<Método>_<Escenario>
func TestService_Create_DuplicateSlug(t *testing.T) { ... }
func TestListRolesByMinLevel_SuperAdmin(t *testing.T) { ... }
```

### Ejemplo de mock

```go
type MockRepository struct { mock.Mock }

func (m *MockRepository) GetBySlug(ctx context.Context, slug string) (*Enterprise, error) {
    args := m.Called(ctx, slug)
    if args.Get(0) == nil { return nil, args.Error(1) }
    return args.Get(0).(*Enterprise), args.Error(1)
}
```

## Dependencias Principales

| Paquete               | Uso                             |
| --------------------- | ------------------------------- |
| `gin-gonic/gin`       | Framework HTTP                  |
| `lib/pq`              | Driver PostgreSQL               |
| `golang-migrate/v4`   | Migraciones de base de datos    |
| `golang-jwt/jwt/v5`   | Autenticación JWT               |
| `stretchr/testify`    | Assertions y mocks para testing |
| `DATA-DOG/go-sqlmock` | Mock de SQL para tests          |
| `joho/godotenv`       | Carga de variables de entorno   |
