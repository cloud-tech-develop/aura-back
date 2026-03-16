# HU-006 - Crear usuario adicional para empresa existente

## 📌 Información General
- ID: HU-006
- Epic: EPIC-002 - Gestión de Usuarios de Empresa
- Prioridad: Alta
- Estado: Backlog
- Porcentaje: 0%
- Autor: QA & Requirements Engineer
- Fecha: 2026-03-15

---

## 👤 Historia de Usuario

**Como** administrador de empresa
**Quiero** crear usuarios adicionales para mi empresa
**Para** asignar diferentes roles y permisos a mi equipo de trabajo

---

## 🧠 Descripción Funcional

El sistema debe permitir crear usuarios adicionales para una empresa ya existente. Al crear un usuario, se generará automáticamente un tercero en el esquema del tenant. El usuario tendrá un email único en public.users y estará asociado a la empresa del usuario autenticado.

---

## ✅ Criterios de Aceptación

### Escenario 1: Creación exitosa de usuario adicional
- Dado que soy un administrador autenticado
- Cuando envío una solicitud POST a /users con datos válidos
- Entonces el sistema crea el usuario en public.users con enterprise_id de mi empresa
- Y crea automáticamente un tercero en el esquema del tenant
- Y asigna los roles especificados al usuario
- Y devuelve 201 Created con los datos del usuario creado

### Escenario 2: Email duplicado
- Dado que existe un usuario con email "usuario@empresa.com"
- Cuando intento crear otro usuario con el mismo email
- Entonces el sistema devuelve 409 Conflict con mensaje "El email ya está registrado"

### Escenario 3: Datos inválidos
- Dado que envío datos incompletos o inválidos
- Cuando envío una solicitud POST a /users
- Entonces el sistema devuelve 400 Bad Request con detalles del error

### Escenario 4: Multi-tenant validation
- Dado que soy administrador de la empresa "empresa_uno"
- Cuando creo un usuario
- Entonces el usuario se asocia automáticamente a "empresa_uno" (no a otra empresa)
- Y el tercero se crea en el esquema "empresa_uno"

---

## ❌ Casos de Error

- Si el email ya existe en public.users → Error 409 Conflict
- Si los roles especificados no existen → Error 400 Bad Request
- Si el usuario no está autenticado → Error 401 Unauthorized
- Si los datos requeridos faltan → Error 400 Bad Request

---

## 🔐 Reglas de Negocio

- Email debe ser único en public.users
- El enterprise_id se toma del usuario autenticado (contexto JWT)
- Cada usuario debe tener al menos un rol
- Password debe tener mínimo 8 caracteres (si se especifica)
- El tercero se crea con is_employee = true por defecto
- El usuario se crea activo por defecto (active = true)

---

## 🎨 Consideraciones UI/UX

- Mostrar lista de roles disponibles en el formulario
- Validación en tiempo real del email (disponibilidad)
- Mensaje de éxito claro: "Usuario creado correctamente"
- Indicador de carga durante la creación
- Campos obligatorios marcados visualmente

---

## 📡 Requisitos Técnicos

- Endpoint: POST /users
- Método HTTP: POST
- Headers: Authorization (Bearer token), X-Tenant (optional, validated by middleware)
- Request body:
```json
{
  "email": "usuario@empresa.com",
  "name": "Juan Pérez",
  "password": "SecurePass123!",
  "roles": [2, 3], // IDs de roles
  "first_name": "Juan",
  "last_name": "Pérez",
  "document_number": "12345678",
  "document_type": "CC",
  "personal_email": "juan.perez@email.com"
}
```
- Response 201:
```json
{
  "data": {
    "id": 100,
    "enterprise_id": 50,
    "name": "Juan Pérez",
    "email": "usuario@empresa.com",
    "active": true,
    "created_at": "2026-03-15T10:00:00Z",
    "updated_at": "2026-03-15T10:00:00Z"
  },
  "success": true,
  "message": "Creado exitosamente"
}
```
- Códigos de error:
  - 400: Datos inválidos o faltantes
  - 401: No autenticado
  - 403: Sin permisos
  - 409: Email ya registrado
  - 500: Error interno

---

## 🧪 Criterios de Testing

- Unit tests: Service layer (create user, email validation, role assignment)
- Integration tests: Handler layer (HTTP request/response, error cases)
- E2E: Flujo completo de creación de usuario con tercero asociado

---

## 📎 Dependencias

- Servicio: UserService
- Librerías: bcrypt (password hashing), golang-migrate (migrations)
- Otra HU relacionada: HU-001 (Login), HU-003 (Roles)

---

## 🚫 Fuera de Alcance

- No incluye creación de empresa (ya existe)
- No incluye asignación de roles no existentes
- No incluye envío de email de bienvenida

---

## 🧠 Generación de Código

Requerir:
- modules/users/domain.go (entidad User y ThirdParty)
- modules/users/service.go (lógica de negocio)
- modules/users/repository.go (operaciones PostgreSQL)
- modules/users/handler.go (endpoints HTTP)
- modules/users/routes.go (registro de rutas)
- Migración para crear tercero (si no existe)
- Tests unitarios e integración
