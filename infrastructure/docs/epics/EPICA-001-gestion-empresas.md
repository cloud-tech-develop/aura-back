# EPICA-001: Gestión de Empresas (Tenants) - Aura POS

## Descripción
Esta épica abarca todo el flujo de gestión de empresas (tenants) en el sistema Aura POS, desde la creación inicial hasta la gestión del ciclo de vida de cada empresa. El sistema utiliza un patrón multi-tenant con esquema por tenant en PostgreSQL.

---

## Estado Actual del Sistema

### Flujo de Creación Existente
```
POST /enterprises (público)
    ├── Handler: modules/enterprise/handler.go
    ├── Service: modules/enterprise/service.go (validaciones)
    └── Manager: tenant/manager.go (transacción DB)
        ├── 1. Crear/Actualizar tenant (public.tenants)
        ├── 2. Crear schema PostgreSQL
        ├── 3. Registrar enterprise (public.enterprises)
        ├── 4. Crear usuario admin (public.users)
        ├── 4.1 Asignar rol ADMIN (user_roles)
        ├── 5. Ejecutar migraciones del tenant
        └── 6. Crear tercero inicial (empleado)
```

### Flujo de Autenticación
```
POST /login
    ├── Valida email y password
    ├── Consulta enterprise del usuario
    └── Genera JWT con: user_id, enterprise_id, tenant_id, slug
```

### Flujo de Resolución de Tenant
```
1. AuthMiddleware valida JWT y extrae slug del token
2. TenantMiddleware usa slug del token (prioridad)
3. Fallback: subdominio si no hay token
```

---

## Historias de Usuario

| HU | Título | Estado | 
|----|--------|--------|
| HU-001 | Crear nueva empresa (Registro público) | ✅ 12/12 |
| HU-002 | Validar unicidad de campos (slug, email, subdominio) | ✅ 12/12 |
| HU-003 | Asignar rol de administrador al usuario inicial | ✅ 3/3 |
| HU-004 | Listar empresas (Admin) | ✅ 4/4 |
| HU-005 | Obtener empresa por slug | ✅ 3/3 |
| HU-006 | Actualizar información de empresa | ✅ 5/5 |
| HU-007 | Cambiar estado de empresa (suspensión) | ✅ 4/4 |
| HU-008 | Validar límite de empresas según plan | ⏳ 0/4 |
| HU-009 | Resolución de tenant por token JWT | ✅ 4/4 |
| HU-010 | Logging de auditoría de creación de empresa | ⏳ 0/3 |

---

## Implementaciones Realizadas

### HU-001: Crear empresa
- ✅ Validación de longitud de slug (3-50 caracteres)
- ✅ Validación de contraseña (mínimo 8 caracteres)
- ✅ Normalización de slug y subdominio a minúsculas
- ✅ Validación de formato de email

### HU-002: Validar unicidad
- ✅ Validación de slug único en enterprises
- ✅ Validación de subdominio único en enterprises
- ✅ Validación de email único en enterprises
- ✅ Validación de email único en users

### HU-003: Asignar rol ADMIN
- ✅ Inserción en user_roles al crear usuario

### HU-004: Listar empresas
- ✅ Paginación (page, limit)
- ✅ Filtros por estado
- ✅ Endpoint GET /enterprises

### HU-005: Obtener empresa
- ✅ Retorno 404 cuando no existe
- ✅ Endpoint GET /enterprises/:slug

### HU-006: Actualizar empresa
- ✅ Endpoint PUT /enterprises/:slug
- ✅ Solo actualizar campos proporcionados
- ✅ No permitir cambio de slug

### HU-007: Cambiar estado
- ✅ Endpoint PATCH /enterprises/:slug/status
- ✅ Validación de estados válidos (ACTIVE, INACTIVE, SUSPENDED, DEBT)

### HU-009: Resolución de Tenant
- ✅ Token JWT incluye tenant_id y slug
- ✅ Middleware usa slug del token (no del header X-Tenant)
- ✅ Fallback a subdominio si no hay token
- ✅ Validación de empresa ACTIVE

---

## Pendiente

### HU-008: Validar plan
- Validación de cuota antes de crear empresa
- Retornar 403 cuando se alcance el límite

### HU-010: Auditoría
- Logging estructurado

---

## Resumen
- **Total de HU**: 10
- **Completadas**: 8 (HU-001 a HU-007, HU-009)
- **Pendientes**: 2 (HU-008, HU-010)
