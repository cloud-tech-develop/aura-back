# HU-012 - Listar roles disponibles

## 📌 Información General
- ID: HU-012
- Epic: EPIC-002 - Gestión de Usuarios de Empresa
- Prioridad: Alta
- Estado: Backlog
- Porcentaje: 0%
- Autor: QA & Requirements Engineer
- Fecha: 2026-03-16

---

## 👤 Historia de Usuario

**Como** administrador de empresa
**Quiero** listar los roles disponibles para asignar
**Para** poder seleccionar los roles apropiados al crear o modificar usuarios

---

## 🧠 Descripción Funcional

El sistema debe proporcionar un endpoint para listar los roles disponibles. Solo se deben mostrar los roles de nivel inferior o igual al nivel del rol del usuario actual (jerarquía de roles).

---

## ✅ Criterios de Aceptación

### Escenario 1: Listar roles como ADMIN (nivel 1)
- Dado que soy ADMIN (nivel 1) autenticado
- Cuando solicito GET /roles
- Entonces el sistema retorna roles de nivel >= 1 (ADMIN, SUPERVISOR, USER, SELLER, CASHIER, ACCOUNTANT)
- Y NO retorna SUPERADMIN (nivel 0)

### Escenario 2: Listar roles como SUPERADMIN (nivel 0)
- Dado que soy SUPERADMIN (nivel 0) autenticado
- Cuando solicito GET /roles
- Entonces el sistema retorna todos los roles (SUPERADMIN, ADMIN, SUPERVISOR, USER, SELLER, CASHIER, ACCOUNTANT)

### Escenario 3: Listar roles como SUPERVISOR (nivel 2)
- Dado que soy SUPERVISOR (nivel 2) autenticado
- Cuando solicito GET /roles
- Entonces el sistema retorna roles de nivel >= 2 (SUPERVISOR, USER, SELLER, CASHIER, ACCOUNTANT)

### Escenario 4: Response exitosa
- El sistema retorna:
```json
{
  "data": [
    {
      "id": 2,
      "name": "ADMIN",
      "description": "Administrator",
      "level": 1
    },
    {
      "id": 3,
      "name": "SUPERVISOR",
      "description": "Supervisor",
      "level": 2
    }
  ],
  "success": true,
  "message": "Operación exitosa"
}
```

---

## ❌ Casos de Error

- Si no hay token de autenticación → Error 401 Unauthorized
- Si el token es inválido → Error 401 Unauthorized

---

## 🔐 Reglas de Negocio

- Solo se listan roles de nivel >= nivel del usuario actual
- Nivel 0 = SUPERADMIN (puede ver todos)
- Nivel 1 = ADMIN (puede ver nivel >= 1)
- Nivel 2 = SUPERVISOR (puede ver nivel >= 2)
- Nivel 3+ = USER/SELLER/CASHIER (solo ve nivel >= 3)

---

## 📡 Requisitos Técnicos

- Endpoint: GET /roles
- Método HTTP: GET
- Headers: Authorization: Bearer <token>
- Response 200:
```json
{
  "data": [
    {
      "id": 1,
      "name": "SUPERADMIN",
      "description": "Super admin",
      "level": 0
    }
  ],
  "success": true,
  "message": "Operación exitosa"
}
```
- Códigos de error:
  - 401: Token inválido o faltante
  - 500: Error interno

---

## 🧪 Criterios de Testing

- Unit tests: Service layer (filter roles by level)
- Integration tests: Handler layer (HTTP request/response)
- Verificar que ADMIN no vea SUPERADMIN

---

## 📎 Dependencias

- Middleware: AuthMiddleware (ya existe)
- Claims: RoleLevel en JWT (nuevo campo)

---

## 🚫 Fuera de Alcance

- No incluye creación de roles
- No incluye modificación de roles
- No incluye eliminación de roles
