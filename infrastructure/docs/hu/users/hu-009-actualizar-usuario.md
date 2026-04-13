# HU-009 - Actualizar datos de usuario

## 📌 Información General
- ID: HU-009
- Epic: EPIC-002 - Gestión de Usuarios de Empresa
- Prioridad: Media
- Estado: Backlog
- Porcentaje: 0%
- Autor: QA & Requirements Engineer
- Fecha: 2026-03-15

---

## 👤 Historia de Usuario

**Como** administrador de empresa
**Quiero** actualizar los datos de un usuario
**Para** mantener la información actualizada y correcta

---

## �ám Descripción Funcional

El sistema debe permitir actualizar los datos de un usuario existente. Los campos actualizables incluyen nombre, email (si no existe en otro usuario), información del tercero asociado.

---

## ✅ Criterios de Aceptación

### Escenario 1: Actualización exitosa
- Dado que soy administrador autenticado
- Cuando envío PUT /users/:id con datos válidos
- Entonces el sistema actualiza el usuario
- Y devuelve los datos actualizados

### Escenario 2: Email duplicado
- Dado que existe "usuario1@empresa.com" y "usuario2@empresa.com"
- Cuando intento actualizar usuario1 con email "usuario2@empresa.com"
- Entonces el sistema devuelve 409 Conflict

### Escenario 3: Usuario no encontrado
- Dado que solicito actualizar usuario inexistente
- Cuando envío PUT /users/99999
- Entonces el sistema devuelve 404 Not Found

### Escenario 4: Multi-tenant validation
- Dado que soy administrador de "empresa_uno"
- Cuando actualizo un usuario de "empresa_dos"
- Entonces el sistema devuelve 404 Not Found

### Escenario 5: Actualización parcial
- Dado que envío solo algunos campos
- Cuando envío PUT /users/:id
- Entonces el sistema actualiza solo los campos enviados

---

## ❌ Casos de Error

- Si el email ya existe → Error 409 Conflict
- Si el usuario no existe → Error 404 Not Found
- Si el usuario no pertenece a la empresa → Error 404 Not Found
- Si datos inválidos → Error 400 Bad Request

---

## 🔐 Reglas de Negocio

- Email debe ser único en public.users (excluyendo el usuario actual)
- Solo se pueden actualizar usuarios de la empresa autenticada
- El tercero asociado también se actualiza si se envían datos relevantes
- Password no se actualiza en este endpoint (usar endpoint específico si es necesario)

---

## 🎨 Consideraciones UI/UX

- Formulario con datos actuales precargados
- Validación en tiempo real del email
- Mensaje de éxito/ error claro
- Indicador de carga

---

## 📡 Requisitos Técnicos

- Endpoint: PUT /users/:id
- Método HTTP: PUT
- Path parameter: id (int64)
- Request body:
```json
{
  "email": "nuevo.email@empresa.com",
  "name": "Nuevo Nombre",
  "first_name": "Nuevo",
  "last_name": "Apellido",
  "document_number": "87654321",
  "document_type": "CC",
  "personal_email": "nuevo.personal@email.com"
}
```
- Response 200:
```json
{
  "data": {
    "id": 1,
    "enterprise_id": 50,
    "name": "Nuevo Nombre",
    "email": "nuevo.email@empresa.com",
    "active": true,
    "created_at": "2026-03-15T10:00:00Z",
    "updated_at": "2026-03-15T12:00:00Z"
  },
  "success": true,
  "message": "Actualizado exitosamente"
}
```
- Códigos de error:
  - 400: Datos inválidos
  - 404: Usuario no encontrado
  - 409: Email ya registrado
  - 500: Error interno

---

## 🧪 Criterios de Testing

- Unit tests: Service layer (update user, email validation)
- Integration tests: Handler layer (HTTP request/response, error cases)
- E2E: Flujo completo de actualización

---

## 📎 Dependencias

- Servicio: UserService
- Repositorio: UserRepository
- Otra HU relacionada: HU-006 (Crear usuario)

---

## 🚫 Fuera de Alcance

- No incluye actualización de password
- No incluye cambio de empresa
- No incluye eliminación de usuario

---

## 🧠 Generación de Código

Requerir:
- UserService.Update
- UserRepository.Update
- Handler para PUT /users/:id
