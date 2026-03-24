# EPICA-002: Gestión de Ventas - Aura POS

## Descripción

Esta épica abarca todo el flujo de ventas en el sistema Aura POS, desde la gestión del catálogo de productos hasta la generación de reportes de ventas. El sistema permite a los cajeros y gerentes procesar ventas de manera eficiente, gestionar inventario y generar documentación fiscal.

---

## Arquitectura del Módulo de Ventas

### Estructura de Módulos

```
modules/
├── products/          # Gestión de productos, categorías, marcas
├── cart/             # Carrito de compras y cotizaciones
├── sales/            # Órdenes de venta y procesamiento
├── payments/         # Procesamiento de pagos
├── invoices/         # Generación de facturas
└── reports/          # Reportes y análisis de ventas
```

### Flujo de Venta Completo

```
1. Gestión de Catálogo (products/)
   ├── Crear/actualizar productos
   ├── Gestionar categorías y marcas
   └── Control de inventario

2. Carrito de Compras (cart/)
   ├── Añadir/eliminar productos
   ├── Aplicar descuentos
   └── Calcular totales e impuestos

3. Procesamiento de Venta (sales/)
   ├── Convertir carrito a orden
   ├── Actualizar inventario
   └── Crear registro de venta

4. Pagos (payments/)
   ├── Múltiples métodos de pago
   ├── Cálculo de cambio
   └── Registro de transacciones

5. Facturación (invoices/)
   ├── Generación automática
   ├── Numeración secuencial
   └── Documentación fiscal

6. Reportes (reports/)
   ├── Análisis de ventas
   ├── Métricas de desempeño
   └── Exportación de datos
```

---

## Historias de Usuario

| HU           | Título                            | Estado        |
| ------------ | --------------------------------- | ------------- |
| HU-SALES-001 | Gestión del catálogo de productos | ✅ Completado |
| HU-SALES-002 | Carrito de compras y cotizaciones | ✅ Completado |
| HU-SALES-003 | Creación de órdenes de venta      | ✅ Completado |
| HU-SALES-004 | Procesamiento de pagos            | ✅ Completado |
| HU-SALES-005 | Generación de facturas            | ✅ Completado |
| HU-SALES-006 | Reportes y análisis de ventas     | ✅ Completado |

---

## Reglas de Negocio Comunes

### Multi-tenancy

- Todos los datos de ventas están scoped al tenant (empresa)
- Las consultas deben filtrar por `enterprise_id`
- La facturación sigue regulaciones locales (Colombia)

### Seguridad y Permisos

- Solo usuarios con rol ADMIN o MANAGER pueden gestionar catálogo
- Los cajeros pueden procesar ventas y generar facturas
- Los reportes sensibles requieren permisos de gerente

### Inventario y Consistencia

- Las actualizaciones de inventario son atómicas
- No se permite venta de productos sin stock suficiente
- El sistema previene el overselling

### Facturación

- Numeración secuencial por sucursal y empresa
- Las facturas son inmutables después de emisión
- Solo soft delete para eliminación de facturas

---

## Implementación

### Entidades Principales

**Producto (Product)**

- SKU único por empresa
- Precios de costo y venta
- Impuestos y descuentos
- Control de inventario mínimo

**Carrito (Cart)**

- Temporal hasta conversión a orden
- Items con cálculo de impuestos
- Descuentos a nivel de item o carrito

**Orden de Venta (SalesOrder)**

- Número de orden único
- Estado: PENDING_PAYMENT, PAID, CANCELLED, COMPLETED
- Vinculada a factura generada

**Pago (Payment)**

- Múltiples métodos: CASH, CARD, TRANSFER, CREDIT
- Registro detallado de transacciones
- Historial de caja registradora

**Factura (Invoice)**

- Numeración con prefijo configurable
- Vinculación a orden de venta
- Documentación fiscal completa

### Endpoints API Principales

**Productos**

- `POST /products` - Crear producto
- `GET /products` - Listar productos
- `GET /products/{id}` - Obtener producto
- `PUT /products/{id}` - Actualizar producto
- `DELETE /products/{id}` - Eliminar producto (soft delete)

**Categorías**

- `POST /categories` - Crear categoría
- `GET /categories` - Listar categorías
- `GET /categories/{id}` - Obtener categoría
- `PUT /categories/{id}` - Actualizar categoría

**Marcas**

- `POST /brands` - Crear marca
- `GET /brands` - Listar marcas
- `GET /brands/{id}` - Obtener marca
- `PUT /brands/{id}` - Actualizar marca

**Carrito**

- `POST /carts` - Crear carrito
- `GET /carts/{id}` - Obtener carrito
- `GET /carts/code/{code}` - Obtener carrito por código
- `POST /carts/{id}/items` - Añadir item
- `PUT /carts/{id}/items/{itemId}` - Actualizar item
- `DELETE /carts/{id}/items/{itemId}` - Eliminar item
- `POST /carts/{id}/convert` - Convertir a venta
- `PUT /carts/{id}/customer` - Asignar cliente
- `POST /carts/{id}/discount` - Aplicar descuento

**Órdenes de Venta**

- `POST /sales-orders` - Crear orden de venta
- `GET /sales-orders` - Listar órdenes
- `GET /sales-orders/{id}` - Obtener orden
- `PUT /sales-orders/{id}/status` - Actualizar estado

**Pagos**

- `POST /payments` - Procesar pago
- `GET /payments` - Listar pagos
- `GET /payments/{id}` - Obtener pago
- `GET /payments/order/{orderId}` - Pagos de orden

**Caja Registradora**

- `POST /cash-drawers` - Abrir caja
- `GET /cash-drawers/active` - Obtener caja activa
- `POST /cash-drawers/{id}/close` - Cerrar caja

**Facturas**

- `POST /invoices` - Crear factura desde orden
- `GET /invoices` - Listar facturas
- `GET /invoices/{id}` - Obtener factura
- `POST /invoices/{id}/cancel` - Cancelar factura

**Prefijos de Factura**

- `POST /invoice-prefixes` - Crear prefijo
- `GET /invoice-prefixes` - Listar prefijos

**Reportes**

- `GET /reports/sales-summary` - Resumen de ventas
- `GET /reports/product-sales` - Ventas por producto
- `GET /reports/payment-methods` - Desglose por método de pago
- `GET /reports/daily-sales` - Ventas diarias
- `GET /reports/top-customers` - Mejores clientes

---

## Dependencias

### Entidades Existentes

- **Empresa**: `enterprise_id` para scope multi-tenant
- **Usuario**: `user_id` para cajeros y gerentes
- **Cliente**: `customer_id` opcional para ventas con clientes registrados
- **Sucursal**: `branch_id` para múltiples ubicaciones

### Migraciones Creadas

1. `000002_products.up.sql` - Tablas de productos, categorías, marcas
2. `000003_cart.up.sql` - Tablas de carrito y items
3. `000004_sales_orders.up.sql` - Tablas de órdenes de venta
4. `000005_payments.up.sql` - Tablas de pagos y caja registradora
5. `000006_invoices.up.sql` - Tablas de facturas y prefijos

---

## Resumen

- **Total de HU**: 6
- **Completadas**: 6
- **Pendientes**: 0
- **Módulos implementados**: 6 (products, cart, sales, payments, invoices, reports)
