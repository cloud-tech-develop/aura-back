# EPICA-016: Gestión de Catálogo - Aura POS

## 📌 Información General

- ID: EPICA-016
- Estado: ✅ Completado
- Prioridad: Alta
- Fecha inicio: 2025-01-15
- Fecha objetivo: 2025-03-01
- Owner: Equipo Backend Aura POS
- Porcentaje: 100%

---

## 🎯 Objetivo de Negocio

Permitir a las empresas gestionar su catálogo de productos, incluyendo la administración de productos, categorías y marcas. Este módulo es fundamental para el funcionamiento del sistema de ventas, ya que todos los productos vendidos deben existir previamente en el catálogo.

**Problema resuelto**: Las empresas necesitan un sistema centralizado para gestionar su inventario de productos con información detallada (precios, SKU, categorías, marcas) que se utilize en el proceso de ventas y reportes.

**Valor generado**:
- Estandarización de productos dentro de cada empresa
- Facilidad para buscar y filtrar productos en ventas
- Control de precios y costos
- Trazabilidad mediante SKU único

---

## 👥 Stakeholders

- **Usuario final**: Administradores, gerentes y cajeros de las empresas
- **Equipo técnico**: Backend developers de Aura POS
- **Producto**: Product Manager de Aura POS

---

## 🧠 Descripción Funcional General

El módulo de catálogo permite la gestión completa de productos, categorías y marcas dentro del contexto multi-tenant del sistema. Cada empresa tiene acceso únicamente a su propio catálogo, con funcionalidades de CRUD completas y búsqueda filtrada.

**Características principales**:
- Productos con precios, SKU único, control de inventario
- Categorías jerárquicas para organizar productos
- Marcas para identificar fabricantes/proveedores
- Búsqueda y filtrado avanzado
- Sincronización de catálogos entre locales

---

## 📦 Alcance

### Incluye

- ✅ CRUD completo de productos
- ✅ CRUD completo de categorías
- ✅ CRUD completo de marcas
- ✅ CRUD completo de unidades
- ✅ Validación de SKU único por empresa
- ✅ Validación de precios (precio venta >= precio costo)
- ✅ Control de estado (ACTIVE, INACTIVE)
- ✅ Soft delete para productos
- ✅ Búsqueda de productos por nombre, SKU, categoría
- ✅ Endpoints RESTful para cada entidad

### No incluye

- ❌ Gestión de inventario (movimientos de stock)
- ❌ Imágenes de productos
- ❌ Variantes de productos (tallas, colores)
- ❌ Proveedores

---

## 🧩 Historias de Usuario Asociadas

| HU | Título | Estado |
|----|--------|--------|
| HU-CATALOG-001 | Crear producto | ✅ Completado |
| HU-CATALOG-002 | Listar productos | ✅ Completado |
| HU-CATALOG-003 | Obtener producto por ID | ✅ Completado |
| HU-CATALOG-004 | Actualizar producto | ✅ Completado |
| HU-CATALOG-005 | Eliminar producto | ✅ Completado |
| HU-CATALOG-006 | Crear categoría | ✅ Completado |
| HU-CATALOG-007 | Listar categorías | ✅ Completado |
| HU-CATALOG-008 | Obtener categoría por ID | ✅ Completado |
| HU-CATALOG-009 | Actualizar categoría | ✅ Completado |
| HU-CATALOG-010 | Eliminar categoría | ✅ Completado |
| HU-CATALOG-011 | Crear marca | ✅ Completado |
| HU-CATALOG-012 | Listar marcas | ✅ Completado |
| HU-CATALOG-013 | Obtener marca por ID | ✅ Completado |
| HU-CATALOG-014 | Actualizar marca | ✅ Completado |
| HU-CATALOG-015 | Eliminar marca | ✅ Completado |
| HU-CATALOG-016 | Crear Unidad | ✅ Completado |
| HU-CATALOG-017 | Listar Unidades | ✅ Completado |
| HU-CATALOG-018 | Obtener Unidad por ID | ✅ Completado |
| HU-CATALOG-019 | Actualizar Unidad | ✅ Completado |
| HU-CATALOG-020 | Eliminar Unidad | ✅ Completado |
| HU-CATALOG-021 | Listar Unidades Paginadas | ✅ Completado |

---

## 🔐 Reglas de Negocio

### Multi-tenancy

- Todos los datos del catálogo están scoped al tenant (empresa)
- Las consultas deben filtrar por `enterprise_id`
- El slug de la empresa se obtiene del JWT o subdominio

### Productos

- **SKU**: Debe ser único dentro de la empresa
- **Precio de venta**: Debe ser mayor o igual al precio de costo
- **Estado inicial**: Los productos se crean con estado ACTIVE
- **Stock inicial**: Los productos se crean con stock en 0
- **Soft delete**: La eliminación de productos es lógica (cambia status a DELETED)

### Categorías y Marcas

- **Nombre único**: No puede haber dos categorías con el mismo nombre en la misma empresa
- **Estado**: Soportan ACTIVE e INACTIVE
- **Eliminación**: Solo si no hay productos asociados (verificar en service)

### Permisos

- Solo usuarios con rol ADMIN o MANAGER pueden gestionar el catálogo
- Los cajeros pueden solo leer el catálogo

---

## 🧱 Arquitectura Relacionada

### Backend

- **Framework**: Go + Gin
- **Módulo**: `modules/catalog/`
- **Sub-módulos**: `products/`, `categories/`, `brands/`, `units/`

### Base de datos

- **Pattern**: Schema-per-tenant
- **Tablas**: `products`, `categories`, `brands` (en tenant schema)
- **Migración**: `000002_products.up.sql`

### Endpoints

```
/products
  ├── POST /products              → Crear producto
  ├── GET /products               → Listar productos
  ├── GET /products/:id          → Obtener producto
  ├── PUT /products/:id          → Actualizar producto
  └── DELETE /products/:id       → Eliminar producto (soft delete)

/categories
  ├── POST /categories            → Crear categoría
  ├── GET /categories            → Listar categorías
  ├── GET /categories/:id        → Obtener categoría
  ├── PUT /categories/:id        → Actualizar categoría
  └── DELETE /categories/:id     → Eliminar categoría

/brands
  ├── POST /brands                → Crear marca
  ├── GET /brands                → Listar marcas
  ├── GET /brands/:id            → Obtener marca
  ├── PUT /brands/:id            → Actualizar marca
  └── DELETE /brands/:id          → Eliminar marca

/units
  ├── POST /units                 → Crear unidad
  ├── GET /units                 → Listar unidades
  ├── GET /units/:id            → Obtener unidad
  ├── PUT /units/:id           → Actualizar unidad
  ├── DELETE /units/:id         → Eliminar unidad
  └── POST /units/page          → Listar paginado
```

---

## 📊 Métricas de Éxito

- ✅ Endpoints implementados y funcionando
- ✅ Tests unitarios pasando
- ✅ Validaciones de negocio aplicadas
- ✅ multi-tenant working correctly

---

## 🚧 Dependencias

### Entidades existentes

- **Empresa**: `enterprise_id` para scope multi-tenant
- **Usuario**: `user_id` para audit trail
- **Roles**: Permisos para gestión de catálogo

### Módulo relacionado

- **EPICA-002 (Gestión de Ventas)**: El catálogo es prerequisito para ventas

---

## 📝 Notas

- Los campos `global_id`, `sync_status`, `last_synced_at` están definidos en la estructura pero no persistidos actualmente
- La eliminación de categorías y marcas debe validar que no haya productos asociados
- El campo `category_id` en productos es opcional

---

## Resumen

- **Total de HU**: 21
- **Completadas**: 21
- **Pendientes**: 0
- **Sub-módulos implementados**: 4 (products, categories, brands, units)