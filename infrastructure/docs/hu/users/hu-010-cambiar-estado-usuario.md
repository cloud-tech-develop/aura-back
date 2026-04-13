# HU-010 - Cambiar estado del usuario (activo/inactivo)

## 📌 Información General
- ID: HU-010
- Epic: EPIC-002 - Gestión de Usuarios de Empresa
- Prioridad: Media
- Estado: Backlog
- Porcentaje: 0%
- Autor: QA & Requirements Engineer
- Fecha: 2026-03-15

---

## 👤 Historia de Usuario

**Como** administrador de empresa
**Quiero** cambiar el estado de un usuario (activo/inactivo)
**Para** controlar el acceso de usuarios sin eliminarlos

---

## 🧠 Descripción Funcional

El sistema debe permitir activar o desactivar un usuario mediante un endpoint específico. Un usuario desactivado no podrá autenticarse.

---

## ✅ Criterios de Aceptación

### Escenario 1: Cambio de estado exitoso
- Dado que soy administrador autenticado
- Cuando envío PATCH /users/:id/status con "active": false
- Entonces el sistema cambia el estado del usuario a inactivo
- Y el usuario no podrá autenticarse

### Escenario 2: Activación de usuario
- Dado que un usuario está inactivo
- Cuando envío PATCH /users/:id/status con "active": true
- Entonces el sistema activa el usuario
- Y el usuario puede autenticarse nuevamente

### Escenario 3: Usuario no encontrado
- Dado que solicito cambiar estado de usuario inexistente
- Cuando envío PATCH /users/99999/status
- Entonces el sistema devuelve 404 Not Found

### Escenario 4: Multi-tenant validation
- Dado que soy administrador de "empresa_uno"
- Cuando cambio estado de usuario de "empresa_dos"
- Entonces el sistema devuelve 404 Not Found

---

## ❌ Casos de Error

- Si el usuario no existe → Error 404 Not Found
- Si el usuario no pertenece a la empresa → Error 404 Not Found
- Si estado inválido → Error 400 Bad Request

---

## 🔐 Reglas de Negocio

- Solo usuarios de la empresa autenticada pueden ser modificados
- El estado debe ser booleano (true/false)
- Usuario inactivo no puede autenticarse (verificado en login)
- No se elimina físicamente el usuario

---

## 🎨 Consideraciones UI/UX

- Toggle switch para cambio de estado
- Confirmación para cambio de estado
- Indicador visual del estado actual
- Mensaje de éxito claro

---

## 📡 Requisitos Técnicos

- Endpoint: PATCH /users/:id/status
- Método HTTP: PATCH
- Path parameter: id (int64)
- Request body:
```json
{
  "active": false
}
```
- Response 200:
```json
{
  "data": {
    "id": 1,
    "active": false,
    "updated_at": "2026-03-15T12:00:00Z"
  },
  "success": true,
  "message": "Actualizado exitosamente"
}
```
- Códigos de error:
  - 400: Estado inválido
  - 404: Usuario no encontrado
  - 500: Error interno

---

## 🧪 Criterios de Testing

- Unit tests: Service layer (change status)
- Integration tests: Handler layer (HTTP request/response)
- E2E: Cambio de estado y verificación en login

---

## 📎 Dependencias

- Servicio: UserService
- Repositorio: UserRepository
- Otra HU relacionada: HU-001 (Login)

---

## 🚫 Fuerra de Alcance

- No incluye eliminación física
- No incluye cambio de otros estados
- No incluye notificación al usuario

---

## 🧠 Generación de Código

Requerir:
- UserService.ChangeStatus
- UserRepository.UpdateStatus
- Handler para PATCH /users/:id/status
