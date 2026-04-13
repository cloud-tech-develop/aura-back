# Users Module Requirements - Aura POS Backend

## Overview
This document contains the complete requirements for the Users module in Aura POS Backend, including Epic, User Stories, API specifications, and business validations.

## Files Created

### 1. Epic
- **File**: `infrastructure/docs/epics/epic-users-management.md`
- **ID**: EPIC-002
- **Title**: Gestión de Usuarios de Empresa
- **Description**: Gestionar usuarios adicionales para una empresa existente

### 2. User Stories (HU)
- **HU-006**: Crear usuario adicional para empresa existente
  - File: `infrastructure/docs/hu/hu-006-crear-usuario-adicional.md`
- **HU-007**: Listar usuarios de una empresa
  - File: `infrastructure/docs/hu/hu-007-listar-usuarios-empresa.md`
- **HU-008**: Obtener detalles de un usuario específico
  - File: `infrastructure/docs/hu/hu-008-obtener-usuario.md`
- **HU-009**: Actualizar datos de usuario
  - File: `infrastructure/docs/hu/hu-009-actualizar-usuario.md`
- **HU-010**: Cambiar estado del usuario (activo/inactivo)
  - File: `infrastructure/docs/hu/hu-010-cambiar-estado-usuario.md`
- **HU-011**: Asignar roles a usuario
  - File: `infrastructure/docs/hu/hu-011-asignar-roles-usuario.md`

### 3. API Specifications
- **File**: `infrastructure/docs/api/users-api-spec.md`
- Endpoints:
  - POST /users (create user)
  - GET /users (list users)
  - GET /users/:id (get user)
  - PUT /users/:id (update user)
  - PATCH /users/:id/status (change status)
  - PATCH /users/:id/roles (assign roles)

### 4. Business Validations
- **File**: `infrastructure/docs/business/user-validations.md`
- Validations:
  - Email uniqueness
  - Password requirements
  - Role validation
  - Third party creation
  - Enterprise association
  - User status
  - Multi-tenancy
  - Soft delete
  - Transaction safety
  - Request validation

## Key Features

### 1. User Creation
- Creates user in `public.users` table
- Automatically creates third party in tenant schema
- Assigns specified roles
- Validates email uniqueness across all enterprises
- Uses transaction for atomicity

### 2. Multi-tenancy
- Users are scoped to enterprise via JWT claims
- All queries include `enterprise_id` filter
- Cannot access users from other enterprises
- Tenant schema used for third party data

### 3. Role Management
- Roles defined in `public.roles`
- Assignment via `public.user_roles`
- Cannot assign SUPERADMIN to regular users
- Roles can be updated (replaces all existing roles)

### 4. Status Management
- Users can be active/inactive
- Inactive users cannot authenticate
- Soft delete via `deleted_at` timestamp

## Database Schema

### public.users
```sql
CREATE TABLE public.users (
    id             BIGSERIAL PRIMARY KEY,
    enterprise_id  BIGINT NOT NULL REFERENCES public.enterprises(id),
    email          TEXT UNIQUE NOT NULL,
    name           TEXT NOT NULL,
    password_hash  VARCHAR(255) NOT NULL,
    active         BOOLEAN NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    deleted_at     TIMESTAMPTZ DEFAULT NULL
);
```

### public.user_roles
```sql
CREATE TABLE public.user_roles (
    user_id     BIGINT NOT NULL REFERENCES public.users(id),
    role_id     BIGINT NOT NULL REFERENCES public.roles(id),
    PRIMARY KEY (user_id, role_id)
);
```

### tenant.third_parties
```sql
CREATE TABLE third_parties (
    id                  BIGSERIAL PRIMARY KEY,
    user_id             BIGINT REFERENCES public.users(id),
    first_name          VARCHAR(100),
    last_name           VARCHAR(100),
    document_number     VARCHAR(50) NOT NULL,
    document_type       VARCHAR(20) NOT NULL,
    personal_email      VARCHAR(150),
    tax_responsibility  VARCHAR(20) NOT NULL,
    is_employee         BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ DEFAULT NULL
);
```

## Implementation Steps

### 1. Create Module Structure
```
modules/users/
├── domain.go          # Entity definitions
├── service.go         # Business logic
├── repository.go      # Database operations
├── handler.go         # HTTP handlers
├── routes.go          # Route registration
└── logger.go          # Event logging
```

### 2. Create Migrations (if needed)
- Check if third_parties table needs updates
- Add any new columns if required

### 3. Implement Domain Entities
- User struct with all fields
- ThirdParty struct for tenant schema
- Repository and Service interfaces

### 4. Implement Service Layer
- Create user with transaction
- List users with pagination
- Get user by ID
- Update user
- Change status
- Assign roles

### 5. Implement Repository Layer
- PostgreSQL operations with context
- Parameterized queries
- Transaction support
- Multi-tenant queries

### 6. Implement HTTP Handlers
- Gin handlers for all endpoints
- Request validation
- Response formatting
- Error handling

### 7. Register Routes
- Add routes to server
- Apply auth middleware
- Apply tenant middleware

### 8. Write Tests
- Unit tests for service layer
- Integration tests for handlers
- E2E tests for complete flows

## Validation Checklist

- [ ] Email uniqueness across all enterprises
- [ ] Password minimum 8 characters
- [ ] Role IDs must exist
- [ ] Third party created automatically
- [ ] Multi-tenancy enforced
- [ ] JWT authentication required
- [ ] Transaction safety
- [ ] Soft delete implemented
- [ ] Status validation
- [ ] Request validation
- [ ] Error handling
- [ ] Performance optimization

## API Examples

### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@empresa.com",
    "name": "Juan Pérez",
    "password": "SecurePass123!",
    "roles": [2, 3],
    "first_name": "Juan",
    "last_name": "Pérez",
    "document_number": "12345678",
    "document_type": "CC",
    "personal_email": "juan.perez@email.com"
  }'
```

### List Users
```bash
curl -X GET "http://localhost:8080/api/v1/users?page=1&limit=10&status=active" \
  -H "Authorization: Bearer <token>"
```

### Get User
```bash
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer <token>"
```

### Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "nuevo.email@empresa.com",
    "name": "Nuevo Nombre"
  }'
```

### Change Status
```bash
curl -X PATCH http://localhost:8080/api/v1/users/1/status \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"active": false}'
```

### Assign Roles
```bash
curl -X PATCH http://localhost:8080/api/v1/users/1/roles \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"role_ids": [2, 3, 4]}'
```

## Dependencies

### Internal
- `shared/response` - HTTP response helpers
- `shared/errors` - Error definitions
- `tenant/auth` - JWT authentication
- `tenant/manager` - Tenant management

### External
- `github.com/gin-gonic/gin` - HTTP framework
- `github.com/lib/pq` - PostgreSQL driver
- `golang.org/x/crypto/bcrypt` - Password hashing
- `github.com/golang-jwt/jwt/v5` - JWT tokens

## Security Considerations

1. **Password Security**: Always hash with bcrypt
2. **JWT Security**: Validate signature and expiration
3. **SQL Injection**: Use parameterized queries
4. **Data Exposure**: Never return sensitive data
5. **Multi-tenancy**: Enforce enterprise isolation

## Performance Considerations

1. **Indexes**: Create indexes on frequently queried columns
2. **Pagination**: Always paginate list endpoints
3. **Transactions**: Use for multi-step operations
4. **Query Optimization**: Avoid N+1 queries

## Next Steps

1. Review and approve requirements
2. Create technical design document
3. Implement domain entities
4. Implement service layer
5. Implement repository layer
6. Implement HTTP handlers
7. Register routes
8. Write tests
9. Code review
10. Deploy to staging
