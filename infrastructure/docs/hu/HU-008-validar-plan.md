# HU-008: Validar límite de empresas según plan

**Como** sistema  
**Quiero** verificar la cuota del plan antes de crear una empresa  
**Para** controlar el uso de recursos

---

## Criterios de Aceptación

- [ ] Antes de crear empresa, verificar el plan activo (**FALTA**)
- [ ] Contar empresas actuales vs límite del plan (**FALTA**)
- [ ] Retornar 403 Forbidden si se alcanzó el límite (**FALTA**)
- [ ] Mensaje claro: "Ha alcanzado el límite de empresas de su plan" (**FALTA**)

---

## Estado: 0/4 implementados

### Issues Técnicos a Resolver
| ID | Descripción | Severidad |
|----|-------------|-----------|
| #008-001 | Crear validación de plan antes de crear empresa | Alta |
| #008-002 | Consultar límite de empresas en tabla plans | Alta |
| #008-003 | Retornar 403 cuando se alcance el límite | Media |

---

## Notas Técnicas

- La tabla `public.plans` existe pero no se usa durante la creación
- No hay verificación de cuota antes de crear una empresa
- El campo `maxEnterprises` debería estar en la tabla plans

---

## API Contract

### Response (403 Forbidden)
```json
{
  "error": "forbidden",
  "message": "Ha alcanzado el límite de empresas de su plan"
}
```
