# HU-007: Cambiar estado de empresa (suspensión)

**Como** administrador del sistema  
**Quiero** suspender una empresa  
**Para** bloquear el acceso cuando sea necesario

---

## Criterios de Aceptación

- [x] El endpoint `PATCH /enterprises/:slug/status` requiere rol ADMIN
- [x] Permitir cambiar a estados: ACTIVE, INACTIVE, SUSPENDED, DEBT
- [x] Si la empresa está SUSPENDED, no puede autenticarse
- [x] Registrar el cambio en logs de auditoría

---

## Estado: ✅ 4/4 IMPLEMENTADO

---

## API Contract

### Request
```
PATCH /enterprises/empresa_uno/status
Content-Type: application/json

{
  "status": "SUSPENDED"
}
```

### Response (200 OK)
```json
{
  "data": {
    "id": 1,
    "slug": "empresa_uno",
    "status": "SUSPENDED",
    "updatedAt": "2026-03-15T12:00:00Z"
  }
}
```

### Estados Válidos
- ACTIVE
- INACTIVE
- SUSPENDED
- DEBT
