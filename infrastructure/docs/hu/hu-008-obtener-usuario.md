# HU-008 - Obtener detalles de un usuario específico

## 📌 Información General
- ID: HU-008
- Epic: EPIC-002 - Gestión de Usuarios de Empresa
- Prioridad: Media
- Estado: Backlog
- Porcentaje: 0%
- Autor: QA & Requirements Engineer
- Fecha: 2026-03-15

---

## 👤 Historia de Usuario

**Como** administrador de empresa
**Quiero** ver los detalles de un usuario específico
**Para** revisar su información y configuración

---

## 🧠 Descripción Funcional

El sistema debe permitir obtener los detalles completos de un usuario específico por su ID, incluyendo información del tercero asociado y sus roles.

---

## ✅ Criterios de Aceptación

### Escenario 1: Obtención exitosa
- Dado que soy administrador autenticado
- Cuando solicito GET /users/:id
- Entonces el sistema devuelve los detalles del usuario
- Y incluye información del tercero asociado
- Y incluye la lista de roles

### Escenario 2: Usuario no encontrado
- Dado que solicito un usuario con ID inexistente
- Cuando envío GET /users/99999
- Entonces el sistema devuelve 404 Not Found

### Escenario 3: Usuario de otra empresa
- Dado que soy administrador de "empresa_uno"
- Cuando solicito GET /users/:id de "empresa_dos"
- Entonces el sistema devuelve 404 Not Found (o 403 Forbidden)

### Escenario 4: Multi-tenant validation
- Dado que soy administrador de "empresa_uno"
- Cuando solicito GET /users/:id
- Entonces el sistema verifica que el usuario pertenezca a "empresa_uno"
- Y devuelve error si pertenece a otra empresa

---

## ❌ Casos de Error

- Si el usuario no existe → Error 404 Not Found
- Si el usuario no pertenece a la empresa autenticada → Error 404 Not Found
- Si no está autenticado → Error 401 Unauthorized

---

## 🔐 Reglas de Negocio

- Solo se pueden ver usuarios de la empresa autenticada
- El tercero se obtiene del esquema del tenant
- Los roles se obtienen de public.user_roles
- No se muestran usuarios eliminados (deleted_at IS NULL)

---

## 🎨 Consideraciones UI/UX

- Mostrar información completa del usuario
- Sección de roles asignados
- Datos del tercero asociado
- Indicador de carga

---

## 📡 Requisitos Técnicos

- Endpoint: GET /users/:id
- Método HTTP: GET
- Path parameter: id (int64)
- Response 200:
```json
{
  "data": {
    "id": 1,
    "enterprise_id": 50,
    "name": "Admin",
    "email": "admin@empresa.com",
    "active": true,
    "created_at": "2026-03-15T10:00:00Z",
    "updated_at": "2026-03-15T10:00:00Z",
    "roles": [
      {
        "id": 1,
        "name": "ADMIN",
        "description": "Administrator"
      }
    ],
    "third_party": {
      "id": 1,
      "first_name": "Admin",
      "last_name": "Empresa",
      "document_number": "123456",
      "document_type": "CC",
      "personal_email": "admin@empresa.com",
      "tax_responsibility": "RESPONSIBLE",
      "is_employee": true
    }
  },
  "success": true,
  "message": "Operación exitosa"
}
```
- Response 404:
```json
{
  "data": null,
  "success": false,
  "message": "Usuario no encontrado"
}
```

---

## 🧪 Criterios de Testing

- Unit tests: Service layer (get user by ID, validation)
- Integration tests: Handler layer (HTTP request/response, 404 cases)
- E2E: Flujo completo de obtención de usuario

---

## 📎 Dependencias

- Servicio: UserService
- Repositorio: UserRepository
- Otra HU relacionada: HU-007 (Listar usuarios)

---

## 🚫 Fuera de Alcance

- No incluye usuarios de otras empresas
- No incluye información sensible (passwords)
- No incluye historial de cambios

---

## 🧠 Generación de Código

Requerir:
- UserService.GetByID
- UserRepository.GetByID
- Handler para GET /users/:id
- JOINs para roles y tercero
