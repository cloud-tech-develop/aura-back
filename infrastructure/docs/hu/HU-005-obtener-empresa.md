# HU-005: Obtener empresa por slug

**Como** usuario autenticado  
**Quiero** obtener los detalles de una empresa por su slug  
**Para** visualizar la información de la empresa

---

## Criterios de Aceptación

- [x] El endpoint `GET /enterprises/:slug` requiere autenticación
- [x] Retornar 404 si la empresa no existe
- [x] Retornar los datos completos de la empresa (excepto password)

---

## Estado: ✅ 3/3 IMPLEMENTADO

---

## API Contract

### Request
```
GET /enterprises/empresa_uno
```

### Response (200 OK)
```json
{
  "data": {
    "id": 1,
    "tenantId": 1,
    "name": "Empresa Uno S.A.S.",
    "slug": "empresa_uno",
    ...
  }
}
```

### Response (404 Not Found)
```json
{
  "error": "not_found",
  "message": "Empresa no encontrada"
}
```
