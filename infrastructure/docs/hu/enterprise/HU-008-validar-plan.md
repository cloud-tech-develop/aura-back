# HU-008: Validar límite de empresas según plan

**Como** sistema  
**Quiero** verificar la cuota del plan antes de crear una empresa  
**Para** controlar el uso de recursos

---

## Criterios de Aceptación

- [x] Antes de crear empresa, verificar el plan activo (**IMPLEMENTADO**)
- [x] Contar empresas actuales vs límite del plan (**IMPLEMENTADO**)
- [x] Retornar 403 Forbidden si se alcanzó el límite (**IMPLEMENTADO**)
- [x] Mensaje claro: "Ha alcanzado el límite de empresas de su plan" (**IMPLEMENTADO**)

---

## Estado: 4/4 implementados

### Issues Técnicos Resueltos
| ID | Descripción | Severidad |
|----|-------------|-----------|
| #008-001 | Crear validación de plan antes de crear empresa | Alta |
| #008-002 | Consultar límite de empresas en tabla plans | Alta |
| #008-003 | Retornar 403 cuando se alcance el límite | Media |

---

## Notas Técnicas

- La tabla `public.plans` ahora tiene el campo `max_enterprises`
- La validación se realiza en el servicio enterprise antes de crear
- Se consulta el plan del tenant y se cuentan las empresas actuales
- Si se alcanza el límite, se retorna error `ErrPlanLimitReached`

---

## Cambios Realizados

### Migration
- `000008_add_max_enterprises_to_plans.up.sql`: Agrega columna `max_enterprises` a `public.plans`

### Código
- `domain.go`: Agregada estructura `Plan` y métodos al `Repository`
- `repository.go`: Implementados `GetPlanByEnterpriseID` y `CountEnterprisesByTenant`
- `service.go`: Agregada validación de límite en `Create`
- `handler.go`: Retorna 403 cuando `ErrPlanLimitReached`

---

## API Contract

### Response (403 Forbidden)
```json
{
  "success": false,
  "message": "Ha alcanzado el límite de empresas de su plan",
  "data": null
}
```
