# HU-SALES-001 - Product Catalog Management

## 📌 General Information
- ID: HU-SALES-001
- Epic: EPIC-SALES-001
- Priority: High
- State: Backlog
- Progress: 0%
- Author: QA Engineer Aura POS
- Date: 2026-03-15

---

## 👤 User Story

**As a** store manager
**I want to** manage the product catalog with categories, brands, and pricing
**So that** cashiers can quickly find and sell products to customers

---

## 🧠 Functional Description

The system must allow complete management of the product catalog including:
- Product creation with SKU, name, description, and images
- Category organization with hierarchical structure
- Brand management
- Pricing with cost, sale price, and discount limits
- Inventory tracking with minimum stock alerts
- Product variants (size, color, etc.)

All products must be filterable and searchable by multiple criteria.

---

## ✅ Acceptance Criteria

### Scenario 1: Create a new product
- Given that I am logged in as a manager with product creation permissions
- When I create a new product with:
  - SKU: "PROD-001"
  - Name: "Wireless Mouse"
  - Category: "Electronics"
  - Brand: "TechBrand"
  - Cost Price: $50,000
  - Sale Price: $80,000
  - Tax Rate: 19%
- Then the product must be saved in the database with:
  - Unique SKU validation
  - Automatic inventory record creation with zero stock
  - Timestamps (created_at, updated_at)
  - Status: ACTIVE

### Scenario 2: Update product information
- Given that a product exists in the catalog
- When I update the product details (price, description, category)
- Then the changes must be persisted with updated_at timestamp
- And the history of price changes must be auditable

### Scenario 3: Search and filter products
- Given that the catalog contains multiple products
- When I search by name, SKU, or category
- Then the system returns matching products with pagination
- And allows filtering by price range, brand, or stock status

### Scenario 4: Manage product categories
- Given that categories exist in a hierarchical structure
- When I create a new subcategory
- Then it must be linked to the parent category
- And products can be assigned to the appropriate category level

### Scenario 5: Inventory tracking
- Given that a product has stock
- When a sale is completed
- Then the inventory must be updated automatically
- And low stock alerts must be triggered when below minimum threshold

---

## ❌ Error Cases

- SKU duplication must return error 400
- Invalid price values (negative or zero) must return error 400
- Non-existent category assignment must return error 404
- Product deletion must be soft delete only
- Required fields missing must return validation errors

---

## 🔐 Business Rules

- Only users with MANAGER or ADMIN role can create/edit products
- Sale price must be greater than or equal to cost price
- SKU must be unique per company
- Products cannot be physically deleted (soft delete only)
- Inventory updates must be atomic and thread-safe

---

## 🗄️ Database Schema (PostgreSQL)

### Table: product
```sql
CREATE TABLE product (
    id BIGSERIAL PRIMARY KEY,
    sku VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    category_id INTEGER NOT NULL REFERENCES category(id),
    brand_id INTEGER REFERENCES brand(id),
    cost_price DECIMAL(12,2) NOT NULL,
    sale_price DECIMAL(12,2) NOT NULL,
    tax_rate DECIMAL(5,2) NOT NULL DEFAULT 19.00,
    min_stock INTEGER DEFAULT 0,
    current_stock INTEGER DEFAULT 0,
    image_url VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE', 'DISCONTINUED')),
    empresa_id INTEGER NOT NULL REFERENCES empresa(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT product_empresa_fk FOREIGN KEY (empresa_id) REFERENCES empresa(id),
    CONSTRAINT product_category_fk FOREIGN KEY (category_id) REFERENCES category(id),
    CONSTRAINT product_brand_fk FOREIGN KEY (brand_id) REFERENCES brand(id),
    CONSTRAINT product_sku_unique UNIQUE (empresa_id, sku),
    CONSTRAINT product_price_check CHECK (sale_price >= cost_price)
);

CREATE INDEX idx_product_empresa ON product(empresa_id);
CREATE INDEX idx_product_category ON product(category_id);
CREATE INDEX idx_product_sku ON product(sku);
CREATE INDEX idx_product_status ON product(status);
CREATE INDEX idx_product_deleted_at ON product(deleted_at) WHERE deleted_at IS NULL;
```

### Table: category
```sql
CREATE TABLE category (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    parent_id INTEGER REFERENCES category(id),
    empresa_id INTEGER NOT NULL REFERENCES empresa(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    
    CONSTRAINT category_empresa_fk FOREIGN KEY (empresa_id) REFERENCES empresa(id),
    CONSTRAINT category_parent_fk FOREIGN KEY (parent_id) REFERENCES category(id),
    CONSTRAINT category_name_unique UNIQUE (empresa_id, name)
);

CREATE INDEX idx_category_empresa ON category(empresa_id);
CREATE INDEX idx_category_parent ON category(parent_id);
```

### Table: brand
```sql
CREATE TABLE brand (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    empresa_id INTEGER NOT NULL REFERENCES empresa(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    
    CONSTRAINT brand_empresa_fk FOREIGN KEY (empresa_id) REFERENCES empresa(id),
    CONSTRAINT brand_name_unique UNIQUE (empresa_id, name)
);

CREATE INDEX idx_brand_empresa ON brand(empresa_id);
```

---

## 🎨 UI/UX Considerations

- Product list with searchable table view
- Image upload and preview functionality
- Category tree selector for product assignment
- Price calculator showing profit margin
- Stock level indicators with color coding

---

## 📡 Technical Requirements

### Entities (Java)

**ProductEntity:**
- id (Long, PK, autoincrement)
- sku (String, VARCHAR 50)
- name (String, VARCHAR 200)
- description (String, TEXT)
- category_id (Long, FK)
- brand_id (Long, FK)
- cost_price (BigDecimal)
- sale_price (BigDecimal)
- tax_rate (BigDecimal)
- min_stock (Integer)
- current_stock (Integer)
- image_url (String, VARCHAR 500)
- status (String): ACTIVE, INACTIVE, DISCONTINUED
- empresa_id (Long, FK)
- created_at (LocalDateTime)
- updated_at (LocalDateTime, nullable)
- deleted_at (LocalDateTime, nullable)

### Services

- `ProductService` with methods:
  - `createProduct(CreateProductDto dto)`
  - `updateProduct(Long id, UpdateProductDto dto)`
  - `getProductById(Long id)`
  - `searchProducts(PageableDto pageable, String search)`
  - `deleteProduct(Long id)` (soft delete)

### Controllers

- `ProductController` with endpoints:
  - `POST /api/products` - Create product
  - `PUT /api/products/{id}` - Update product
  - `GET /api/products/{id}` - Get product details
  - `POST /api/products/page` - Paginated search
  - `DELETE /api/products/{id}` - Soft delete product

---

## 🚫 Out of Scope

- Product variant management (separate HU)
- Advanced image management with CDN
- Product reviews and ratings
- Supplier management integration
- Advanced pricing rules engine

---

## 📎 Dependencies

- EPIC-SALES-001: Sales Module Epic
- Existing Empresa and Usuario entities
- Spring Security configuration
