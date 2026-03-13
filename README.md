# AURA BACKEND


---

## Flujo completo

```
Arranque del servidor
    ↓
MigratePublic()  →  crea public.tenants
    ↓
MigrateAll()     →  migra esquemas de todos los tenants existentes
    ↓
POST /tenants    →  Create() → nuevo esquema + migraciones automáticas
    ↓
GET /me con X-Tenant: empresa_uno
    →  Middleware valida → handler opera en esquema empresa_uno
```

## Agregar una nueva migración

# Estando en la raíz de tu proyecto
```
migrate create -ext sql -dir migrations -seq nombre_de_la_migracion
```

Esto genera automáticamente dos archivos:
```
migrations/
  000001_nombre_de_la_migracion.up.sql    ← lo que se aplica
  000001_nombre_de_la_migracion.down.sql  ← cómo revertirlo
```

Solo creas los archivos y al reiniciar el servidor se aplican solas a todos los tenants:

```
migrations/
  000003_orders.up.sql    ← nueva tabla
  000003_orders.down.sql  ← rollback
```


### Estructura
```
aura-back/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── tenant/
│   │   ├── manager.go       # Crear/gestionar tenants
│   │   └── middleware.go    # Resolver tenant por request
│   ├── migration/
│   │   ├── migrator.go      # Lógica de migraciones
│   │   └── migrations/      # Archivos SQL
│   │       ├── 001_init.sql
│   │       └── 002_users.sql
│   └── db/
│       └── db.go            # Pool de conexiones
├── go.mod
└── go.sum
```


## 1. Crear una migración
> Estando en la raíz de tu proyecto
```
migrate create -ext sql -dir migrations -seq nombre_de_la_migracion
```

Esto genera automáticamente dos archivos:
```
migrations/
  000001_nombre_de_la_migracion.up.sql    ← lo que se aplica
  000001_nombre_de_la_migracion.down.sql  ← cómo revertirlo
```


## 2. Ejemplos prácticos 
> Estando en la raíz de tu proyecto
```
migrate create -ext sql -dir migrations -seq nombre_de_la_migracion
```

Esto genera automáticamente dos archivos:
```
migrations/
  000001_nombre_de_la_migracion.up.sql    ← lo que se aplica
  000001_nombre_de_la_migracion.down.sql  ← cómo revertirlo
```
