# HU-009: Resolución de tenant por token JWT

**Como** usuario  
**Quiero** que el sistema determine mi tenant automáticamente desde el token JWT  
**Para** no tener que especificar el tenant en cada solicitud

---

## Criterios de Aceptación

- [x] El token JWT debe incluir `tenant_id` y `slug`
- [x] El middleware de autenticación debe extraer el tenant del token
- [x] El middleware de tenant debe usar el slug del token (prioridad sobre header)
- [x] Si no hay token, usar subdominio como fallback

---

## Estado: ✅ 4/4 IMPLEMENTADO

---

## Flujo de Resolución de Tenant

```
1. AuthMiddleware valida el JWT y extrae el slug del token
2. TenantMiddleware obtiene el slug del contexto (del token)
3. Si no hay slug en token, usa subdominio como fallback
4. Valida que la empresa exista y esté ACTIVE
5. Establece el contexto de tenant
```

### Prioridad de Resolución
1. **JWT Token** (slug del token) - Mayor prioridad
2. **Subdominio** (fallback)
3. **Error** si no se encuentra ningún tenant

---

## Estructura del Token JWT

```json
{
  "user_id": 1,
  "enterprise_id": 1,
  "tenant_id": 1,
  "slug": "empresa_uno",
  "email": "admin@empresa.com",
  "roles": ["ADMIN"],
  "ip": "192.168.1.1",
  "exp": 1700000000,
  "iat": 1700000000
}
```

---

## Código - Middleware de Tenant

```go
func Middleware(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Primero obtener tenant del token JWT
        slug, hasSlug := SlugFromContext(c)
        
        // 2. Si no está en token, usar subdominio
        if !hasSlug || slug == "" {
            slug = resolveSubDomain(c)
        }
        
        // 3. Validar que existe y está ACTIVE
        // ...
        
        c.Set(string(TenantKey), slug)
        c.Next()
    }
}
```

---

## Notas Técnicas

- El header `X-Tenant` ya no es necesario para endpoints autenticados
- El token JWT contiene toda la información del tenant
- El subdominio se usa como fallback cuando no hay autenticación
- La validación de empresa ACTIVE se mantiene
