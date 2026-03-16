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

| HU | Título | Estado |
|----|--------|--------|
| HU-SALES-001 | Gestión del catálogo de productos | ⏳ Pendiente |
| HU-SALES-002 | Carrito de compras y cotizaciones | ⏳ Pendiente |
| HU-SALES-003 | Creación de órdenes de venta | ⏳ Pendiente |
| HU-SALES-004 | Procesamiento de pagos | ⏳ Pendiente |
| HU-SALES-005 | Generación de facturas | ⏳ Pendiente |
| HU-SALES-006 | Reportes y análisis de ventas | ⏳ Pendiente |

---

## Reglas de Negocio Comunes

### Multi-tenancy
- Todos los datos de ventas están scoped al tenant (empresa)
- Las consultas deben filtrar por `empresa_id`
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
- `POST /api/products` - Crear producto
- `GET /api/products` - Listar productos
- `GET /api/products/{id}` - Obtener producto
- `PUT /api/products/{id}` - Actualizar producto
- `DELETE /api/products/{id}` - Eliminar producto (soft delete)

**Carrito**
- `POST /api/carts` - Crear carrito
- `POST /api/carts/{id}/items` - Añadir item
- `PUT /api/carts/{id}/items/{itemId}` - Actualizar item
- `POST /api/carts/{id}/convert` - Convertir a venta

**Órdenes de Venta**
- `POST /api/sales-orders` - Crear desde carrito
- `GET /api/sales-orders/{id}` - Obtener orden
- `PUT /api/sales-orders/{id}/status` - Actualizar estado

**Pagos**
- `POST /api/payments` - Procesar pago
- `GET /api/payments/order/{orderId}` - Pagos de orden

**Facturas**
- `POST /api/invoices` - Generar factura
- `GET /api/invoices/{id}` - Obtener factura
- `GET /api/invoices/{id}/pdf` - Generar PDF

**Reportes**
- `GET /api/reports/sales/daily` - Ventas diarias
- `POST /api/reports/sales/product` - Ventas por producto
- `POST /api/reports/sales/employee` - Ventas por empleado

---

## Dependencias

### Entidades Existentes
- **Empresa**: `empresa_id` para scope multi-tenant
- **Usuario**: `user_id` para cajeros y gerentes
- **Cliente**: `customer_id` opcional para ventas con clientes registrados
- **Sucursal**: `branch_id` para múltiples ubicaciones

### Migraciones Requeridas
1. Crear tablas de productos, categorías, marcas
2. Crear tablas de carrito y items
3. Crear tablas de órdenes de venta
4. Crear tablas de pagos y caja registradora
5. Crear tablas de facturas y prefijos

---

## Resumen

- **Total de HU**: 6
- **Completadas**: 0
- **Pendientes**: 6
- **Módulos a implementar**: 6 (products, cart, sales, payments, invoices, reports)
