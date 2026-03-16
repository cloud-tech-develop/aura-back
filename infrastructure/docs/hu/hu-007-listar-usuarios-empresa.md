# HU-007 - Listar usuarios de una empresa

## 📌 Información General
- ID: HU-007
- Epic: EPIC-002 - Gestión de Usuarios de Empresa
- Prioridad: Media
- Estado: Backlog
- Porcentaje: 0%
- Autor: QA & Requirements Engineer
- Fecha: 2026-03-15

---

## 👤 Historia de Usuario

**Como** administrador de empresa
**Quiero** ver la lista de usuarios de mi empresa
**Para** gestionar y monitorear el acceso de mi equipo

---

## 🧠 Descripción Funcional

El sistema debe permitir listar los usuarios de la empresa autenticada con paginación y filtros. La lista solo debe incluir usuarios de la empresa del usuario autenticado (multi-tenant validation).

---

## ✅ Criterios de Aceptación

### Escenario 1: Listado exitoso
- Dado que soy administrador autenticado
- Cuando solicito GET /users
- Entonces el sistema devuelve la lista de usuarios de mi empresa
- Y la respuesta incluye datos de paginación

### Escenario 2: Con filtros
- Dado que existen usuarios activos e inactivos
- Cuando solicito GET /users?status=active
- Entonces el sistema devuelve solo los usuarios activos

### Escenario 3: Con paginación
- Dado que existen más usuarios que el límite por página
- Cuando solicito GET /users?page=2&limit=5
- Entonces el sistema devuelve los usuarios de la página 2 (5 por página)

### Escenario 4: Sin usuarios
- Dado que la empresa no tiene usuarios (solo el administrador inicial)
- Cuando solicito GET /users
- Entonces el sistema devuelve lista vacía o solo el admin

### Escenario 5: Multi-tenant validation
- Dado que soy administrador de "empresa_uno"
- Cuando solicito GET /users
- Entonces solo veo usuarios de "empresa_uno"
- Y no veo usuarios de otras empresas

---

## ❌ Casos de Error

- Si no hay usuarios → Lista vacía (no error)
- Si el usuario no está autenticado → Error 401 Unauthorized
- Si parámetros de paginación inválidos → Usa valores por defecto

---

## 🔐 Reglas de Negocio

- Solo usuarios de la empresa autenticada se muestran
- No se incluyen usuarios eliminados (deleted_at IS NULL)
- Paginación por defecto: page=1, limit=10
- Los filtros son opcionales
- El usuario autenticado siempre se incluye en la lista

---

## 🎨 Consideraciones UI/UX

- Mostrar estado (activo/inactivo) visualmente
- Indicador de carga durante la consulta
- Mensaje amigable si no hay usuarios
- Paginación clara con números de página

---

## 📡 Requisitos Técnicos

- Endpoint: GET /users
- Método HTTP: GET
- Query parameters:
  - page (int, default: 1)
  - limit (int, default: 10, max: 100)
  - status (string, "active" o "inactive")
- Response 200:
```json
{
  "data": {
    "data": [
      {
        "id": 1,
        "enterprise_id": 50,
        "name": "Admin",
        "email": "admin@empresa.com",
        "active": true,
        "created_at": "2026-03-15T10:00:00Z",
        "updated_at": "2026-03-15T10:00:00Z",
        "roles": ["ADMIN"],
        "third_party": {
          "first_name": "Admin",
          "last_name": "Empresa",
          "document_number": "123456"
        }
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "totalPages": 1
    }
  },
  "success": true,
  "message": "Operación exitosa"
}
```

---

## 🧪 Criterios de Testing

- Unit tests: Service layer (list users, pagination, filters)
- Integration tests: Handler layer (HTTP request/response)
- E2E: Listado de usuarios de empresa específica

---

## 📎 Dependencias

- Servicio: UserService
- Repositorio: UserRepository
- Otra HU relacionada: HU-006 (Crear usuario)

---

## 🚫 Fuera de Alcance

- No incluye búsqueda por nombre/email
- No incluye exportación de datos
- No incluye usuarios de otras empresas

---

## 🧠 Generación de Código

Requerir:
- Actualizar UserService.List
- Actualizar UserRepository.List
- Handler para GET /users
