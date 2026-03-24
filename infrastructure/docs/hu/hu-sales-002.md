# HU-SALES-002 - Shopping Cart / Quote Management

## 📌 General Information

- ID: HU-SALES-002
- Epic: EPIC-SALES-001
- Priority: High
- State: Backlog
- Progress: 0%
- Author: QA Engineer Aura POS
- Date: 2026-03-15

---

## 👤 User Story

**As a** cashier or salesperson
**I want to** create and manage shopping carts/quotes with multiple products
**So that** I can prepare sales orders for customers efficiently

---

## 🧠 Functional Description

The system must allow creating temporary shopping carts that can be converted to sales orders. Carts should support:

- Adding/removing products with quantities
- Applying discounts at cart or item level
- Calculating taxes automatically
- Saving quotes for later retrieval
- Converting quotes to sales orders

Carts are temporary until converted to sales orders, after which they become permanent records.

---

## ✅ Acceptance Criteria

### Scenario 1: Create a new shopping cart

- Given that I am logged in as a cashier
- When I start a new sale
- Then a new cart must be created with:
  - Unique cart identifier
  - Empty item list
  - Timestamp of creation
  - User ID of the cashier
  - Branch ID from JWT

### Scenario 2: Add products to cart

- Given that a cart exists
- When I add a product with quantity
- Then the product must be added to the cart items
- And the cart total must be recalculated including taxes

### Scenario 3: Update cart item quantities

- Given that a product is in the cart
- When I change the quantity
- Then the item total must be recalculated
- And the cart grand total must be updated

### Scenario 4: Apply discount to cart

- Given that a cart has items
- When I apply a discount percentage or fixed amount
- Then the discount must be applied to the total
- And the discount must be validated against user permissions

### Scenario 5: Save quote for later

- Given that I have a cart with items
- When I save the cart as a quote
- Then the quote must be stored with customer information
- And I can retrieve it later to continue editing

### Scenario 6: Convert quote to sale

- Given that I have a saved quote
- When I convert it to a sale
- Then a sales order must be created
- And the cart must be cleared or archived

---

## ❌ Error Cases

- Adding non-existent product must return error 404
- Adding product with zero/negative quantity must return error 400
- Discount exceeding maximum allowed must return error 400
- Cart conversion without items must return error 400
- Accessing another user's cart must return error 403

---

## 🔐 Business Rules

- Carts are tied to the creating user and branch
- Discounts require appropriate permissions (manager approval for >10%)
- Cart items are immutable after conversion to sale
- Quotes can be retrieved by the creating user or managers
- Cart expiration: 24 hours for temporary carts

---

## 🗄️ Database Schema (PostgreSQL)

### Table: cart

```sql
CREATE TABLE cart (
    id BIGSERIAL PRIMARY KEY,
    cart_code VARCHAR(50) NOT NULL,
    customer_id INTEGER REFERENCES customer(id),
    user_id INTEGER NOT NULL REFERENCES usuario(id),
    branch_id INTEGER NOT NULL REFERENCES branch(id),
    enterprise_id INTEGER NOT NULL REFERENCES empresa(id),
    subtotal DECIMAL(12,2) NOT NULL DEFAULT 0,
    discount DECIMAL(12,2) NOT NULL DEFAULT 0,
    tax_total DECIMAL(12,2) NOT NULL DEFAULT 0,
    total DECIMAL(12,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'SAVED', 'CONVERTED', 'EXPIRED')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,

    CONSTRAINT cart_empresa_fk FOREIGN KEY (enterprise_id) REFERENCES empresa(id),
    CONSTRAINT cart_user_fk FOREIGN KEY (user_id) REFERENCES usuario(id),
    CONSTRAINT cart_branch_fk FOREIGN KEY (branch_id) REFERENCES branch(id),
    CONSTRAINT cart_customer_fk FOREIGN KEY (customer_id) REFERENCES customer(id),
    CONSTRAINT cart_code_unique UNIQUE (enterprise_id, cart_code)
);

CREATE INDEX idx_cart_empresa ON cart(enterprise_id);
CREATE INDEX idx_cart_user ON cart(user_id);
CREATE INDEX idx_cart_status ON cart(status);
```

### Table: cart_item

```sql
CREATE TABLE cart_item (
    id BIGSERIAL PRIMARY KEY,
    cart_id BIGINT NOT NULL REFERENCES cart(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES product(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(12,2) NOT NULL,
    discount_percent DECIMAL(5,2) DEFAULT 0,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    tax_rate DECIMAL(5,2) NOT NULL,
    tax_amount DECIMAL(12,2) NOT NULL,
    total DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,

    CONSTRAINT cart_item_cart_fk FOREIGN KEY (cart_id) REFERENCES cart(id) ON DELETE CASCADE,
    CONSTRAINT cart_item_product_fk FOREIGN KEY (product_id) REFERENCES product(id)
);

CREATE INDEX idx_cart_item_cart ON cart_item(cart_id);
CREATE INDEX idx_cart_item_product ON cart_item(product_id);
```

---

## 🎨 UI/UX Considerations

- Real-time cart updates with visual feedback
- Product search with barcode scanner support
- Discount calculation preview
- Customer selection dialog
- Quote save/retrieve interface

---

## 📡 Technical Requirements

### Entities (Java)

**CartEntity:**

- id (Long, PK)
- cart_code (String, unique)
- customer_id (Long, FK, nullable)
- user_id (Long, FK)
- branch_id (Long, FK)
- enterprise_id (Long, FK)
- subtotal (BigDecimal)
- discount (BigDecimal)
- tax_total (BigDecimal)
- total (BigDecimal)
- status (String): ACTIVE, SAVED, CONVERTED, EXPIRED
- created_at (LocalDateTime)
- updated_at (LocalDateTime)

**CartItemEntity:**

- id (Long, PK)
- cart_id (Long, FK)
- product_id (Long, FK)
- quantity (Integer)
- unit_price (BigDecimal)
- discount_percent (BigDecimal)
- discount_amount (BigDecimal)
- tax_rate (BigDecimal)
- tax_amount (BigDecimal)
- total (BigDecimal)
- created_at (LocalDateTime)
- updated_at (LocalDateTime)

### Controllers

- `CartController` with endpoints:
  - `POST /api/carts` - Create new cart
  - `GET /api/carts/{id}` - Get cart details
  - `POST /api/carts/{cartId}/items` - Add item to cart
  - `PUT /api/carts/{cartId}/items/{itemId}` - Update cart item
  - `DELETE /api/carts/{cartId}/items/{itemId}` - Remove item
  - `POST /api/carts/{cartId}/discount` - Apply discount
  - `POST /api/carts/{cartId}/convert` - Convert to sale

---

## 🚫 Out of Scope

- Complex pricing rules engine
- Customer credit limit checking
- Integration with external CRM
- Advanced quote templates

---

## 📎 Dependencies

- HU-SALES-001: Product Catalog
- Existing Customer entity (from EP-002)
- Branch and Empresa entities
