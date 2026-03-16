# HU-006: Actualizar información de empresa

**Como** administrador de empresa  
**Quiero** actualizar los datos de mi empresa  
**Para** mantener la información actualizada

---

## Criterios de Aceptación

- [x] El endpoint `PUT /enterprises/:slug` requiere autenticación
- [x] Solo el propietario puede actualizar su empresa
- [x] Validar unicidad del email si se modifica
- [x] No permitir cambiar el slug una vez creado
- [x] Retornar 200 con los datos actualizados

---

## Estado: ✅ 5/5 IMPLEMENTADO

---

## API Contract

### Request
```
PUT /enterprises/empresa_uno
Content-Type: application/json

{
  "name": "Empresa Uno Modificada S.A.S.",
  "commercialName": "Mi Nueva Tienda",
  "phone": "+573009999999",
  "municipality": "Medellín"
}
```

### Response (200 OK)
```json
{
  "data": {
    "id": 1,
    "name": "Empresa Uno Modificada S.A.S.",
    "slug": "empresa_uno",
    ...
  }
}
```
