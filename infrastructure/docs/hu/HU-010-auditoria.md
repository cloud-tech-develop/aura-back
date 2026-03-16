# HU-010: Logging de auditoría de creación de empresa

**Como** administrador  
**Quiero** registrar quién creó cada empresa  
**Para** mantener trazabilidad y seguridad

---

## Criterios de Aceptación

- [ ] Registrar: timestamp, IP, email del fundador, slug de empresa (**FALTA**)
- [ ] Guardar en tabla de auditoría o logs estructurados (**FALTA**)
- [ ] Incluir en la respuesta el ID del tenant creado (**FALTA**)

---

## Estado: 0/3 implementados

### Issues Técnicos a Resolver
| ID | Descripción | Severidad |
|----|-------------|-----------|
| #010-001 | Agregar logging estructurado al crear empresa | Media |
| #010-002 | Capturar IP del cliente | Media |

---

## Datos a Registrar

| Campo | Descripción |
|-------|-------------|
| timestamp | Fecha y hora ISO 8601 |
| ip_address | IP del cliente |
| email | Email del usuario fundador |
| slug | Slug de la empresa creada |
| tenant_id | ID del tenant creado |
| action | "enterprise_created" |

---

## Notas Técnicas

- No hay logging estructurado al crear empresa
- El service publica eventos pero no registra auditoría

---

## API Contract - Response Actual



### Header Necesario
```
X-Tenant-ID: 1
```
