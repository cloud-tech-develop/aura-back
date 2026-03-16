# EPIC-002 - Gestión de Usuarios de Empresa

## 📌 Información General
- ID: EPIC-002
- Estado: Backlog
- Prioridad: Alta
- Fecha inicio: 2026-03-15
- Fecha objetivo: 2026-03-22
- Owner: QA & Requirements Engineer
- Porcentaje: 0%

---

## 🎯 Objetivo de Negocio

Gestionar usuarios adicionales para una empresa existente en Aura POS Backend.

¿Qué problema resuelve?
- Permite crear usuarios adicionales (no administradores) para empresas ya creadas
- Asigna roles específicos a cada usuario para control de acceso
- Genera automáticamente el tercero en el esquema del tenant para cada usuario

¿Qué valor genera?
- Mejor gestión de usuarios por empresa
- Control de acceso basado en roles
- Integración con módulo de terceros del tenant

---

## 👥 Stakeholders

- Usuario final: Administrador de empresa
- Equipo técnico: Backend developers
- Producto: Product Owner

---

## 🧠 Descripción Funcional General

Módulo de gestión de usuarios para empresas existentes en Aura POS. Permite crear, listar, actualizar y gestionar el estado de usuarios adicionales (no administradores iniciales). Cada usuario se asocia a una empresa específica y automáticamente se crea un tercero en el esquema del tenant.

---

## 📦 Alcance

Incluye:
- Creación de usuarios adicionales para empresas existentes
- Generación automática de tercero en esquema del tenant
- Asignación de roles a usuarios
- Validación de email único en public.users
- Listado de usuarios por empresa
- Actualización de usuarios
- Cambio de estado (activo/inactivo)
- Gestión de roles por usuario

No incluye:
- Creación de empresa (ya existe en EPIC-001)
- Autenticación de usuarios (ya existe en módulo auth)
- Eliminación física de usuarios (soft delete via deleted_at)

---

## 🧩 Historias de Usuario Asociadas

- [ ] HU-006 - Crear usuario adicional para empresa existente
- [ ] HU-007 - Listar usuarios de una empresa
- [ ] HU-008 - Obtener detalles de un usuario específico
- [ ] HU-009 - Actualizar datos de usuario
- [ ] HU-010 - Cambiar estado del usuario (activo/inactivo)
- [ ] HU-011 - Asignar roles a usuario

---

## 🐞 Bugs Asociados

- Ninguno identificado

---

## 🔐 Reglas de Negocio Globales

- Email debe ser único en public.users
- Cada usuario debe estar asociado a una empresa (enterprise_id)
- Al crear usuario, se genera automáticamente un tercero en el esquema del tenant
- Roles se definen en public.roles y se asignan via public.user_roles
- Soft delete: usuarios se marcan con deleted_at en lugar de eliminarse
- Usuarios inactivos no pueden autenticarse
- El usuario inicial (admin) se crea durante la creación de empresa
- Usuarios adicionales se crean después de la creación de empresa

---

## 🧱 Arquitectura Relacionada

Frontend: Componentes de gestión de usuarios en dashboard de administración
Backend: Módulo users en modules/users/
Base de datos: 
  - public.users (usuarios globales)
  - public.user_roles (asignación de roles)
  - tenant.third_parties (información del tercero)
Autenticación: JWT middleware (tenant/auth.go)

---

## 📊 Métricas de Éxito

- 100% de usuarios creados con tercero asociado
- Tiempo de respuesta < 200ms para operaciones CRUD
- Tasa de error < 1% en creación de usuarios
- Validación de email único correctamente implementada

---

## 🚧 Riesgos

- Concurrencia en creación de usuarios con mismo email (mitigado con UNIQUE constraint)
- Consistencia entre public.users y tenant.third_parties (mitigado con transacciones)
- Validación de roles disponibles (mitigado con lookup de roles existentes)
