# Sincronización de Datos Online ↔ Offline

El sistema de sincronización de Aura POS funciona de forma **bidireccional** en tiempo real, permitiendo que las versiones online y offline mantengan sus datos siempre actualizados sin necesidad de intervención manual.

> **IMPORTANTE**: El sistema garantiza el **aislamiento total entre empresas** usando el `tenant_slug`. Los datos de una empresa nunca se mezclan con otra.

---

## Arquitectura General

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           AURA POS BACKEND                                  │
│                                                                             │
│  ┌─────────────┐         ┌─────────────┐         ┌─────────────┐            │
│  │   ONLINE    │◄──────► │  RABBITMQ   │◄──────► │   OFFLINE   │            │
│  │ (PostgreSQL)│  HTTP   │   (Broker)  │  Event  │   (SQLite)  │            │
│  └─────────────┘  POST   └─────────────┘  Topics └─────────────┘            │
│       │                                    │             │                  │
│       │                                    │             │                  │
│       │     POST /offline/sync/products    │             │                  │
│       └───────────────────────────────────►│             │                  │
│                                            │             │                  │
│                                            │   product.created.{tenant}     │
│                                            │   product.updated.{tenant}     │
│                                            │   product.deleted.{tenant}     │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## 1. Aislamiento por Tenant/Slug

El sistema garantiza que los datos de cada empresa nunca se mezclen. Esto se logra en múltiples capas:

### 1.1 Aislamiento en RabbitMQ (Routing Key)

Cada evento incluye el tenant_slug en el routing key:

```
Exchange: aura.prod

├── COLA: empresa_uno
│   └── recibe: product.*.empresa_uno
│
└── COLA: empresa_dos
    └── recibe: product.*.empresa_dos
```

| Evento                  | Routing Key                   | Qué recibe              |
| ----------------------- | ----------------------------- | ----------------------- |
| Empresa 1 crea producto | `product.created.empresa_uno` | Solo app de empresa_uno |
| Empresa 2 crea producto | `product.created.empresa_dos` | Solo app de empresa_dos |

### 1.2 Aislamiento en HTTP API

El endpoint recibe el `tenant_slug` en el request:

```json
POST /offline/sync/products
{
  "tenant_slug": "empresa_uno",
  "products": [...]
}
```

El handler obtiene el slug del JWT token - **no se puede acceder a datos de otra empresa**.

### 1.3 Aislamiento en Base de Datos

**PostgreSQL (Online):**

```sql
-- Schema público
public.enterprises
├── id: 1, slug: empresa_uno
└── id: 2, slug: empresa_dos

-- Schema por tenant
empresa_uno.products  -- Solo productos de empresa_uno
empresa_dos.products  -- Solo productos de empresa_dos
```

**SQLite (Offline):**

- Tablas con `enterprise_id` para filtrar datos por empresa

### 1.4 Código que garantiza el aislamiento

```go
// Repository siempre recibe tenantSlug como parámetro
func (r *repository) GetByID(ctx context.Context, tenantSlug string, id int64) (*Product, error) {
    // Query con SCHEMA específico del tenant
    query := fmt.Sprintf("SELECT * FROM %s.products WHERE id = $1", tenantSlug)
}
```

---

## 2. Online → Offline (RabbitMQ)

### Flujo

```
┌─────────────┐                    ┌─────────────┐                ┌─────────────┐
│   ONLINE    │                    │  RABBITMQ   │                │   OFFLINE   │
│ (PostgreSQL)│                    │   (Broker)  │                │  (SQLite)   │
└──────┬──────┘                    └──────┬──────┘                └──────┬──────┘
       │                                  │                              │
       │ 1. CREATE/UPDATE/DELETE          │                              │
       │    producto                      │                              │
       ├─────────────────────────────────►│                              │
       │                                  │                              │
       │                          2. Publica evento                      │
       │                          "product.created.empresa_uno"          │
       │                                    │                            │
       │                                    ├───────────────────────────►│
       │                                    │                          │
       │                                    │       3. Suscriptor      │
       │                                    │       recibe evento      │
       │                                    │                          │
       │                                    │       4. Actualiza       │
       │                                    │       SQLite local       │
```

### Cómo funciona

1. **Evento trigger**: Cuando un usuario crea, actualiza o elimina un producto en la versión online, el `product_service.go` ejecuta la operación en PostgreSQL.

2. **Publicación de evento**: Inmediatamente después de la operación, el servicio publica un evento a RabbitMQ:
   - `product.created.{tenant_slug}` - Producto creado
   - `product.updated.{tenant_slug}` - Producto actualizado
   - `product.deleted.{tenant_slug}` - Producto eliminado

3. **Distribución**: RabbitMQ distribuye el evento a todas las colas suscritas (apps offline conectadas).

4. **Consumo offline**: La app offline (suscrita a la cola de su tenant) recibe el evento.

5. **Actualización local**: El `SyncHandler` procesa el evento y actualiza la base de datos SQLite local.

### Código relevante

```go
// En modules/catalog/products/service.go

// Create publica ambos eventos
func (s *service) Create(ctx context.Context, tenantSlug string, p *Product) error {
    // ... lógica de creación ...

    // Publica evento de dominio
    s.publish(NewCreatedEvent(p), tenantSlug)

    // Publica evento de sincronización para offline
    s.publishSync(NewSyncCreatedEvent(tenantSlug, p))

    return nil
}
```

---

## 2. Offline → Online (HTTP POST)

### Flujo

```
┌─────────────┐                    ┌─────────────┐
│   OFFLINE   │     HTTP POST      │   ONLINE    │
│  (SQLite)   │ ──────────────────►│ (PostgreSQL)│
└─────────────┘   /offline/sync    └─────────────┘
       │            /products
       │
       │ 1. App offline detecta
       │    cambios locales
       │
       │ 2. Envía productos modificados
       │    al endpoint de sincronización
       │
       │         3. SyncHandler procesa
       │            cada producto
       │
       │         4. Resuelve conflictos
       │            por timestamp
       │
       │         5. Actualiza PostgreSQL
```

### Cómo funciona

1. **Detección de cambios**: La app offline detecta que tiene productos pendientes de sincronizar (creados, editados o eliminados mientras estaba desconectada).

2. **Envío de datos**: Cuando tiene conectividad, envía un `POST /offline/sync/products` con la lista de productos modificados.

3. **Procesamiento**: El `SyncHandler` procesa cada producto:
   - **Create**: Si no existe en online, lo crea; si existe, compara timestamps
   - **Update**: Actualiza si la versión offline es más reciente
   - **Delete**: Elimina si la versión offline es más reciente

4. **Respuesta**: Devuelve el resultado de cada operación (éxito o error).

### Endpoint

```
POST /offline/sync/products
```

---

## 3. Resolución de Conflictos

### Estrategia: Last Write Wins (Última escritura gana)

```
┌─────────────────┐      ┌─────────────────┐
│     ONLINE      │      │    OFFLINE      │
│  updated_at:    │      │  updated_at:    │
│  2026-05-04     │      │ 2026-05-03      │
│     10:00:00    │      │    15:00:00     │
└────────┬────────┘      └────────┬────────┘
         │                        │
         │    ┌───────────────────┘
         │    │ Comparación de timestamps
         │    ▼
    ┌────┴────┐
    │ WINNER  │
    └────┬────┘
         │
         ▼
   GANA ONLINE (timestamp más reciente)
```

### Reglas

| Escenario                                  | Resultado                                                   |
| ------------------------------------------ | ----------------------------------------------------------- |
| `updated_at` online > `updated_at` offline | **Gana online** (se mantiene la versión online)             |
| `updated_at` offline > `updated_at` online | **Gana offline** (se actualiza online con datos de offline) |
| Delete con conflicto                       | No se elimina si online tiene versión más reciente          |

### Código relevante

```go
// En modules/catalog/products/sync_handler.go

func isOnlineNewer(onlineTimestamp *vo.DateTime, offlineTimestamp time.Time) bool {
    if onlineTimestamp == nil {
        return false
    }
    return convertToTime(onlineTimestamp).After(offlineTimestamp)
}
```

---

## 4. Ejemplos de Request/Response

### Request (Offline → Online)

```json
POST /offline/sync/products
{
  "tenant_slug": "empresa_uno",
  "products": [
    {
      "product_id": 100,
      "sku": "PROD-001",
      "barcode": "123456789",
      "name": "Nuevo Producto",
      "description": "Descripción del producto",
      "cost_price": 100.00,
      "sale_price": 150.00,
      "active": true,
      "action": "create",
      "timestamp": "2026-05-04T10:00:00Z",
      "source": "offline"
    },
    {
      "product_id": 101,
      "sku": "PROD-002",
      "name": "Producto actualizado",
      "action": "update",
      "timestamp": "2026-05-04T11:30:00Z",
      "source": "offline"
    },
    {
      "product_id": 102,
      "action": "delete",
      "timestamp": "2026-05-04T12:00:00Z",
      "source": "offline"
    }
  ]
}
```

### Response

```json
{
  "data": {
    "results": [
      {
        "product_id": 100,
        "sku": "PROD-001",
        "action": "create",
        "status": "success",
        "message": "Product created successfully"
      },
      {
        "product_id": 101,
        "sku": "PROD-002",
        "action": "update",
        "status": "success",
        "message": "Product updated successfully"
      },
      {
        "product_id": 102,
        "action": "delete",
        "status": "success",
        "message": "Product deleted successfully"
      }
    ],
    "total_synced": 3,
    "success_count": 3,
    "error_count": 0,
    "sync_time": "2026-05-04T12:00:00Z"
  },
  "message": "Sincronizacion completada",
  "success": true
}
```

---

## 5. Estructura de Eventos

### Eventos de Sincronización

| Evento                    | Dirección        | Routing Key                      | Descripción                     |
| ------------------------- | ---------------- | -------------------------------- | ------------------------------- |
| `product.created`         | Online → Offline | `product.created.{slug}`         | Producto creado en online       |
| `product.updated`         | Online → Offline | `product.updated.{slug}`         | Producto actualizado en online  |
| `product.deleted`         | Online → Offline | `product.deleted.{slug}`         | Producto eliminado en online    |
| `product.offline.created` | Offline → Online | `product.offline.created.{slug}` | Producto creado en offline      |
| `product.offline.updated` | Offline → Online | `product.offline.updated.{slug}` | Producto actualizado en offline |
| `product.offline.deleted` | Offline → Online | `product.offline.deleted.{slug}` | Producto eliminado en offline   |

### Formato de Routing Key en RabbitMQ

```
aura.{env}.{event}.{tenant}

Ejemplo:
aura.prod.product.created.empresa_uno
aura.dev.product.updated.miempresa
```

---

## 6. Uso en el Código

### Configuración del Handler

```go
// Con sincronización habilitada
syncHandler := products.NewSyncHandler(products.NewRepository(db), eventBus)
handler := products.NewHandlerWithSync(productSvc, syncHandler)

// Sin sincronización (retrocompatible)
handler := products.NewHandler(productSvc)
```

### Suscripción en App Offline

```go
// La app offline se suscribe a los eventos de su tenant
bus.Subscribe("product.created."+tenantSlug, syncHandler.SyncHandlerFunc())
bus.Subscribe("product.updated."+tenantSlug, syncHandler.SyncHandlerFunc())
bus.Subscribe("product.deleted."+tenantSlug, syncHandler.SyncHandlerFunc())
```

---

## 7. Inicio de Sincronización en App Offline

La aplicación offline puede iniciar la sincronización de dos formas:

### 7.1 Sincronización Automática (al iniciar la app)

Cuando la app offline se inicia, puede detectar automáticamente si hay una enterprise configurada en SQLite y usar su slug para sincronizar:

```
┌─────────────────────────────────────────────────────────────────┐
│                    AURA POS OFFLINE                             │
│                                                                  │
│   1. App inicia                                                 │
│      │                                                           │
│      ▼                                                           │
│   2. Verificar tabla 'enterprises' en SQLite                   │
│      │                                                           │
│      │   SELECT * FROM enterprises LIMIT 1                      │
│      │                                                            │
│      ▼                                                           │
│   3. ¿Existe enterprise?                                         │
│      │                                                           │
│      ├────── SÍ ──────► Usar slug para sync automática         │
│      │                POST /offline/ping                         │
│      │                                                            │
│      └────── NO ──────► Esperar sync manual del usuario         │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 7.2 Sincronización Manual (endpoint /offline/ping)

El endpoint `GET /offline/ping` permite iniciar la sincronización:

#### Flujo del Endpoint

```
GET /offline/ping
Authorization: Bearer <token>

1. Verificar modo offline (SQLite)
   │
   ├──► SI: Continuar
   │
   └──► NO: Retornar error 403

2. Obtener slug del token JWT
   │
   ├──► SI tiene slug ──► Usar slug del token
   │
   └──► NO tiene slug ──► Buscar en tabla enterprises
       │
       ├──► SI existe ──► Usar slug de la enterprise
       │
       └──► NO existe ──► Retornar error: "No hay enterprise configurada"

3. Obtener URL de producción (URL_PROD)
   │
   └──► Default: http://localhost:8081

4. Ejecutar sincronización completa
   │
   └──► SyncAllBySlug(prodURL, token, slug)

5. Retornar resultado
```

#### Response Exitoso

```json
{
  "data": {
    "slug": "empresa_uno",
    "source": "http://localhost:8081",
    "mode": "offline",
    "sync_source": "local_db",
    "result": {
      "enterprises": 1,
      "plans": 1,
      "users": 5,
      "categories": 10,
      "brands": 8,
      "units": 5,
      "products": 150,
      "presentations": 200
    },
    "message": "Sincronización completada"
  },
  "success": true
}
```

#### Response con error (sin enterprise)

```json
{
  "data": "No hay enterprise configurada. Primero sincronice con /offline/sync-tenant",
  "success": false
}
```

### 7.3 Código del Handler

```go
// En modules/offline/handler.go

func (h *Handler) Ping(c *gin.Context) {
    // 1. Verificar modo offline
    isOffline := driver == "sqlite" || dsn == ""
    if !isOffline {
        response.Forbidden(c, "Endpoint solo disponible en modo offline")
        return
    }

    // 2. Obtener slug del token JWT
    slug, ok := tenant.SlugFromContext(c)
    if ok && slug != "" {
        // Usar slug del token
        syncSource = "token"
    } else {
        // 3. Fallback: obtener de SQLite
        enterprise, err := h.svc.GetActiveEnterprise(ctx)
        if err != nil {
            response.BadRequest(c, "No hay enterprise configurada")
            return
        }
        slug = enterprise.Slug
        syncSource = "local_db"
    }

    // 4. Ejecutar sincronización
    result, err := h.svc.SyncAllBySlug(ctx, prodURL, token, slug)

    // 5. Retornar resultado
    response.OK(c, gin.H{
        "slug":        slug,
        "sync_source": syncSource,
        "result":      result,
    })
}
```

---

## Resumen

### Sincronización Automática

Cuando creas un producto en cualquier ambiente, **se sincroniza automáticamente**:

| Acción                    | Qué sucede                                                          |
| ------------------------- | ------------------------------------------------------------------- |
| Crear producto en Online  | → Publica evento a RabbitMQ → App offline recibe y actualiza SQLite |
| Crear producto en Offline | → Envía a `/offline/sync/products` → PostgreSQL se actualiza        |
| Actualizar producto       | → Sincroniza cambios en tiempo real                                 |
| Eliminar producto         | → Sincroniza eliminación en ambos ambientes                         |

### Aislamiento entre Empresas

**Los datos nunca se mezclan entre empresas** gracias a:

| Capa              | Mecanismo                                               |
| ----------------- | ------------------------------------------------------- |
| **RabbitMQ**      | Routing key con slug (`product.created.empresa_uno`)    |
| **HTTP API**      | JWT con slug en claims                                  |
| **Base de datos** | Schema per tenant (PostgreSQL) + enterprise_id (SQLite) |

---

## Resumen Técnico

| Dirección            | Mecanismo                | Cuándo                  | Contenido                    |
| -------------------- | ------------------------ | ----------------------- | ---------------------------- |
| **Online → Offline** | RabbitMQ (eventos)       | Tiempo real             | Productos modificados        |
| **Offline → Online** | HTTP POST                | Cuando hay conectividad | Productos pendientes         |
| **Conflictos**       | Timestamp (`updated_at`) | Automático              | Última versión gana          |
| **Aislamiento**      | tenant_slug en todo      | Siempre                 | Datos protegidos por empresa |

---

## Ventajas del Sistema

1. **Tiempo real**: Los cambios en online se reflejan inmediatamente en offline
2. **Bidireccional**: Ambos sistemas pueden iniciar cambios
3. **Conflictos automáticos**: No requiere intervención manual
4. **Aislamiento garantizado**: Cada empresa tiene sus datos aislados
5. **Escalable**: Funciona con múltiples tenants
6. **Tolerante a fallos**: Funciona cuando hay conectividad intermitente
