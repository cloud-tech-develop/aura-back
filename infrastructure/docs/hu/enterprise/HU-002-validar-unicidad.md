# HU-002: Validar unicidad de campos (slug, email, subdominio)

**Como** sistema  
**Quiero** verificar que slug, email y subdominio sean únicos antes de crear una empresa  
**Para** evitar conflictos de nombres y garantizar datos únicos en el sistema

---

## Criterios de Aceptación

### Slug
- [x] Verificar que el slug no exista en `public.enterprises`
- [x] El slug debe cumplir regex `^[a-z0-9_]+$`
- [x] Retornar 409 Conflict si el slug ya está en uso
- [x] Mensaje: "El identificador {slug} ya está en uso"

### Email
- [x] Verificar que el email no exista en `public.enterprises`
- [x] Verificar que el email no exista en `public.users` ✅
- [x] Retornar 409 Conflict si el email ya está registrado
- [x] Mensaje: "El correo electrónico {email} ya está registrado"

### Subdominio (opcional)
- [x] Si se proporciona subDomain, verificar unicidad en `public.enterprises`
- [x] El subDomain debe cumplir el mismo regex que el slug
- [x] Retornar 409 Conflict si el subdominio ya está en uso
- [x] Mensaje: "El subdominio {subdomain} ya está en uso"

### General
- [x] Todas las verificaciones dentro de la transacción
- [x] Retornar el primer error de validación encontrado

---

## Estado: ✅ 12/12 IMPLEMENTADO

---

## API Contract

### Response (409 Conflict - Slug)
```json
{
  "error": "conflict",
  "message": "El identificador empresa_ejemplo ya está en uso"
}
```

### Response (409 Conflict - Email)
```json
{
  "error": "conflict",
  "message": "El correo electrónico admin@empresa.com ya está registrado"
}
```

### Response (409 Conflict - Subdominio)
```json
{
  "error": "conflict",
  "message": "El subdominio mitienda ya está en uso"
}
```
