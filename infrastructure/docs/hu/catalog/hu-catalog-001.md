# HU-CATALOG-001: Crear Producto

## 📌 Información General
- ID: HU-CATALOG-001
- Epic: EPICA-016 (Gestión de Catálogo)
- Prioridad: Alta
- Estado: ✅ Completado (Actualizado)
- Porcentaje: 100%
- Autor: Equipo Backend Aura POS
- Fecha: 2026-04-16

---

## 👤 Historia de Usuario

**Como** administrador/gerente de una empresa  
**Quiero** crear nuevos productos en el catálogo  
**Para** poder venderlos posteriormente en las transacciones de venta

---

## 🧠 Descripción Funcional

El sistema debe permitir crear productos con información completa incluyendo SKU único, código de barras, nombre, descripción, categoría, marca, unidad de medida, tipo de producto, precios múltiples, configuración de impuestos, y opciones de control de inventario.

---

## ✅ Criterios de Aceptación

### Escenario 1: Crear producto exitosamente
- Dado que estoy autenticado como usuario con rol ADMIN o MANAGER
- Cuando envío una solicitud POST a `/products` con datos válidos del producto
- Entonces el sistema crea el producto con estado activo (active = true)
- Y retorna el producto creado con código 201

### Escenario 2: Crear producto con categoría y marca
- Dado que existen una categoría y marca activas en la empresa
- Cuando incluyo `categoriaId` y `marcaId` válidos en la solicitud
- Entonces el sistema asocia el producto a esa categoría y marca

### Escenario 3: Crear producto con precios múltiples
- Dado que proporciono precios en diferentes niveles (precio, precio2, precio3)
- Cuando creo el producto
- Entonces el sistema guarda todos los precios
- Y el precio principal es el precio de venta al público

### Escenario 4: Crear producto con tipo ESTANDAR
- Dado que el tipo de producto es ESTANDAR
- Cuando creo el producto
- Entonces el tipo por defecto es ESTANDAR

### Escenario 5: Crear producto con control de inventario
- Dado que `manejaInventario` es true
- Cuando creo el producto
- Entonces el sistema habilita el control de inventario
- Y el stock inicial es 0

---

## ❌ Casos de Error

- **SKU duplicado**: Si el SKU ya existe para la empresa, retorna error 400
- **Código de barras duplicado**: Si el barcode ya existe para la empresa, retorna error 400
- **Tipo de producto inválido**: Si `tipoProducto` no es ESTANDAR, SERVICIO, COMBO o RECETA, retorna error 400
- **Precio negativo**: Si `precio` o `costo` son negativos, retorna error 400
- **Categoría no encontrada**: Si `categoriaId` no existe, retorna error 400
- **Marca no encontrada**: Si `marcaId` no existe, retorna error 400
- **Unidad de medida no encontrada**: Si `unidadMedidaBaseId` no existe, retorna error 400
- **Campos requeridos faltantes**: Si `sku`, `nombre`, `costo`, `precio` o `unidadMedidaBaseId` faltan, retorna error 400

---

## 🔐 Reglas de Negocio

- SKU debe ser único por empresa
- Código de barras debe ser único por empresa (si se proporciona)
- Los precios no pueden ser negativos
- El tipo de producto debe ser uno de: ESTANDAR, SERVICIO, COMBO, RECETA
- El tipo por defecto es ESTANDAR
- Por defecto el producto está activo y visible en POS
- La unidad de medida es requerida
- El IVA por defecto es 19%
- Si maneja inventario, el stock inicial es 0

---

## 📡 Requisitos Técnicos

- **Endpoint**: `/products`
- **Método HTTP**: POST
- **Request**:
```json
{
  "nombre": "Producto tets",
  "sku": "sk-u1",
  "codigoBarras": "1123255241",
  "descripcion": "descripcion del producto",
  "imagenUrl": null,
  "categoriaId": 8,
  "marcaId": 4,
  "unidadMedidaBaseId": 6,
  "tipoProducto": "ESTANDAR",
  "activo": true,
  "precio": 18558,
  "costo": 17000,
  "precio2": 17500,
  "precio3": null,
  "ivaPorcentaje": 10,
  "impoconsumo": 5,
  "manejaInventario": true,
  "manejaLotes": false,
  "manejaSerial": false,
  "permitirStockNegativo": true,
  "visibleEnPos": true
}
```
- **Response**: 201 Created con el producto creado
- **Códigos de error**: 400 (bad request), 409 (conflicto SKU/barcode duplicado)

---

## 🧪 Criterios de Testing

- Unit tests para validaciones de servicio
- Tests de repositorio para inserción en DB
- Tests de handler para validación de request
- Tests de integración para flujo completo

---

## 📎 Dependencias

- Migración: 000003_products.up.sql
- Servicios: products.Service
- Módulo relacionado: EPICA-016 (Gestión de Catálogo)

---

## 📝 Notas de Actualización

**Fecha**: 2026-04-16  
**Cambios realizados**:
- Se agregaron nuevos campos: barcode, codigoBarras, tipoProducto, unidadMedidaBaseId, precio2, precio3, ivaPorcentaje, impoconsumo, manejaInventario, manejaLotes, manejaSerial, permitirStockNegativo, visibleEnPos
- Se modificó el campo `status` a `activo` (boolean)
- Se cambió `sale_price` a `precio` en el request
- Se cambió `cost_price` a `costo` en el request
- Se cambió `tax_rate` a `ivaPorcentaje` e `impoconsumo` en el request
- Se agregaron restricciones CHECK para product_type
- Se agregaron comentarios en español para cada columna