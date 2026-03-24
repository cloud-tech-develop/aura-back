# EPICA-005: Shopping Cart / Quotations

## 📌 General Information
- ID: EPICA-005
- State: Completed
- Priority: High
- Start Date: 2026-03-23
- Target Date: 2026-05-15
- Owner: Aura POS Backend Team
- Percentage: 100%

---

## 🎯 Business Objective

Provide shopping cart and quotation management functionality for Aura POS. The system must allow creating carts, adding/removing items, applying discounts, calculating taxes, and converting carts to sales orders or quotations.

**What problem does it solve?**
- Enables pre-sale quote preparation
- Supports multiple pricing rules per customer
- Calculates totals with taxes and discounts
- Prepares items before finalizing sale

**What value does it generate?**
- Faster checkout process
- Accurate pricing with customer-specific rules
- Sales opportunity through quotations
- Flexible discount management

---

## 👥 Stakeholders

- End User: Cashiers, sales representatives, store managers
- Technical Team: Backend developers
- Product: Product Manager Aura POS

---

## 🧠 Functional Description

The Cart module handles:

1. **Cart Creation**: New cart for a customer/session
2. **Item Management**: Add, update, remove items
3. **Pricing**: Product prices, customer prices, volume pricing
4. **Discounts**: Item-level, cart-level, customer-specific
5. **Taxes**: IVA, RETEICA calculations
6. **Conversion**: Cart to sales order or quotation

---

## 📦 Scope

### Included:
- Cart creation and management
- Add/update/remove items
- Quantity adjustments
- Discount application (item and cart level)
- Tax calculation (IVA, RETEICA)
- Customer-specific pricing
- Volume discounts
- Cart expiration
- Quotation validity period
- Cart to order conversion

### Not Included:
- Shopping cart persistence across sessions (after logout)
- Wishlist functionality
- Price comparison tools
- Advanced promotion engine

---

## 🧩 User Stories

| HU | Title | State |
|----|-------|-------|
| HU-CART-001 | Create Cart | ✅ Completed |
| HU-CART-002 | Add Item to Cart | ✅ Completed |
| HU-CART-003 | Update Cart Item | ✅ Completed |
| HU-CART-004 | Remove Item from Cart | ✅ Completed |
| HU-CART-005 | Apply Discounts | ✅ Completed |
| HU-CART-006 | Calculate Totals and Taxes | ✅ Completed |
| HU-CART-007 | Convert Cart to Quotation | ✅ Completed |
| HU-CART-008 | Convert Cart to Sale Order | ✅ Completed |

---

## 🐞 Associated Bugs

None identified

---

## 🔐 Global Business Rules

- Cart belongs to a specific branch and user
- Cart expires after configurable period (default 24 hours)
- Discount priority: customer price > volume discount > item discount > cart discount
- Tax calculation follows Colombian regulations
- Cart must be validated before conversion
- All prices are in COP
- Quantities must be positive integers
- Items reference product variants when applicable

---

## 🧱 Related Architecture

**Backend:** Go 1.26.1 with Gin framework
**Database:** PostgreSQL with schema-per-tenant
**Authentication:** JWT with tenant context

### Database Schema (Tenant Schema)

**Table: cart**
```sql
CREATE TABLE cart (
    id BIGSERIAL PRIMARY KEY,
    cart_code VARCHAR(50) NOT NULL,
    cart_type VARCHAR(20) NOT NULL DEFAULT 'SALE' CHECK (cart_type IN ('SALE', 'QUOTATION')),
    customer_id BIGINT REFERENCES third_party(id),
    user_id BIGINT NOT NULL REFERENCES public.users(id),
    branch_id BIGINT NOT NULL REFERENCES public.branches(id),
    enterprise_id BIGINT NOT NULL,
    subtotal DECIMAL(12,2) NOT NULL DEFAULT 0,
    discount DECIMAL(12,2) NOT NULL DEFAULT 0,
    tax_total DECIMAL(12,2) NOT NULL DEFAULT 0,
    total DECIMAL(12,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'SAVED', 'CONVERTED', 'EXPIRED', 'CANCELLED')),
    notes TEXT,
    valid_until TIMESTAMP,
    converted_at TIMESTAMP,
    reference_id BIGINT,
    reference_type VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

**Table: cart_item**
```sql
CREATE TABLE cart_item (
    id BIGSERIAL PRIMARY KEY,
    cart_id BIGINT NOT NULL REFERENCES cart(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES product(id),
    product_variant_id BIGINT REFERENCES product_variant(id),
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(12,2) NOT NULL,
    discount_type VARCHAR(20) CHECK (discount_type IN ('PERCENTAGE', 'FIXED')),
    discount_value DECIMAL(12,2) DEFAULT 0,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    tax_rate DECIMAL(5,2) NOT NULL DEFAULT 19.00,
    tax_amount DECIMAL(12,2) DEFAULT 0,
    line_total DECIMAL(12,2) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);
```

---

## 📊 Success Metrics

- Cart conversion rate > 70%
- Pricing accuracy 100%
- Tax calculation accuracy 100%
- Cart processing time < 200ms

---

## 🚧 Risks

- Concurrent cart updates by multiple users
- Price calculation complexity with multiple rules
- Cart expiration handling
- Large cart performance

---

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/carts` | Create cart or quotation |
| GET | `/carts` | List carts with filters |
| GET | `/carts/:id` | Get cart by ID |
| GET | `/carts/code/:code` | Get cart by code |
| DELETE | `/carts/:id` | Delete cart (soft delete) |
| POST | `/carts/:id/items` | Add item to cart |
| PUT | `/carts/:id/items/:itemId` | Update cart item |
| DELETE | `/carts/:id/items/:itemId` | Remove item from cart |
| POST | `/carts/:id/items/:itemId/discount` | Apply item discount |
| POST | `/carts/:id/convert` | Convert cart to sale |
| POST | `/carts/:id/quotation` | Convert cart to quotation |
| PUT | `/carts/:id/customer` | Set customer for cart |
| POST | `/carts/:id/discount` | Apply cart discount |

### Query Parameters for List
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 20)
- `type` - Filter by type: sale, quotation
- `status` - Filter by status
- `branch_id` - Filter by branch

---

## 📁 Module Structure

```
modules/cart/
├── domain.go     # Entity, Repository & Service interfaces
├── service.go    # Repository & Service implementation
├── handler.go    # HTTP handlers
└── routes.go     # Route registration
```

### Migration Updated
- `tenant/migrations/tenant/000003_cart.up.sql` - Enhanced with cart_type, quotation fields

---

## Resumen

- **Total de HU**: 8
- **Completadas**: 8
- **Pendientes**: 0
- **Módulo implementado**: cart (with quotation support)
