# EPICA-010: Purchases

## 📌 General Information
- ID: EPICA-010
- State: Completed
- Priority: Medium
- Start Date: 2026-03-23
- Target Date: 2026-07-30
- Owner: Aura POS Backend Team
- Percentage: 100%

---

## 🎯 Business Objective

Manage the purchasing process from suppliers, including purchase order creation, goods receipt, purchase payments, and accounts payable tracking. This module complements the sales process by handling inventory replenishment.

**What problem does it solve?**
- Manual purchase tracking
- No supplier order history
- Purchase payment tracking
- Inventory replenishment visibility

**What value does it generate?**
- Complete purchase records
- Supplier relationship management
- Accounts payable automation
- Better inventory management

---

## 👥 Stakeholders

- End User: Store managers, inventory clerks, accountants
- Technical Team: Backend developers
- Product: Product Manager Aura POS

---

## 🧠 Functional Description

The Purchases module handles:

1. **Purchase Orders**: Create and manage orders to suppliers
2. **Goods Receipt**: Accept delivered goods
3. **Purchase Payments**: Record payments to suppliers
4. **Accounts Payable**: Track supplier debts
5. **Purchase History**: Complete supplier transaction history

---

## 📦 Scope

### Included:
- Purchase order creation
- Purchase order status management
- Goods receipt processing
- Purchase payment recording
- Partial receipts
- Partial payments
- Purchase cancellation
- Supplier purchase history

### Not Included:
- Advanced purchase planning
- Supplier portal
- Automatic reorder
- Purchase approval workflows

---

## 🧩 User Stories

| HU | Title | State |
|----|-------|-------|
| HU-PUR-001 | Create Purchase Order | ✅ Implemented |
| HU-PUR-002 | Receive Goods | ✅ Implemented |
| HU-PUR-003 | Record Purchase Payment | ✅ Implemented |
| HU-PUR-004 | Cancel Purchase | ✅ Implemented |
| HU-PUR-005 | View Purchase History | ✅ Implemented |
| HU-PUR-006 | Supplier Account Summary | ✅ Implemented |

---

## 🐞 Associated Bugs

None identified

---

## 🔐 Global Business Rules

- Purchase orders can be partially received
- Inventory is updated upon goods receipt
- Purchases to credit generate accounts payable
- Purchases can have multiple payments
- Cancellation creates inverse inventory movement
- Suppliers are third parties with type SUPPLIER

---

## 🧱 Related Architecture

**Backend:** Go 1.26.1 with Gin framework
**Database:** PostgreSQL with schema-per-tenant
**Authentication:** JWT with tenant context

### Database Schema (Tenant Schema)

**Table: purchase_order**
```sql
CREATE TABLE purchase_order (
    id BIGSERIAL PRIMARY KEY,
    order_number VARCHAR(50) NOT NULL,
    supplier_id BIGINT NOT NULL REFERENCES third_party(id),
    branch_id BIGINT NOT NULL REFERENCES branch(id),
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    order_date DATE NOT NULL,
    expected_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'PARTIAL', 'RECEIVED', 'CANCELLED')),
    subtotal DECIMAL(12,2) NOT NULL DEFAULT 0,
    discount_total DECIMAL(12,2) NOT NULL DEFAULT 0,
    tax_total DECIMAL(12,2) NOT NULL DEFAULT 0,
    total DECIMAL(12,2) NOT NULL DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_purchase_order_number ON purchase_order(order_number);
CREATE INDEX idx_purchase_order_supplier ON purchase_order(supplier_id);
CREATE INDEX idx_purchase_order_status ON purchase_order(status);
CREATE INDEX idx_purchase_order_date ON purchase_order(order_date);
```

**Table: purchase_order_item**
```sql
CREATE TABLE purchase_order_item (
    id BIGSERIAL PRIMARY KEY,
    purchase_order_id BIGINT NOT NULL REFERENCES purchase_order(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES product(id),
    quantity DECIMAL(10,2) NOT NULL,
    received_quantity DECIMAL(10,2) NOT NULL DEFAULT 0,
    unit_cost DECIMAL(12,2) NOT NULL,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    tax_rate DECIMAL(5,2) NOT NULL DEFAULT 0,
    line_total DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_po_item_order ON purchase_order_item(purchase_order_id);
```

**Table: purchase**
```sql
CREATE TABLE purchase (
    id BIGSERIAL PRIMARY KEY,
    purchase_number VARCHAR(50) NOT NULL,
    purchase_order_id BIGINT REFERENCES purchase_order(id),
    supplier_id BIGINT NOT NULL REFERENCES third_party(id),
    branch_id BIGINT NOT NULL REFERENCES branch(id),
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    purchase_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'COMPLETED' CHECK (status IN ('COMPLETED', 'PARTIAL', 'CANCELLED')),
    subtotal DECIMAL(12,2) NOT NULL DEFAULT 0,
    discount_total DECIMAL(12,2) NOT NULL DEFAULT 0,
    tax_total DECIMAL(12,2) NOT NULL DEFAULT 0,
    total DECIMAL(12,2) NOT NULL DEFAULT 0,
    paid_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    pending_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_purchase_number ON purchase(purchase_number);
CREATE INDEX idx_purchase_supplier ON purchase(supplier_id);
CREATE INDEX idx_purchase_status ON purchase(status);
CREATE INDEX idx_purchase_date ON purchase(purchase_date);
```

---

## 📊 Success Metrics

- Purchase processing time < 2 seconds
- 100% inventory update accuracy
- Payment tracking accuracy 100%
- Zero negative inventory from purchases

---

## 🚧 Risks

- Partial receipt complexity
- Inventory updates with batch products
- Supplier credit limit validation
- Purchase cancellation with payments

---

## 📡 API Endpoints

### Purchase Orders
```
POST   /purchases/orders                → Crear orden de compra (HU-PUR-001)
GET    /purchases/orders/:id            → Obtener orden de compra
GET    /purchases/orders                → Listar órdenes de compra
```

### Goods Receipt
```
POST   /purchases/receive               → Recibir mercancía (HU-PUR-002)
```

### Purchases
```
GET    /purchases/:id                   → Obtener compra
GET    /purchases                       → Listar historial de compras (HU-PUR-005)
POST   /purchases/:id/cancel            → Cancelar compra (HU-PUR-004)
```

### Payments
```
POST   /purchases/payments              → Registrar pago (HU-PUR-003)
```

### Supplier Summary
```
GET    /purchases/suppliers/:id/summary → Resumen de cuenta proveedor (HU-PUR-006)
```

### Query Parameters
- `page` - Número de página
- `limit` - Items por página
- `status` - Filtrar por estado: PENDING, PARTIAL, RECEIVED, CANCELLED
- `supplier_id` - Filtrar por proveedor
- `start_date` - Fecha inicio
- `end_date` - Fecha fin

---

## 📁 Module Structure

```
modules/purchases/
├── domain.go     # Entity, Repository & Service interfaces
├── repository.go # Repository implementation
├── service.go    # Service implementation
├── handler.go    # HTTP handlers
└── routes.go     # Route registration
```

---

## Resumen

- **Total de HU**: 6
- **Completadas**: 6
- **Pendientes**: 0
- **Módulo implementado**: purchases
