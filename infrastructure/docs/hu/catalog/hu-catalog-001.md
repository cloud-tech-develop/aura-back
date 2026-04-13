# HU-CATALOG-001: Crear Producto

## 📌 Información General
- ID: HU-CATALOG-001
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Alta
- Estado: ✅ Completado
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2025-01-15

---

## 👤 Historia de Usuario

**Como** administrador/gerente de una empresa  
**Quiero** crear nuevos productos en el catálogo  
**Para** poder venderlos posteriormente en las transacciones de venta

---

## 🧠 Descripción Funcional

El sistema debe permitir crear productos con información completa incluyendo SKU único, nombre, descripción, categoría, marca, precios de costo y venta, tasa de impuestos, stock mínimo e imagen.

---

## ✅ Criterios de Aceptación

### Escenario 1: Crear producto exitosamente
- Dado que estoy autenticado como usuario con rol ADMIN o MANAGER
- Cuando envío una solicitud POST a `/products` con datos válidos del producto
- Entonces el sistema crea el producto con estado ACTIVE
- Y retorna el producto creado con código 201
- Y el stock inicial es 0

### Escenario 2: Crear producto con categoría y marca
- Dado que existen una categoría y marca activas en la empresa
- Cuando incluyo `category_id` y `brand_id` válidos en la solicitud
- Entonces el sistema asocia el producto a esa categoría y marca

### Escenario 3: Crear producto con precio mayor al costo
- Dado que proporciono `sale_price` mayor a `cost_price`
- Cuando creo el producto
- Entonces el sistema acepta la operación exitosamente

---

## ❌ Casos de Error

- **SKU duplicado**: Si el SKU ya existe para la empresa, retorna error 400
- **Precio venta menor al costo**: Si `sale_price < cost_price`, retorna error 400
- **Categoría no encontrada**: Si `category_id` no existe, retorna error 400
- **Marca no encontrada**: Si `brand_id` no existe, retorna error 400
- **Campos requeridos faltantes**: Si `sku`, `name`, `cost_price` o `sale_price` faltan, retorna error 400

---

## 🔐 Reglas de Negocio

- SKU debe ser único por empresa
- `sale_price` debe ser >= `cost_price`
- El estado inicial del producto es ACTIVE
- El stock inicial es 0
- `category_id` y `brand_id` son opcionales

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/products`
- **Método HTTP**: POST
- **Request**:
```json
{
  "sku": "SKU-001",
  "name": "Producto Ejemplo",
  "description": "Descripción del producto",
  "category_id": 1,
  "brand_id": 1,
  "cost_price": 100.00,
  "sale_price": 150.00,
  "tax_rate": 19.0,
  "min_stock": 10,
  "image_url": "https://..."
}
```
- **Response**: 201 Created con el producto creado
- **Códigos de error**: 400 (bad request), 409 (conflicto SKU duplicado)

---

## 🧪 Criterios de Testing

- Unit tests para validaciones de servicio
- Tests de repositorio para inserción en DB
- Tests dehandler para validación de request
- Tests de integración para flujo completo

---

## 📎 Dependencias

- Migración: 000002_products.up.sql
- Servicios: products.Service
- Módulo relacionado: EPICA-002 (Gestión de Ventas)