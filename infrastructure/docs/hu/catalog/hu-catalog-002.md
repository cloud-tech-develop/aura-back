# HU-CATALOG-002: Listar Productos

## 📌 Información General
- ID: HU-CATALOG-002
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Alta
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2025-01-15

---

## 👤 Historia de Usuario

**Como** usuario del sistema (admin, manager, cajero)  
**Quiero** listar los productos del catálogo de mi empresa  
**Para** poder buscar y seleccionar productos para vender

---

## 🧠 Descripción Funcional

El sistema debe devolver una lista paginada de productos con opciones de búsqueda por nombre o SKU, y filtrado por categoría o marca.

---

## ✅ Criterios de Aceptación

### Escenario 1: Listar todos los productos
- Dado que la empresa tiene productos registrados
- Cuando solicito GET a `/products` sin parámetros
-Entonces el sistema retorna los primeros 10 productos ordenados por ID

### Escenario 2: Pagination
- Dado que existen más de 10 productos
- Cuando especifico `page=2&limit=20`
-Entonces el sistema retorna los productos 21-40

### Escenario 3: Búsqueda por nombre
- Dado que hay productos con "cafe" en el nombre
- Cuando busco con `search=cafe`
-Entonces el sistema retorna solo los productos que contienen "cafe"

### Escenario 4: Filtrar por categoría
- Dado que existen productos en diferentes categorías
- Cuando filtro con `category_id=1`
-Entonces el sistema retorna solo productos de esa categoría

### Escenario 5: Filtrar por marca
- Dado que existen productos de diferentes marcas
- Cuando filtro con `brand_id=2`
-Entonces el sistema retorna solo productos de esa marca

---

## ❌ Casos de Error

- **Empresa sin productos**: Retorna array vacío con metadata de paginación
- **Parámetros inválidos**: Si page o limit no son números, usar valores por defecto

---

## 🔐 Reglas de Negocio

- Solo se retornan productos de la empresa del usuario autenticado
- Los productos con status DELETED no se retornan
- La paginación tiene valores por defecto: page=1, limit=10

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/products`
- **Método HTTP**: GET
- **Query Parameters**:
  - `page` (opcional): número de página, default 1
  - `limit` (opcional): elementos por página, default 10
  - `search` (opcional): búsqueda por nombre o SKU
  - `category_id` (opcional): filtrar por categoría
  - `brand_id` (opcional): filtrar por marca
- **Response**:
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 50,
    "total_pages": 5
  }
}
```

---

## 🧪 Criterios de Testing

- Tests de paginación
- Tests de filtros (search, category, brand)
- Tests de autorización multi-tenant