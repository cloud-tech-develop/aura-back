# HU-SALES-003 - Sale/Order Creation

## 📌 General Information
- ID: HU-SALES-003
- Epic: EPIC-SALES-001
- Priority: High
- State: Backlog
- Progress: 0%
- Author: QA Engineer Aura POS
- Date: 2026-03-15

---

## 👤 User Story

**As a** cashier or salesperson
**I want to** create a sales order from a shopping cart
**So that** I can process customer purchases and update inventory

---

## 🧠 Functional Description

The system must convert a shopping cart into a permanent sales order. The sales order includes:
- Customer information (if applicable)
- Complete item list with pricing
- Tax calculations
- Payment information
- Inventory updates
- Invoice generation

The process must be atomic to ensure data consistency.

---

## ✅ Acceptance Criteria

### Scenario 1: Create sales order from cart
- Given that a cart exists with items
- When I convert the cart to a sales order
- Then a sales order must be created with:
  - Unique order number
  - Customer information
  - Complete item details
  - Calculated taxes
  - Total amount
  - Status: PENDING_PAYMENT
- And inventory must be updated for all items

### Scenario 2: Tax calculation
- Given that products have different tax rates
- When the sales order is created
- Then taxes must be calculated per item
- And the total tax must be the sum of all item taxes

### Scenario 3: Inventory update
- Given that products have current stock
- When a sales order is completed
- Then stock levels must be decremented
- And low stock alerts must be triggered if needed

### Scenario 4: Order status management
- Given that a sales order exists
- When payment is processed
- Then the order status must update to PAID
- And the order must be marked as completed

### Scenario 5: Multi-payment support
- Given that a customer wants to split payment
- When multiple payment methods are used
- Then the sales order must track all payments
- And the total must equal the sum of payments

---

## ❌ Error Cases

- Converting empty cart must return error 400
- Insufficient stock must prevent order creation
- Invalid customer data must return error 400
- Order creation without required fields must return validation errors

---

## 🔐 Business Rules

- Sales orders are tied to the creating user and branch
- Inventory is updated atomically to prevent overselling
- Orders with insufficient stock are rejected
- Tax calculations follow Colombian regulations
- Each order gets a unique sequential number per branch

---

## 🗄️ Database Schema (PostgreSQL)

### Table: sales_order
```sql
CREATE TABLE sales_order (
    id BIGSERIAL PRIMARY KEY,
    order_number VARCHAR(50) NOT NULL,
    customer_id INTEGER REFERENCES customer(id),
    user_id INTEGER NOT NULL REFERENCES usuario(id),
    branch_id INTEGER NOT NULL REFERENCES branch(id),
    empresa_id INTEGER NOT NULL REFERENCES empresa(id),
    subtotal DECIMAL(12,2) NOT NULL,
    discount DECIMAL(12,2) NOT NULL DEFAULT 0,
    tax_total DECIMAL(12,2) NOT NULL,
    total DECIMAL(12,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING_PAYMENT' CHECK (status IN ('PENDING_PAYMENT', 'PAID', 'CANCELLED', 'COMPLETED')),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    
    CONSTRAINT sales_order_empresa_fk FOREIGN KEY (empresa_id) REFERENCES empresa(id),
    CONSTRAINT sales_order_user_fk FOREIGN KEY (user_id) REFERENCES usuario(id),
    CONSTRAINT sales_order_branch_fk FOREIGN KEY (branch_id) REFERENCES branch(id),
    CONSTRAINT sales_order_customer_fk FOREIGN KEY (customer_id) REFERENCES customer(id),
    CONSTRAINT sales_order_number_unique UNIQUE (empresa_id, branch_id, order_number)
);

CREATE INDEX idx_sales_order_empresa ON sales_order(empresa_id);
CREATE INDEX idx_sales_order_branch ON sales_order(branch_id);
CREATE INDEX idx_sales_order_customer ON sales_order(customer_id);
CREATE INDEX idx_sales_order_status ON sales_order(status);
```

### Table: sales_order_item
```sql
CREATE TABLE sales_order_item (
    id BIGSERIAL PRIMARY KEY,
    sales_order_id BIGINT NOT NULL REFERENCES sales_order(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES product(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(12,2) NOT NULL,
    discount_percent DECIMAL(5,2) DEFAULT 0,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    tax_rate DECIMAL(5,2) NOT NULL,
    tax_amount DECIMAL(12,2) NOT NULL,
    total DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT sales_order_item_order_fk FOREIGN KEY (sales_order_id) REFERENCES sales_order(id) ON DELETE CASCADE,
    CONSTRAINT sales_order_item_product_fk FOREIGN KEY (product_id) REFERENCES product(id)
);

CREATE INDEX idx_sales_order_item_order ON sales_order_item(sales_order_id);
CREATE INDEX idx_sales_order_item_product ON sales_order_item(product_id);
```

---

## 🎨 UI/UX Considerations

- Order summary with item breakdown
- Tax calculation preview
- Stock availability indicators
- Customer selection with search
- Order confirmation dialog

---

## 📡 Technical Requirements

### Entities (Java)

**SalesOrderEntity:**
- id (Long, PK)
- order_number (String, unique)
- customer_id (Long, FK, nullable)
- user_id (Long, FK)
- branch_id (Long, FK)
- empresa_id (Long, FK)
- subtotal (BigDecimal)
- discount (BigDecimal)
- tax_total (BigDecimal)
- total (BigDecimal)
- status (String): PENDING_PAYMENT, PAID, CANCELLED, COMPLETED
- notes (String, TEXT)
- created_at (LocalDateTime)
- updated_at (LocalDateTime)

**SalesOrderItemEntity:**
- id (Long, PK)
- sales_order_id (Long, FK)
- product_id (Long, FK)
- quantity (Integer)
- unit_price (BigDecimal)
- discount_percent (BigDecimal)
- discount_amount (BigDecimal)
- tax_rate (BigDecimal)
- tax_amount (BigDecimal)
- total (BigDecimal)
- created_at (LocalDateTime)

### Services

- `SalesOrderService` with methods:
  - `createFromCart(CartDto cart)`
  - `getOrderById(Long id)`
  - `updateOrderStatus(Long id, String status)`
  - `cancelOrder(Long id)`

### Controllers

- `SalesOrderController` with endpoints:
  - `POST /api/sales-orders` - Create from cart
  - `GET /api/sales-orders/{id}` - Get order details
  - `PUT /api/sales-orders/{id}/status` - Update status
  - `POST /api/sales-orders/page` - Paginated list

---

## 🚫 Out of Scope

- Advanced order editing after creation
- Order splitting and merging
- Complex shipping logistics
- Returns and exchanges management

---

## 📎 Dependencies

- HU-SALES-001: Product Catalog
- HU-SALES-002: Shopping Cart
- Existing Customer entity
- Inventory management system
