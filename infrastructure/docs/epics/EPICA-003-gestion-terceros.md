# EPICA-003: Third Parties Management

## 📌 General Information
- ID: EPICA-003
- State: Completed
- Priority: High
- Start Date: 2026-03-23
- Target Date: 2026-04-15
- Owner: Aura POS Backend Team
- Percentage: 100%

---

## 🎯 Business Objective

Manage third parties (clients, suppliers, and employees) for each tenant in Aura POS. Every third party must be stored in the tenant schema with full contact information, document validation, and categorization. This module is the foundation for sales, purchases, and payroll operations.

**What problem does it solve?**
- Centralized management of all business relationships
- Unified contact information across all modules
- Support for different third party types with specific attributes

**What value does it generate?**
- Data consistency across sales and purchases
- Efficient customer lookup during sales
- Better vendor management for procurement

---

## 👥 Stakeholders

- End User: Store managers, cashiers, accountants
- Technical Team: Backend developers
- Product: Product Manager Aura POS

---

## 🧠 Functional Description

The Third Parties module manages all external and internal entities that interact with the business:

1. **Clients (Customers)**: Individuals or companies that purchase products/services
2. **Suppliers (Vendors)**: Companies or individuals that provide products/services
3. **Employees**: Staff members including cashiers, managers, and technicians

Each third party belongs to a tenant schema, contains contact details, document information, and can be associated with sales or purchases.

---

## 📦 Scope

### Included:
- CRUD operations for third parties
- Document validation (NIT, CC, CE, passport)
- Contact information management (address, phone, email)
- Third party categorization by type
- Search and filter capabilities
- Bulk import/export
- Customer loyalty points tracking

### Not Included:
- Advanced CRM features
- Marketing campaign management
- Customer segmentation analytics
- Integration with external ID validation services

---

## 🧩 User Stories

| HU | Title | State |
|----|-------|-------|
| HU-THIRD-001 | Create Third Party | ✅ Completed |
| HU-THIRD-002 | List Third Parties | ✅ Completed |
| HU-THIRD-003 | Get Third Party Details | ✅ Completed |
| HU-THIRD-004 | Update Third Party | ✅ Completed |
| HU-THIRD-005 | Change Third Party Status | ✅ Completed |
| HU-THIRD-006 | Search Third Parties | ✅ Completed |

---

## 🐞 Associated Bugs

None identified

---

## 🔐 Global Business Rules

- All third parties belong to a tenant schema
- Document numbers must be unique per tenant and type
- Email validation follows RFC 5322 standard
- Phone numbers follow E.164 format
- Soft delete preserves data integrity
- Third party types: CLIENT, SUPPLIER, EMPLOYEE
- Each employee must be linked to a user account

---

## 🧱 Related Architecture

**Backend:** Go 1.26.1 with Gin framework
**Database:** PostgreSQL with schema-per-tenant
**Authentication:** JWT with tenant context

### Database Schema (Tenant Schema)

**Table: third_parties**
```sql
CREATE TABLE IF NOT EXISTS third_parties (
    id                  BIGSERIAL PRIMARY KEY,
    user_id             BIGINT REFERENCES public.users(id),
    first_name          VARCHAR(100),
    last_name           VARCHAR(100),
    document_number     VARCHAR(50) NOT NULL,
    document_type       VARCHAR(20) NOT NULL,
    personal_email      VARCHAR(150),
    commercial_name     VARCHAR(255),
    address             VARCHAR(255),
    phone               VARCHAR(20),
    additional_email    VARCHAR(150),
    tax_responsibility  VARCHAR(20) NOT NULL CHECK (tax_responsibility IN ('RESPONSIBLE', 'NOT-RESPONSIBLE')),
    is_client           BOOLEAN NOT NULL DEFAULT FALSE,
    is_provider         BOOLEAN NOT NULL DEFAULT FALSE,
    is_employee         BOOLEAN NOT NULL DEFAULT FALSE,
    municipality_id     VARCHAR(10),
    municipality        VARCHAR(255),
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ DEFAULT NULL
);
```

---

## 📊 Success Metrics

- 100% of third parties with valid documents
- Response time < 150ms for queries
- Data accuracy > 99%
- Unique document validation 100%

---

## 🚧 Risks

- Document number validation complexity (different formats per type)
- Data migration from legacy systems
- Duplicate third party detection

---

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/third-parties` | Create a new third party |
| GET | `/third-parties` | List third parties with filters |
| GET | `/third-parties/:id` | Get third party by ID |
| GET | `/third-parties/document/:documentNumber` | Get third party by document number |
| PUT | `/third-parties/:id` | Update third party |
| DELETE | `/third-parties/:id` | Delete third party (soft delete) |

### Query Parameters for List
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 20)
- `type` - Filter by type: client, provider, employee
- `search` - Search by name or document number

---

## 📁 Module Structure

```
modules/third-parties/
├── domain.go     # Entity, Repository & Service interfaces
├── service.go    # Repository & Service implementation
├── handler.go    # HTTP handlers
└── routes.go     # Route registration
```

---

## Resumen

- **Total de HU**: 6
- **Completadas**: 6
- **Pendientes**: 0
- **Módulo implementado**: third-parties
