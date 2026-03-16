# HU-001: Crear nueva empresa (Registro público)

**Como** usuario fundador  
**Quiero** registrar una nueva empresa en el sistema  
**Para** obtener acceso a mi propio entorno de POS

---

## Criterios de Aceptación

- [x] El endpoint `POST /enterprises` debe ser accesible sin autenticación
- [x] Debe aceptar: name, slug, email, password (obligatorios)
- [x] Debe aceptar: commercialName, subDomain, dv, phone, municipality (opcionales)
- [x] El slug debe cumplir regex `^[a-z0-9_]+$`
- [x] El slug debe tener entre 3 y 50 caracteres ✅
- [x] El email debe ser válido (formato válido)
- [x] La contraseña debe tener mínimo 8 caracteres ✅
- [x] Validar unicidad de slug, email y subdominio (ver HU-002)
- [x] Debe crear automáticamente el schema PostgreSQL con el nombre del slug
- [x] Debe crear el usuario administrador inicial en `public.users`
- [x] Debe asignar rol ADMIN al usuario fundador ✅
- [x] Debe ejecutar migraciones del tenant (tablas iniciales)
- [x] Debe crear un tercero inicial de tipo empleado
- [x] Retornar 201 con los datos de la empresa creada (sin password)

---

## Estado: ✅ 12/12 IMPLEMENTADO

---

## API Contract

### Request
```json
POST /enterprises
Content-Type: application/json

{
  "name": "Empresa de Ejemplo S.A.S.",
  "slug": "empresa_ejemplo",
  "email": "admin@empresa.com",
  "password": "securePassword123",
  "commercialName": "Mi Tienda",
  "subDomain": "mitienda",
  "dv": "12345678-9",
  "phone": "+573001234567",
  "municipality": "Bogotá"
}
```

### Response (201 Created)
```json
{
  "data": {...},
  "success": true,
  "message": "Creado exitosamente"
}
```

### Headers
```
X-Tenant-ID: 1
```

### Errores
- 400: Datos inválidos o faltantes
- 409: Slug, email o subdominio ya existe
