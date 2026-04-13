# HU-004: Listar empresas (Admin)

**Como** administrador del sistema  
**Quiero** ver todas las empresas registradas  
**Para** gestionar el宏观层面 del sistema

---

## Criterios de Aceptación

- [x] El endpoint `GET /enterprises` requiere autenticación JWT
- [x] Retornar lista paginada de empresas
- [x] Incluir filtros por estado (ACTIVE, INACTIVE, SUSPENDED)
- [x] Solo accesible para usuarios con rol ADMIN del tenant public

---

## Estado: ✅ 4/4 IMPLEMENTADO

---

## API Contract

### Request
```
GET /enterprises?page=1&limit=10&status=ACTIVE
```

### Response (200 OK)
```json
{
  "data": [
    {
      "id": 1,
      "name": "Empresa Uno",
      "slug": "empresa_uno",
      "email": "admin@empresa1.com",
      "status": "ACTIVE",
      "createdAt": "2026-01-15T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "totalPages": 3
  }
}
```
