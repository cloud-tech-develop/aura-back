# GUIA DE IMPLEMENTACION - AURA POS BACKEND

## Alcance del Proyecto

AURA es un sistema POS (Point of Sale) multi-tenant construido con Go que permite gestionar multiples empresas desde una sola instalacion. Cada empresa opera en su propio esquema de PostgreSQL, garantizando aislamiento total de datos.

### Stack Tecnologico

| Componente | Tecnologia |
|------------|------------|
| Lenguaje | Go 1.26.1 |
| Framework HTTP | Gin |
| Base de datos | PostgreSQL |
| Driver DB | lib/pq |
| Migraciones | golang-migrate/v4 |
| Autenticacion | JWT con validacion de IP |
| Testing | testify + go-sqlmock |

### Arquitectura

```
┌─────────────────────────────────────────────────────────┐
│                    CLIENTES (HTTP)                       │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│              GIN ROUTER + MIDDLEWARES                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐ │
│  │ Auth (JWT)  │→ │   Tenant    │→ │  Rate Limiter   │ │
│  └─────────────┘  └─────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                 MODULOS DE NEGOCIO                       │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌───────────┐  │
│  │Enterprise│ │  Users   │ │ Products │ │   Sales   │  │
│  └──────────┘ └──────────┘ └──────────┘ └───────────┘  │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌───────────┐  │
│  │   Cart   │ │ Payments │ │ Invoices │ │Inventory  │  │
│  └──────────┘ └──────────┘ └──────────┘ └───────────┘  │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌───────────┐  │
│  │  Cash    │ │Purchases │ │Shrinkage │ │Transfers  │  │
│  └──────────┘ └──────────┘ └──────────┘ └───────────┘  │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐               │
│  │Commissions│ │ Payroll  │ │ Reports  │               │
│  └──────────┘ └──────────┘ └──────────┘               │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                   POSTGRESQL                             │
│  ┌────────────────────────────────────────────────────┐ │
│  │ public schema (shared)                             │ │
│  │ enterprises, users, roles, user_roles, plans       │ │
│  └────────────────────────────────────────────────────┘ │
│  ┌────────────┐ ┌────────────┐ ┌────────────────────┐  │
│  │ empresa_a  │ │ empresa_b  │ │ empresa_c          │  │
│  │ (tenant)   │ │ (tenant)   │ │ (tenant)           │  │
│  └────────────┘ └────────────┘ └────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

---

## Modulos del Sistema

### 1. Enterprise (Gestion de Empresas)

**Proposito**: Registro y administracion de empresas que usan el sistema.

**Entidades principales**:
- `Enterprise` - Empresa con datos fiscales, configuracion y estado
- `Plan` - Plan de suscripcion con limites (usuarios, empresas)

**Funcionalidades**:
- Crear empresa (genera esquema, usuario admin, migraciones)
- Listar, actualizar y eliminar empresas
- Gestion de estados: ACTIVE, INACTIVE, SUSPENDED, DEBT
- Validacion de slug unico (^[a-z0-9_]+$, 3-50 caracteres)
- Validacion de email unico en todo el sistema
- Validacion de subdominio unico
- Control de limites por plan

**Endpoints**:
```
POST   /enterprises              [Publico] Crear empresa
GET    /enterprises              Listar empresas
GET    /enterprises/:slug        Obtener empresa por slug
PUT    /enterprises/:slug        Actualizar empresa
PATCH  /enterprises/:slug/status Cambiar estado
```

---

### 2. Users (Gestion de Usuarios y Roles)

**Proposito**: Administracion de usuarios del sistema y control de acceso por roles.

**Entidades principales**:
- `User` - Usuario con credenciales y datos de tercero asociado
- `Role` - Rol con nivel de privilegio (0=superadmin, 1=admin, 2=supervisor, 3+=user)

**Funcionalidades**:
- Crear usuario (inserta en public.users + tenant.third_parties)
- Listar usuarios por empresa
- Actualizar datos de usuario
- Activar/desactivar usuarios
- Asignar roles (validacion de nivel de privilegio)
- Listar roles disponibles segun nivel del usuario autenticado

**Jerarquia de roles**:
| Nivel | Rol | Acceso |
|-------|-----|--------|
| 0 | SUPERADMIN | Acceso total al sistema |
| 1 | ADMIN | Administrador de empresa |
| 2 | SUPERVISOR | Permisos extendidos |
| 3+ | USER | Usuario estandar |

**Endpoints**:
```
POST   /users                    Crear usuario
GET    /users                    Listar usuarios
GET    /users/:id                Obtener usuario
PUT    /users/:id                Actualizar usuario
PATCH  /users/:id/status         Activar/desactivar
PATCH  /users/:id/roles          Asignar roles
GET    /roles                    Listar roles disponibles
```

---

### 3. Third Parties (Terceros)

**Proposito**: Registro de clientes, proveedores y empleados.

**Entidades principales**:
- `ThirdParty` - Tercero con tipo (cliente, proveedor, empleado)

**Funcionalidades**:
- CRUD completo de terceros
- Filtrado por tipo (cliente, proveedor, empleado)
- Busqueda por documento de identidad
- Tipos de documento: CC, CE, NIT, PASSPORT, RUT
- Responsabilidades fiscales: RESPONSIBLE, NOT-RESPONSIBLE

**Endpoints**:
```
POST   /third-parties                           Crear tercero
GET    /third-parties                           Listar terceros
GET    /third-parties/:id                       Obtener tercero
GET    /third-parties/document/:documentNumber  Buscar por documento
PUT    /third-parties/:id                       Actualizar tercero
DELETE /third-parties/:id                       Eliminar tercero
```

---

### 4. Products (Productos)

**Proposito**: Gestion del catalogo de productos con categorias y marcas.

**Entidades principales**:
- `Product` - Producto con SKU, precios, impuestos y stock minimo
- `Category` - Categoria con jerarquia (parent)
- `Brand` - Marca de productos

**Funcionalidades**:
- CRUD de productos, categorias y marcas
- Gestion de precios (costo, venta)
- Configuracion de impuestos por producto
- Stock minimo para alertas
- Filtrado por categoria y marca

**Endpoints**:
```
POST   /products                Crear producto
GET    /products                Listar productos
GET    /products/:id            Obtener producto
PUT    /products/:id            Actualizar producto
DELETE /products/:id            Eliminar producto

POST   /categories              Crear categoria
GET    /categories              Listar categorias
GET    /categories/:id          Obtener categoria
PUT    /categories/:id          Actualizar categoria

POST   /brands                  Crear marca
GET    /brands                  Listar marcas
GET    /brands/:id              Obtener marca
PUT    /brands/:id              Actualizar marca
```

---

### 5. Cart (Carrito de Compras)

**Proposito**: Gestion del carrito de ventas y cotizaciones.

**Entidades principales**:
- `Cart` - Carrito con tipo (SALE/QUOTATION), totales y estado
- `CartItem` - Item con cantidad, precio, descuento e impuestos

**Estados del carrito**: ACTIVE, SAVED, CONVERTED, EXPIRED, CANCELLED

**Funcionalidades**:
- Crear carrito de venta o cotizacion
- Agregar, actualizar y eliminar items
- Aplicar descuentos (porcentaje o fijo) a carrito o items
- Asignar cliente al carrito
- Convertir carrito a orden de venta
- Convertir carrito a cotizacion con fecha de vigencia
- Calculo automatico de totales (subtotal, descuento, impuestos, total)

**Endpoints**:
```
POST   /carts                              Crear carrito
GET    /carts                              Listar carritos
GET    /carts/:id                          Obtener carrito
GET    /carts/code/:code                   Buscar por codigo
DELETE /carts/:id                          Eliminar carrito
POST   /carts/:id/items                    Agregar item
PUT    /carts/:id/items/:itemId            Actualizar item
DELETE /carts/:id/items/:itemId            Eliminar item
POST   /carts/:id/items/:itemId/discount   Descuento a item
POST   /carts/:id/convert                  Convertir a venta
POST   /carts/:id/quotation                Convertir a cotizacion
PUT    /carts/:id/customer                 Asignar cliente
POST   /carts/:id/discount                 Descuento al carrito
```

---

### 6. Sales (Ventas)

**Proposito**: Gestion de ordenes de venta.

**Entidades principales**:
- `SalesOrder` - Orden de venta con items, totales y estado
- `SalesOrderItem` - Item de la orden

**Estados de orden**: PENDING_PAYMENT, PAID, CANCELLED, COMPLETED

**Funcionalidades**:
- Crear orden desde carrito convertido
- Listar y consultar ordenes
- Cambiar estado de orden
- Cancelar y completar ordenes

**Endpoints**:
```
POST   /sales-orders                Crear orden
GET    /sales-orders                Listar ordenes
GET    /sales-orders/:id            Obtener orden
PUT    /sales-orders/:id/status     Cambiar estado
POST   /sales-orders/:id/cancel     Cancelar orden
POST   /sales-orders/:id/complete   Completar orden
```

---

### 7. Payments (Pagos)

**Proposito**: Procesamiento de pagos y gestion de cajas registradoras.

**Entidades principales**:
- `Payment` - Transaccion de pago con metodo y monto
- `PaymentTransaction` - Log de transacciones (cargo, reembolso, chargeback)
- `CashDrawer` - Caja registradora por sucursal
- `CashMovement` - Movimiento de efectivo (entrada/salida)

**Metodos de pago**: CASH, DEBIT_CARD, CREDIT_CARD, BANK_TRANSFER, CREDIT, VOUCHER, CHECK

**Funcionalidades**:
- Procesar pagos individuales y multiples
- Calcular cambio automaticamente
- Cancelar pagos con motivo
- Abrir/cerrar cajas registradoras
- Registrar entradas y salidas de efectivo
- Listar pagos por referencia (orden, factura)

**Endpoints**:
```
POST   /payments                           Procesar pago
POST   /payments/batch                     Pago multiple
GET    /payments                           Listar pagos
GET    /payments/:id                       Obtener pago
GET    /payments/reference/:type/:id       Pagos por referencia
POST   /payments/:id/cancel                Cancelar pago

POST   /cash-drawers                       Abrir caja
GET    /cash-drawers                       Listar cajas
GET    /cash-drawers/open                  Caja abierta actual
GET    /cash-drawers/:id                   Obtener caja
POST   /cash-drawers/:id/close             Cerrar caja
POST   /cash-drawers/:id/cash-in           Entrada de efectivo
POST   /cash-drawers/:id/cash-out          Salida de efectivo
```

---

### 8. Cash (Caja y Turnos)

**Proposito**: Gestion de turnos de caja con arqueo y conciliacion.

**Entidades principales**:
- `CashDrawer` - Configuracion de caja por sucursal
- `CashShift` - Turno de caja con montos de apertura/cierre
- `CashMovement` - Movimiento durante el turno

**Estados de turno**: OPEN, CLOSED, AUDITED

**Razones de movimiento**: SALE, OPENING, CLOSING, EXPENSE, DROPS, WITHDRAWAL, ADJUSTMENT, REFUND

**Funcionalidades**:
- Configurar caja por sucursal
- Abrir/cerrar turnos
- Registrar movimientos de caja
- Conciliar turno (comparar esperado vs real)
- Generar resumen de turno
- Listar historial de turnos

**Endpoints**:
```
GET    /cash/drawer/:branchID              Obtener caja por sucursal
POST   /cash/drawer                        Configurar caja
POST   /cash/shift/open                    Abrir turno
POST   /cash/shift/:shiftID/close          Cerrar turno
GET    /cash/shift/active                  Turno activo del usuario
GET    /cash/shift/:shiftID                Resumen de turno
GET    /cash/shifts                        Listar turnos
POST   /cash/movement                      Registrar movimiento
POST   /cash/shift/:shiftID/reconcile      Conciliar turno
```

---

### 9. Invoices (Facturacion)

**Proposito**: Generacion y gestion de facturas de venta.

**Entidades principales**:
- `Invoice` - Factura con datos fiscales, impuestos y totales
- `InvoiceItem` - Item de factura
- `InvoicePrefix` - Numeracion de facturas por resolucion
- `InvoiceLog` - Auditoria de cambios

**Tipos de factura**: SALE, CREDIT_NOTE, DEBIT_NOTE

**Estados**: DRAFT, ISSUED, SENT, VIEWED, CANCELLED

**Funcionalidades**:
- Generar factura desde orden de venta
- Crear factura manual
- Emitir factura (cambia estado)
- Cancelar factura con motivo
- Gestion de prefijos de numeracion (resoluciones DIAN)
- Registro de auditoria
- Calculo de impuestos: IVA 19%, IVA 5%, ReteICA, Retefuente

**Endpoints**:
```
POST   /invoices/generate          Generar desde venta
POST   /invoices                   Crear factura
GET    /invoices                   Listar facturas
GET    /invoices/:id               Obtener factura
GET    /invoices/number/:number    Buscar por numero
POST   /invoices/:id/issue         Emitir factura
POST   /invoices/:id/cancel        Cancelar factura
GET    /invoices/:id/logs          Ver auditoria

POST   /invoice-prefixes           Crear prefijo
GET    /invoice-prefixes           Listar prefijos
```

---

### 10. Inventory (Inventario)

**Proposito**: Control de stock por sucursal con kardex de movimientos.

**Entidades principales**:
- `Inventory` - Stock actual por producto y sucursal
- `InventoryMovement` - Movimiento de kardex (entrada/salida/ajuste)
- `MovementReason` - Configuracion de tipos de movimiento

**Tipos de movimiento**: ENTRY, EXIT, ADJUSTMENT

**Razones**: SALE, PURCHASE, SHRINKAGE, TRANSFER_IN, TRANSFER_OUT, ADJUSTMENT, RETURN, INITIAL, DAMAGE, THEFT, EXPIRED

**Funcionalidades**:
- Consultar stock por producto y sucursal
- Listar inventario completo
- Alertas de stock bajo (bajo minimo)
- Ver kardex de un producto
- Registrar movimientos de inventario
- Listar historial de movimientos

**Endpoints**:
```
GET    /inventory                              Listar inventario
GET    /inventory/low-stock                    Productos con stock bajo
GET    /inventory/:productId/:branchId         Stock especifico
GET    /inventory/product/:productId           Stock por sucursal
GET    /inventory/kardex/:productId/:branchId  Kardex del producto
POST   /inventory/movements                    Registrar movimiento
GET    /movements                              Listar movimientos
GET    /movements/:id                          Detalle de movimiento
GET    /movement-reasons                       Razones de movimiento
```

---

### 11. Purchases (Compras)

**Proposito**: Gestion de ordenes de compra, recepcion de mercancia y pagos a proveedores.

**Entidades principales**:
- `PurchaseOrder` - Orden de compra a proveedor
- `PurchaseOrderItem` - Item de la orden
- `Purchase` - Recepcion de mercancia
- `PurchaseItem` - Item recibido
- `PurchasePayment` - Pago a proveedor
- `SupplierSummary` - Resumen de cuenta del proveedor

**Estados**: PENDING, PARTIAL, RECEIVED, CANCELLED, COMPLETED

**Funcionalidades**:
- Crear orden de compra
- Recibir mercancia (genera movimientos de inventario)
- Registrar pagos a proveedores
- Cancelar compras
- Historial de compras por proveedor
- Resumen de cuenta por proveedor

**Endpoints**:
```
POST   /purchases/orders                     Crear orden
GET    /purchases/orders/:id                 Obtener orden
GET    /purchases/orders                     Listar ordenes
POST   /purchases/receive                    Recibir mercancia
GET    /purchases/:id                        Obtener compra
GET    /purchases                            Listar compras
POST   /purchases/:id/cancel                 Cancelar compra
POST   /purchases/payments                   Registrar pago
GET    /purchases/suppliers/:id/summary      Resumen proveedor
```

---

### 12. Shrinkage (Mermas)

**Proposito**: Registro y control de mermas con autorizacion.

**Entidades principales**:
- `ShrinkageReason` - Motivo de merma con umbral de autorizacion
- `Shrinkage` - Registro de merma
- `ShrinkageItem` - Producto afectado
- `ShrinkageReportItem` - Datos para reportes

**Estados**: PENDING, APPROVED, REJECTED, CANCELLED

**Funcionalidades**:
- Registrar merma con productos afectados
- Configurar razones de merma
- Flujo de autorizacion (para mermas sobre umbral)
- Aprobar/rechazar mermas
- Cancelar mermas
- Reporte de mermas por periodo

**Endpoints**:
```
POST   /shrinkage                  Registrar merma
GET    /shrinkage/:id              Detalle de merma
GET    /shrinkage                  Listar mermas
POST   /shrinkage/:id/authorize    Autorizar merma
POST   /shrinkage/:id/cancel       Cancelar merma
POST   /shrinkage/reasons          Crear razon
GET    /shrinkage/reasons          Listar razones
GET    /shrinkage/report           Reporte de mermas
```

---

### 13. Transfers (Traslados)

**Proposito**: Transferencia de inventario entre sucursales.

**Entidades principales**:
- `Transfer` - Traslado entre sucursales
- `TransferItem` - Producto a transferir

**Estados**: PENDING, APPROVED, SHIPPED, PARTIAL, RECEIVED, CANCELLED

**Funcionalidades**:
- Crear solicitud de traslado
- Aprobar traslado
- Marcar como enviado (reduce stock origen)
- Recibir traslado (aumenta stock destino, permite recepcion parcial)
- Cancelar traslado

**Endpoints**:
```
POST   /transfers                Crear traslado
GET    /transfers/:id            Detalle de traslado
GET    /transfers                Listar traslados
POST   /transfers/:id/approve    Aprobar traslado
POST   /transfers/:id/ship       Enviar traslado
POST   /transfers/:id/receive    Recibir traslado
POST   /transfers/:id/cancel     Cancelar traslado
```

---

### 14. Commissions (Comisiones)

**Proposito**: Calculo y liquidacion de comisiones por ventas.

**Entidades principales**:
- `CommissionRule` - Regla de comision (porcentaje venta, margen, monto fijo)
- `Commission` - Comision calculada por venta
- `CommissionSummary` - Resumen por empleado

**Tipos de comision**: PERCENTAGE_SALE, PERCENTAGE_MARGIN, FIXED_AMOUNT

**Estados**: PENDING, SETTLED, CANCELLED

**Funcionalidades**:
- Configurar reglas de comision (por producto, categoria, empleado)
- Calcular comisiones automaticamente al completar venta
- Liquidar comisiones por periodo
- Reporte de comisiones por empleado
- Totales de comisiones pendientes y liquidadas

**Endpoints**:
```
POST   /commissions/rules                  Crear regla
PUT    /commissions/rules/:id              Actualizar regla
DELETE /commissions/rules/:id              Eliminar regla
GET    /commissions/rules                  Listar reglas
POST   /commissions/calculate              Calcular comisiones
GET    /commissions                        Listar comisiones
GET    /commissions/:id                    Detalle de comision
POST   /commissions/settle                 Liquidar comisiones
GET    /commissions/report/summary         Reporte resumen
GET    /commissions/report/totals          Totales
```

---

### 15. Payroll (Nomina)

**Proposito**: Gestion completa de nomina, empleados, deducciones y prestaciones.

**Entidades principales**:
- `Employee` - Empleado con datos laborales y bancarios
- `Salary` - Historial salarial
- `PayrollPeriod` - Periodo de nomina (semanal, quincenal, mensual)
- `Payroll` - Nomina individual por empleado
- `PayrollDetail` - Detalle de conceptos (devengos, deducciones)
- `DeductionType` - Tipo de deduccion (ISR, IMSS, etc.)
- `AdditionType` - Tipo de devengo (bono, horas extra, etc.)
- `EmployeeLoan` - Anticipos y prestamos
- `Overtime` - Horas extra
- `Bonus` - Bonos
- `PayrollPayment` - Pago de nomina
- `LeaveType` - Tipo de permiso
- `EmployeeLeave` - Solicitud de permiso
- `LeaveBalance` - Saldo de permisos

**Estados de periodo**: OPEN, PROCESSING, APPROVED, PAID, CLOSED

**Funcionalidades**:
- Gestion de empleados (alta, baja, modificacion)
- Configuracion salarial con historial
- Crear periodos de nomina
- Calcular nomina automaticamente
- Registrar horas extra, bonos y prestamos
- Flujo de aprobacion de nomina
- Procesar pagos de nomina
- Importar comisiones a nomina
- Gestion de permisos y vacaciones
- Reportes fiscales (ISR, IMSS)
- Reporte de deducciones

**Flujo de nomina**:
```
1. Crear periodo de nomina
2. Registrar horas extra, bonos, prestamos (durante el periodo)
3. Cerrar periodo
4. Calcular nomina (genera registros por empleado)
5. Preview/revision
6. Aprobar nomina
7. Procesar pagos
8. Cerrar periodo definitivamente
```

**Endpoints**:
```
# Empleados
POST   /payroll/employees                          Crear empleado
GET    /payroll/employees/:id                      Obtener empleado
PUT    /payroll/employees/:id                      Actualizar empleado
POST   /payroll/employees/:id/terminate            Dar de baja
GET    /payroll/employees                          Listar empleados
POST   /payroll/employees/:id/salary               Asignar salario

# Periodos
POST   /payroll/periods                            Crear periodo
GET    /payroll/periods/:id                        Obtener periodo
GET    /payroll/periods                            Listar periodos
POST   /payroll/periods/:periodId/calculate         Calcular nomina
POST   /payroll/periods/:periodId/approve           Aprobar nomina
POST   /payroll/periods/:periodId/close             Cerrar periodo
POST   /payroll/periods/:periodId/reopen            Reabrir periodo
POST   /payroll/periods/:periodId/reject            Rechazar nomina

# Registros de nomina
GET    /payroll/periods/:periodId/payrolls          Listar nominas
GET    /payroll/periods/:periodId/preview           Preview de nomina
GET    /payroll/periods/:periodId/summary           Resumen del periodo
GET    /payroll/payslips/:payrollId                 Recibo de nomina

# Conceptos
POST   /payroll/overtime                           Registrar horas extra
POST   /payroll/bonuses                            Registrar bono
POST   /payroll/loans                              Registrar prestamo

# Pagos
POST   /payroll/payments                           Procesar pago
GET    /payroll/payments                           Listar pagos

# Comisiones
POST   /payroll/periods/:periodId/import-commissions Importar comisiones

# Reportes
GET    /payroll/periods/:periodId/reports/tax       Reporte fiscal
GET    /payroll/periods/:periodId/reports/deductions Reporte deducciones
GET    /payroll/employees/:employeeId/earnings      Historial de ingresos

# Tipos
GET    /payroll/deduction-types                    Tipos de deduccion
GET    /payroll/addition-types                     Tipos de devengo

# Permisos
POST   /payroll/leave-types                        Crear tipo de permiso
GET    /payroll/leave-types                        Listar tipos
POST   /payroll/leaves                             Solicitar permiso
POST   /payroll/leaves/:leaveId/approve            Aprobar/rechazar
GET    /payroll/employees/:employeeId/leave-balance Saldo de permisos
```

---

### 16. Reports (Reportes)

**Proposito**: Generacion de reportes y exportacion a PDF/Excel.

**Reportes disponibles**:
- Ventas por periodo
- Ventas por producto
- Ventas por empleado
- Metodos de pago
- Ventas diarias
- Top clientes
- Estado de inventario
- Movimientos de inventario

**Funcionalidades**:
- Reportes con filtros (fecha, sucursal, etc.)
- Paginacion de resultados
- Exportacion a PDF
- Exportacion a Excel

**Endpoints**:
```
GET    /reports/sales-summary               Resumen de ventas
GET    /reports/product-sales               Ventas por producto
GET    /reports/payment-methods             Metodos de pago
GET    /reports/daily-sales                 Ventas diarias
GET    /reports/top-customers               Top clientes
GET    /reports/sales                       Ventas por periodo
GET    /reports/sales/products              Ventas detalladas por producto
GET    /reports/sales/employees             Ventas por empleado
GET    /reports/inventory                   Estado de inventario
GET    /reports/inventory/movements         Movimientos de inventario
POST   /reports/:type/export/pdf            Exportar a PDF
POST   /reports/:type/export/excel          Exportar a Excel
```

---

## Flujo de Negocio Principal

### Venta Completa (Punto de Venta)

```
1. Usuario inicia sesion → JWT con tenant, roles, IP
2. Abre turno de caja → Cash.OpenShift
3. Crea carrito de venta → Cart.CreateCart
4. Agrega productos al carrito → Cart.AddItem
5. (Opcional) Asigna cliente → Cart.SetCustomer
6. (Opcional) Aplica descuentos → Cart.ApplyDiscount
7. Convierte carrito a venta → Cart.ConvertToSale
   └── Genera SalesOrder automaticamente
8. Procesa pago(s) → Payment.ProcessPayment
   └── Puede ser pago multiple (efectivo + tarjeta)
9. Genera factura → Invoice.GenerateInvoiceFromSale
10. Emite factura → Invoice.IssueInvoice
11. Actualiza inventario → Inventory.UpdateStock (automatico)
12. Cierra turno de caja → Cash.CloseShift
```

### Compra a Proveedor

```
1. Crear orden de compra → Purchase.CreatePurchaseOrder
2. Recibir mercancia → Purchase.ReceiveGoods
   └── Actualiza inventario automaticamente
3. Registrar pago → Purchase.RecordPayment
4. Consultar estado de cuenta → Purchase.GetSupplierSummary
```

### Traslado entre Sucursales

```
1. Solicitar traslado → Transfer.CreateTransfer
2. Aprobar traslado → Transfer.ApproveTransfer
3. Enviar mercancia → Transfer.ShipTransfer
   └── Reduce stock en sucursal origen
4. Recibir mercancia → Transfer.ReceiveTransfer
   └── Aumenta stock en sucursal destino
```

### Liquidacion de Nomina

```
1. Crear periodo de nomina → Payroll.CreatePeriod
2. Registrar horas extra → Payroll.RegisterOvertime
3. Registrar bonos → Payroll.RegisterBonus
4. Importar comisiones → Payroll.ImportCommissions
5. Cerrar periodo → Payroll.ClosePeriod
6. Calcular nomina → Payroll.CalculatePayroll
7. Aprobar nomina → Payroll.ApprovePayroll
8. Procesar pagos → Payroll.ProcessPayment
9. Cerrar periodo → Payroll.ClosePeriod (definitivo)
```

---

## Tablas por Esquema

### Esquema PUBLIC (compartido entre todos los tenants)

| Tabla | Descripcion |
|-------|-------------|
| enterprises | Registro de empresas |
| tenants | Registro de tenants |
| users | Usuarios del sistema |
| roles | Roles y permisos |
| user_roles | Relacion usuario-rol |
| plans | Planes de suscripcion |

### Esquema TENANT (por empresa)

| Tabla | Descripcion |
|-------|-------------|
| third_parties | Clientes, proveedores, empleados |
| products | Catalogo de productos |
| categories | Categorias de productos |
| brands | Marcas de productos |
| carts | Carritos de compra |
| cart_items | Items del carrito |
| sales_orders | Ordenes de venta |
| sales_order_items | Items de la orden |
| payments | Transacciones de pago |
| payment_transactions | Log de transacciones |
| cash_drawers | Cajas registradoras |
| cash_movements | Movimientos de caja |
| cash_shifts | Turnos de caja |
| invoices | Facturas |
| invoice_items | Items de factura |
| invoice_prefixes | Prefijos de numeracion |
| invoice_logs | Auditoria de facturas |
| inventory | Stock por sucursal |
| inventory_movements | Kardex de movimientos |
| purchase_orders | Ordenes de compra |
| purchase_order_items | Items de orden de compra |
| purchases | Recepciones de mercancia |
| purchase_items | Items recibidos |
| purchase_payments | Pagos a proveedores |
| shrinkage_reasons | Razones de merma |
| shrinkages | Registros de merma |
| shrinkage_items | Items de merma |
| transfers | Traslados entre sucursales |
| transfer_items | Items de traslado |
| commission_rules | Reglas de comision |
| commissions | Comisiones calculadas |
| employees | Empleados |
| salaries | Historial salarial |
| payroll_periods | Periodos de nomina |
| payrolls | Nominas por empleado |
| payroll_details | Detalle de conceptos |
| deduction_types | Tipos de deduccion |
| addition_types | Tipos de devengo |
| employee_loans | Prestamos y anticipos |
| overtime | Horas extra |
| bonuses | Bonos |
| payroll_payments | Pagos de nomina |
| leave_types | Tipos de permiso |
| employee_leaves | Permisos y vacaciones |
