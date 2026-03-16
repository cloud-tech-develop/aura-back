# Users API Specification

## Base URL
```
/api/v1
```

## Authentication
All endpoints require JWT authentication via `Authorization: Bearer <token>` header.

## Multi-tenancy
The system uses JWT claims to determine the tenant/enterprise context. All queries are scoped to the authenticated user's enterprise.

---

## Endpoints

### 1. POST /users
Create a new user for the authenticated enterprise.

**Request:**
```http
POST /api/v1/users
Authorization: Bearer <token>
Content-Type: application/json

{
  "email": "usuario@empresa.com",
  "name": "Juan Pérez",
  "password": "SecurePass123!",
  "roles": [2, 3],
  "first_name": "Juan",
  "last_name": "Pérez",
  "document_number": "12345678",
  "document_type": "CC",
  "personal_email": "juan.perez@email.com"
}
```

**Response 201 Created:**
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

**Response 400 Bad Request:**
```json
{
  "data": null,
  "success": false,
  "message": "Campos requeridos faltantes o inválidos"
}
```

**Response 409 Conflict:**
```json
{
  "data": null,
  "success": false,
  "message": "El email ya está registrado"
}
```

**Validation Rules:**
- `email`: Required, valid email format, unique across public.users
- `name`: Required, max 255 characters
- `password`: Required, min 8 characters
- `roles`: Required, array of valid role IDs
- `first_name`: Required for third party creation
- `last_name`: Required for third party creation
- `document_number`: Required, 5-20 alphanumeric characters
- `document_type`: Required, one of: CC, CE, NIT, PASSPORT

---

### 2. GET /users
List users for the authenticated enterprise with pagination and filters.

**Request:**
```http
GET /api/v1/users?page=1&limit=10&status=active
Authorization: Bearer <token>
```

**Parameters:**
- `page` (integer, default: 1): Page number
- `limit` (integer, default: 10, max: 100): Items per page
- `status` (string, optional): Filter by status ("active" or "inactive")

**Response 200 OK:**
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
          "document_number": "123456",
          "document_type": "CC"
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

### 3. GET /users/:id
Get details of a specific user.

**Request:**
```http
GET /api/v1/users/1
Authorization: Bearer <token>
```

**Response 200 OK:**
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

**Response 404 Not Found:**
```json
{
  "data": null,
  "success": false,
  "message": "Usuario no encontrado"
}
```

---

### 4. PUT /users/:id
Update user data.

**Request:**
```http
PUT /api/v1/users/1
Authorization: Bearer <token>
Content-Type: application/json

{
  "email": "nuevo.email@empresa.com",
  "name": "Nuevo Nombre",
  "first_name": "Nuevo",
  "last_name": "Apellido",
  "document_number": "87654321",
  "document_type": "CC",
  "personal_email": "nuevo.personal@email.com"
}
```

**Response 200 OK:**
```json
{
  "data": {
    "id": 1,
    "enterprise_id": 50,
    "name": "Nuevo Nombre",
    "email": "nuevo.email@empresa.com",
    "active": true,
    "created_at": "2026-03-15T10:00:00Z",
    "updated_at": "2026-03-15T12:00:00Z"
  },
  "success": true,
  "message": "Actualizado exitosamente"
}
```

**Response 409 Conflict:**
```json
{
  "data": null,
  "success": false,
  "message": "El email ya está registrado"
}
```

---

### 5. PATCH /users/:id/status
Change user status (active/inactive).

**Request:**
```http
PATCH /api/v1/users/1/status
Authorization: Bearer <token>
Content-Type: application/json

{
  "active": false
}
```

**Response 200 OK:**
```json
{
  "data": {
    "id": 1,
    "active": false,
    "updated_at": "2026-03-15T12:00:00Z"
  },
  "success": true,
  "message": "Actualizado exitosamente"
}
```

**Response 400 Bad Request:**
```json
{
  "data": null,
  "success": false,
  "message": "Estado inválido. Valor esperado: boolean"
}
```

---

### 6. PATCH /users/:id/roles
Assign roles to a user (replaces all existing roles).

**Request:**
```http
PATCH /api/v1/users/1/roles
Authorization: Bearer <token>
Content-Type: application/json

{
  "role_ids": [2, 3, 4]
}
```

**Response 200 OK:**
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
      },
      {
        "id": 4,
        "name": "USER",
        "description": "Standard user"
      }
    ]
  },
  "success": true,
  "message": "Roles actualizados exitosamente"
}
```

**Response 400 Bad Request:**
```json
{
  "data": null,
  "success": false,
  "message": "IDs de roles inválidos o inexistentes"
}
```

---

## Common Error Responses

### 401 Unauthorized
```json
{
  "data": null,
  "success": false,
  "message": "Token inválido o expirado"
}
```

### 403 Forbidden
```json
{
  "data": null,
  "success": false,
  "message": "Acceso denegado"
}
```

### 500 Internal Server Error
```json
{
  "data": null,
  "success": false,
  "message": "Error interno del servidor"
}
```

---

## Database Schema Reference

### public.users
| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | BIGSERIAL | PRIMARY KEY | User ID |
| enterprise_id | BIGINT | FOREIGN KEY | Reference to enterprise |
| email | TEXT | UNIQUE, NOT NULL | User email |
| name | TEXT | NOT NULL | User full name |
| password_hash | VARCHAR(255) | NOT NULL | Hashed password |
| active | BOOLEAN | NOT NULL DEFAULT TRUE | Active status |
| created_at | TIMESTAMPTZ | DEFAULT NOW() | Creation timestamp |
| deleted_at | TIMESTAMPTZ | NULL | Soft delete timestamp |

### public.user_roles
| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| user_id | BIGINT | FOREIGN KEY, PRIMARY KEY | User ID |
| role_id | BIGINT | FOREIGN KEY, PRIMARY KEY | Role ID |

### tenant.third_parties
| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | BIGSERIAL | PRIMARY KEY | Third party ID |
| user_id | BIGINT | FOREIGN KEY | Reference to public.users |
| first_name | VARCHAR(100) | | First name |
| last_name | VARCHAR(100) | | Last name |
| document_number | VARCHAR(50) | NOT NULL | Document number |
| document_type | VARCHAR(20) | NOT NULL | Document type |
| personal_email | VARCHAR(150) | | Personal email |
| tax_responsibility | VARCHAR(20) | NOT NULL | Tax responsibility |
| is_employee | BOOLEAN | NOT NULL DEFAULT FALSE | Employee flag |
| created_at | TIMESTAMPTZ | DEFAULT NOW() | Creation timestamp |
| deleted_at | TIMESTAMPTZ | NULL | Soft delete timestamp |

---

## Available Roles
| ID | Name | Description | Level |
|----|------|-------------|-------|
| 1 | SUPERADMIN | Super admin | 0 |
| 2 | ADMIN | Administrator | 1 |
| 3 | SUPERVISOR | Supervisor | 2 |
| 4 | USER | Standard user | 3 |
| 5 | SELLER | Sales and customers access | 3 |
| 6 | CASHIER | Cashier | 3 |
| 7 | ACCOUNTANT | Accountant | 3 |

---

## Multi-tenancy Implementation

The system ensures data isolation through:
1. JWT claims containing `enterprise_id` and `slug`
2. Middleware sets search_path to tenant schema
3. All queries include `WHERE enterprise_id = ?` or use tenant schema
4. Email uniqueness validated across all enterprises in public.users
