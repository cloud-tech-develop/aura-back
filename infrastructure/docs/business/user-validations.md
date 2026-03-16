# Business Validations for Users Module

## 1. Email Validation

### Rule
- Email must be unique across all enterprises in `public.users`
- Email format must be valid
- Email is required when creating a user

### Implementation
- **Database Level**: UNIQUE constraint on `public.users.email`
- **Service Level**: Check email existence before insertion
- **Error Response**: 409 Conflict - "El email ya está registrado"

### Validation Flow
```
1. User submits email
2. Service checks: SELECT id FROM public.users WHERE email = $1 AND deleted_at IS NULL
3. If found → Return conflict error
4. If not found → Proceed with creation
```

---

## 2. Password Validation

### Rule
- Password is required when creating a user
- Minimum length: 8 characters
- Must be hashed using bcrypt before storage

### Implementation
- **Service Level**: Validate length before hashing
- **Repository Level**: Store hashed password only
- **Error Response**: 400 Bad Request - "La contraseña debe tener al menos 8 caracteres"

### Validation Flow
```
1. User submits password
2. Service validates: len(password) >= 8
3. If invalid → Return validation error
4. If valid → Hash with bcrypt.GenerateFromPassword
5. Store hash in database
```

---

## 3. Role Validation

### Rule
- User must have at least one role
- Roles must exist in `public.roles`
- Role IDs must be valid integers
- Cannot assign SUPERADMIN role (level 0) to regular users

### Implementation
- **Database Level**: FOREIGN KEY constraint on `public.user_roles.role_id`
- **Service Level**: Validate role IDs exist before assignment
- **Error Response**: 400 Bad Request - "IDs de roles inválidos o inexistentes"

### Validation Flow
```
1. User submits role_ids array
2. For each role_id:
   a. Check existence in public.roles
   b. Verify role is not SUPERADMIN (unless creating admin user)
3. If any role_id invalid → Return validation error
4. If valid → Insert into public.user_roles
```

### Available Roles
| ID | Name | Level | Can Assign |
|----|------|-------|------------|
| 1 | SUPERADMIN | 0 | No (only during enterprise creation) |
| 2 | ADMIN | 1 | Yes |
| 3 | SUPERVISOR | 2 | Yes |
| 4 | USER | 3 | Yes |
| 5 | SELLER | 3 | Yes |
| 6 | CASHIER | 3 | Yes |
| 7 | ACCOUNTANT | 3 | Yes |

---

## 4. Third Party Creation Validation

### Rule
- When creating a user, a third party must be created in the tenant schema
- Third party requires: first_name, last_name, document_number, document_type
- Document number must be 5-20 alphanumeric characters
- Document type must be one of: CC, CE, NIT, PASSPORT
- is_employee defaults to true for user-associated third parties

### Implementation
- **Database Level**: Validation constraints in tenant.third_parties
- **Service Level**: Transaction ensures both user and third party are created
- **Error Response**: 400 Bad Request for invalid data

### Validation Flow
```
1. Begin transaction
2. Create user in public.users
3. Get user ID
4. Create third party in tenant.third_parties
5. If any step fails → Rollback transaction
6. If all succeed → Commit transaction
```

### Document Validation
- **document_number**: 5-20 alphanumeric characters or hyphens
- **document_type**: Must match predefined types
- **tax_responsibility**: Must be RESPONSIBLE or NOT-RESPONSIBLE

---

## 5. Enterprise Association Validation

### Rule
- Every user must be associated with exactly one enterprise
- enterprise_id is taken from JWT claims (authenticated user's enterprise)
- Cannot create users for other enterprises

### Implementation
- **Middleware**: Extract enterprise_id from JWT and set in context
- **Service Level**: Use enterprise_id from context, not from request body
- **Repository Level**: Always include enterprise_id in WHERE clauses

### Validation Flow
```
1. User authenticates → JWT contains enterprise_id
2. Auth middleware extracts enterprise_id from JWT
3. Service gets enterprise_id from context (not request)
4. User created with enterprise_id from context
5. Query to list users filters by enterprise_id
```

---

## 6. User Status Validation

### Rule
- Users can be active (active = true) or inactive (active = false)
- Inactive users cannot authenticate
- Default status is active = true
- Status can be changed via PATCH /users/:id/status

### Implementation
- **Database Level**: Boolean column with default TRUE
- **Auth Layer**: Check user.active in login function
- **Error Response**: 403 Forbidden - "Usuario inactivo"

### Validation Flow (Login)
```
1. User attempts login
2. Query: SELECT active FROM public.users WHERE email = $1
3. If active = false → Return 403 Forbidden
4. If active = true → Continue authentication
```

### Validation Flow (Status Change)
```
1. User submits {"active": false}
2. Validate boolean value
3. Update user.active in database
4. Return updated user
```

---

## 7. Multi-tenancy Validation

### Rule
- Users can only see and modify users from their own enterprise
- Query results must be scoped to enterprise_id from JWT
- Cannot access users from other enterprises (returns 404)

### Implementation
- **Middleware**: Sets search_path to tenant schema
- **Service Layer**: Always include enterprise_id in queries
- **Repository Layer**: Add WHERE enterprise_id = ? to all queries

### Validation Flow
```
1. Request arrives with JWT token
2. Auth middleware validates token and extracts enterprise_id
3. Service receives enterprise_id from context
4. All repository queries include: WHERE enterprise_id = ?
5. If user tries to access other enterprise → Returns 404
```

### Example Query
```sql
-- List users (safe - only own enterprise)
SELECT * FROM public.users 
WHERE enterprise_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- Get user by ID (safe - validates enterprise)
SELECT * FROM public.users 
WHERE id = $1 AND enterprise_id = $2 AND deleted_at IS NULL;
```

---

## 8. Soft Delete Validation

### Rule
- Users are never physically deleted
- Deleted users have deleted_at timestamp set
- Deleted users are excluded from all queries
- Cannot update deleted users

### Implementation
- **Database Level**: deleted_at TIMESTAMPTZ column (NULL = active)
- **Repository Level**: Always include `AND deleted_at IS NULL` in queries
- **Service Level**: Cannot update deleted users

### Validation Flow
```
1. All SELECT queries include: AND deleted_at IS NULL
2. When "deleting": UPDATE users SET deleted_at = NOW()
3. Cannot update user with deleted_at IS NOT NULL
```

---

## 9. Transaction Safety

### Rule
- Creating a user requires multiple database operations
- All operations must succeed or all must fail (atomicity)
- User and third party must be created together

### Implementation
- **Service Level**: Use database transactions
- **Error Handling**: Rollback on any failure
- **Retry Logic**: Not implemented (use idempotent operations)

### Validation Flow
```
1. Begin transaction
2. Create user in public.users
3. Get user ID (returning id)
4. Create third party in tenant.third_parties
5. Assign roles in public.user_roles
6. If any error → Rollback
7. If all succeed → Commit
```

---

## 10. Request Validation

### Rule
- All request bodies must be validated
- Required fields must be present
- Field formats must be correct
- Business rules must be enforced

### Implementation
- **Handler Layer**: Gin binding with validation tags
- **Service Layer**: Additional business rule validation
- **Error Response**: 400 Bad Request with specific error messages

### Required Fields by Endpoint

**POST /users:**
- email (required, valid format)
- name (required)
- password (required, min 8 chars)
- roles (required, array)
- first_name (required)
- last_name (required)
- document_number (required, 5-20 chars)
- document_type (required, valid type)

**PUT /users/:id:**
- email (optional, valid format if provided)
- name (optional)
- first_name (optional)
- last_name (optional)
- document_number (optional, 5-20 chars if provided)
- document_type (optional, valid type if provided)

**PATCH /users/:id/status:**
- active (required, boolean)

**PATCH /users/:id/roles:**
- role_ids (required, array of integers)

---

## 11. Error Handling Summary

| Error Type | HTTP Code | Message | Cause |
|------------|-----------|---------|-------|
| Email duplicado | 409 | "El email ya está registrado" | Email exists in public.users |
| Datos inválidos | 400 | "Campos requeridos faltantes o inválidos" | Missing/invalid fields |
| Usuario no encontrado | 404 | "Usuario no encontrado" | User not found or not in enterprise |
| Roles inválidos | 400 | "IDs de roles inválidos o inexistentes" | Role IDs don't exist |
| No autenticado | 401 | "Token inválido o expirado" | Missing/invalid JWT |
| Sin permisos | 403 | "Acceso denegado" | Trying to access other enterprise |
| Estado inválido | 400 | "Valor esperado: boolean" | Invalid status value |
| Error interno | 500 | "Error interno del servidor" | Database/other errors |

---

## 12. Security Considerations

### Password Handling
- Never store plain text passwords
- Always hash with bcrypt before storage
- Never return password hash in API responses
- Validate password strength (min 8 chars)

### JWT Claims
- Include enterprise_id in JWT claims
- Validate JWT signature and expiration
- Include user_id in JWT claims
- Include roles in JWT claims

### Data Exposure
- Never return password_hash in responses
- Never return internal IDs in third_party if not needed
- Filter sensitive data from responses

### SQL Injection Prevention
- Always use parameterized queries
- Never concatenate user input into SQL
- Use prepared statements

---

## 13. Performance Considerations

### Query Optimization
- Create indexes on:
  - public.users(email)
  - public.users(enterprise_id)
  - public.users(deleted_at)
  - public.user_roles(user_id, role_id)
  - tenant.third_parties(user_id)

### Pagination
- Always paginate list endpoints
- Use LIMIT/OFFSET for pagination
- Return pagination metadata

### Caching
- Consider caching role lookups
- Consider caching enterprise info from JWT
- Do not cache user lists (always fresh)

---

## 14. Testing Validations

### Unit Tests
- Email validation (unique, format)
- Password validation (length, hashing)
- Role validation (existence, assignment)
- Third party creation validation
- Multi-tenancy validation

### Integration Tests
- HTTP request/response validation
- Transaction rollback on error
- JWT authentication and authorization
- Database constraint enforcement

### E2E Tests
- Complete user creation flow
- User authentication with new user
- Role-based access control
- Multi-tenant data isolation
