# HU-011 - Asignar roles a usuario

## 📌 Información General
- ID: HU-011
- Epic: EPIC-002 - Gestión de Usuarios de Empresa
- Prioridad: Alta
- Estado: Backlog
- Porcentaje: 0%
- Autor: QA & Requirements Engineer
- Fecha: 2026-03-15

---

## 👤 Historia de Usuario

**Como** administrador de empresa
**Quiero** asignar roles a un usuario
**Para** controlar qué permisos tiene cada miembro del equipo

---

## 🧠 Descripción Funcional

El sistema debe permitir asignar roles específicos a un usuario. Los roles se definen en public.roles y se asignan via public.user_roles.

---

## ✅ Criterios de Aceptación

### Escenario 1: Asignación exitosa
- Dado que soy administrador autenticado
- Cuando envío PATCH /users/:id/roles con IDs de roles válidos
- Entonces el sistema asigna los roles al usuario
- Y elimina roles anteriores (reemplazo total)

### Escenario 2: Roles inexistentes
- Dado que envío IDs de roles que no existen
- Cuando envío PATCH /users/:id/roles
- Entonces el sistema devuelve 400 Bad Request

### Escenario 3: Usuario no encontrado
- Dado que solicito asignar roles a usuario inexistente
- Cuando envío PATCH /users/99999/roles
- Entonces el sistema devuelve 404 Not Found

### Escenario 4: Multi-tenant validation
- Dado que soy administrador de "empresa_uno"
- Cuando asigno roles a usuario de "empresa_dos"
- Entonces el sistema devuelve 404 Not Found

### Escenario 5: Asignación de roles vacía
- Dado que envío lista vacía de roles
- Cuando envío PATCH /users/:id/roles con []
- Entonces el sistema elimina todos los roles del usuario

---

## ❌ Casos de Error

- Si roles no existen → Error 400 Bad Request
- Si usuario no existe → Error 404 Not Found
- Si usuario no pertenece a empresa → Error 404 Not Found
- Si IDs inválidos → Error 400 Bad Request

---

## 🔐 Reglas de Negocio

- Solo usuarios de la empresa autenticada pueden ser modificados
- Los roles se reemplazan completamente (no se añaden)
- Los IDs de roles deben existir en public.roles
- No se permite asignar roles de nivel superior (ej: SUPERADMIN)
- El usuario debe tener al menos un rol

---

## 🎨 Consideraciones UI/UX

- Selector de roles multi-select
- Mostrar roles disponibles
- Confirmación de cambio de roles
- Indicador de roles actuales

---

## 📡 Requisitos Técnicos

- Endpoint: PATCH /users/:id/roles
- Método HTTP: PATCH
- Path parameter: id (int64)
- Request body:
```json
{
  "role_ids": [2, 3, 4]
}
```
- Response 200:
```json
{
  "data": {
    "user_id": 1,
    "roles": [
      {
        "id": 2,
        "name": "ADMIN",
        "description": "Administrator"
      },
      {
        "id": 3,
        "name": "SUPERVISOR",
        "description": "Supervisor"
      }
    ]
  },
  "success": true,
  "message": "Roles actualizados exitosamente"
}
```
- Códigos de error:
  - 400: Roles inválidos o faltantes
  - 404: Usuario no encontrado
  - 500: Error interno

---

## 🧪 Criterios de Testing

- Unit tests: Service layer (assign roles, validate role IDs)
- Integration tests: Handler layer (HTTP request/response, error cases)
- E2E: Asignación de roles y verificación en JWT

---

## 📎 Dependencias

- Servicio: UserService
- Repositorio: UserRepository
- Otra HU relacionada: HU-006 (Crear usuario), HU-001 (Login)

---

## 🚫 Fuera de Alcance

- No incluye creación de nuevos roles
- No incluye permisos granular por módulo
- No incluye herencia de roles

---

## 🧠 Generación de Código

Requerir:
- UserService.AssignRoles
- UserRepository.AssignRoles
- Handler para PATCH /users/:id/roles
- Validación de IDs de roles existentes
