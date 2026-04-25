package offline

import (
	"context"
	"encoding/json"
	"time"
 
	"github.com/cloud-tech-develop/aura-back/internal/db"
)

type repository struct {
	db *db.DB
}

func NewRepository(database *db.DB) Repository {
	return &repository{db: database}
}

// ─── Enterprise Operations ─────────────────────────────────────────────────────

func (r *repository) UpsertEnterprise(ctx context.Context, e *Enterprise) error {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM enterprises WHERE id = ?)", e.ID).Scan(&exists)
	if err != nil {
		return err
	}

	settingsJSON, _ := json.Marshal(e.Settings)
	if settingsJSON == nil || string(settingsJSON) == "null" {
		settingsJSON = []byte("{}")
	}

	if exists {
		_, err = r.db.ExecContext(ctx,
			`UPDATE enterprises SET 
			 tenant_id = ?, name = ?, commercial_name = ?, slug = ?, sub_domain = ?, email = ?, 
			 document = ?, dv = ?, phone = ?, municipality_id = ?, municipality = ?, 
			 status = ?, settings = ?, updated_at = ?
			 WHERE id = ?`,
			e.TenantID, e.Name, e.CommercialName, e.Slug, e.SubDomain, e.Email,
			e.Document, e.DV, e.Phone, e.MunicipalityID, e.Municipality,
			e.Status, settingsJSON, e.UpdatedAt, e.ID,
		)
	} else {
		_, err = r.db.ExecContext(ctx,
			`INSERT INTO enterprises 
			 (id, tenant_id, name, commercial_name, slug, sub_domain, email, document, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at, deleted_at) 
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			e.ID, e.TenantID, e.Name, e.CommercialName, e.Slug, e.SubDomain,
			e.Email, e.Document, e.DV, e.Phone, e.MunicipalityID, e.Municipality,
			e.Status, settingsJSON, e.CreatedAt, e.UpdatedAt, nil,
		)
	}
	return err
}

func (r *repository) updateEnterprise(ctx context.Context, e *Enterprise) error {
	e.UpdatedAt = time.Now()

	settingsJSON, err := json.Marshal(e.Settings)
	if err != nil || string(settingsJSON) == "null" {
		settingsJSON = []byte("{}")
	}

	_, err = r.db.ExecContext(ctx,
		`UPDATE enterprises SET 
		 tenant_id = ?, name = ?, commercial_name = ?, slug = ?, sub_domain = ?, email = ?, document = ?, dv = ?, phone = ?, municipality_id = ?, municipality = ?, status = ?, settings = ?, updated_at = ?
		 WHERE slug = ?`,
		e.TenantID, e.Name, e.CommercialName, e.Slug, e.SubDomain,
		e.Email, e.Document, e.DV, e.Phone, e.MunicipalityID, e.Municipality,
		e.Status, settingsJSON, e.UpdatedAt, e.Slug,
	)
	return err
}

func (r *repository) GetEnterpriseBySlug(ctx context.Context, slug string) (*Enterprise, error) {
	var e Enterprise
	var settingsJSON []byte
	err := r.db.QueryRowContext(ctx,
		`SELECT id, tenant_id, name, commercial_name, slug, sub_domain, email, document, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at 
		 FROM enterprises WHERE slug = ? AND deleted_at IS NULL`,
		slug,
	).Scan(&e.ID, &e.TenantID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain,
		&e.Email, &e.Document, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality,
		&e.Status, &settingsJSON, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(settingsJSON, &e.Settings)
	return &e, nil
}

func (r *repository) ListEnterprises(ctx context.Context) ([]Enterprise, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, tenant_id, name, commercial_name, slug, sub_domain, email, document, dv, phone, municipality_id, municipality, status, settings, created_at, updated_at 
		 FROM enterprises WHERE deleted_at IS NULL ORDER BY name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Enterprise
	for rows.Next() {
		var e Enterprise
		var settingsJSON []byte
		if err := rows.Scan(&e.ID, &e.TenantID, &e.Name, &e.CommercialName, &e.Slug, &e.SubDomain,
			&e.Email, &e.Document, &e.DV, &e.Phone, &e.MunicipalityID, &e.Municipality,
			&e.Status, &settingsJSON, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(settingsJSON, &e.Settings)
		list = append(list, e)
	}
	return list, nil
}

// ─── Plan Operations ─────────────────────────────────────────────────────

func (r *repository) UpsertPlan(ctx context.Context, p *Plan) error {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM plans WHERE id = ?)", p.ID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = r.db.ExecContext(ctx,
			`UPDATE plans SET enterprise_id = ?, max_users = ?, max_enterprises = ?, trial_until = ?, updated_at = ?
			 WHERE id = ?`,
			p.EnterpriseID, p.MaxUsers, p.MaxEnterprises, p.TrialUntil, time.Now(), p.ID,
		)
	} else {
		p.CreatedAt = time.Now()
		p.UpdatedAt = time.Now()
		_, err = r.db.ExecContext(ctx,
			`INSERT INTO plans (id, enterprise_id, max_users, max_enterprises, trial_until, created_at, updated_at, deleted_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			p.ID, p.EnterpriseID, p.MaxUsers, p.MaxEnterprises, p.TrialUntil, p.CreatedAt, p.UpdatedAt, nil,
		)
	}
	return err
}

// ─── User Operations ─────────────────────────────────────────────────────

func (r *repository) UpsertUser(ctx context.Context, u *User) error {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", u.ID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = r.db.ExecContext(ctx,
			`UPDATE users SET enterprise_id = ?, name = ?, email = ?, active = ?, password_hash = ?, updated_at = ?
			 WHERE id = ?`,
			u.EnterpriseID, u.Name, u.Email, u.Active, u.PasswordHash, time.Now(), u.ID,
		)
	} else {
		u.CreatedAt = time.Now()
		u.UpdatedAt = time.Now()
		_, err = r.db.ExecContext(ctx,
			`INSERT INTO users (id, enterprise_id, name, email, active, password_hash, created_at, updated_at, deleted_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			u.ID, u.EnterpriseID, u.Name, u.Email, u.Active, u.PasswordHash, u.CreatedAt, u.UpdatedAt, nil,
		)
	}
	return err
}

// ─── UserRole Operations ─────────────────────────────────────────────────────

func (r *repository) UpsertUserRole(ctx context.Context, ur *UserRole) error {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM user_roles WHERE id = ?)", ur.ID).Scan(&exists)
	if err != nil {
		return err
	}

if exists {
		_, err = r.db.ExecContext(ctx,
			`UPDATE user_roles SET user_id = ?, role_id = ? WHERE id = ?`,
			ur.UserID, ur.RoleID, ur.ID,
		)
	} else {
		_, err = r.db.ExecContext(ctx,
			`INSERT INTO user_roles (id, user_id, role_id) VALUES (?, ?, ?)`,
			ur.ID, ur.UserID, ur.RoleID,
		)
	}
	return err
}

// ─── Tenant Operations: ThirdParty ─────────────────────────────────────────────────────

func (r *repository) UpsertThirdParty(ctx context.Context, tp *ThirdParty) error {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM third_parties WHERE id = ?)", tp.ID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = r.db.ExecContext(ctx,
			`UPDATE third_parties SET 
			 user_id = ?, first_name = ?, last_name = ?, document_number = ?, document_type = ?, 
			 personal_email = ?, commercial_name = ?, address = ?, phone = ?, additional_email = ?, 
			 tax_responsibility = ?, is_client = ?, is_provider = ?, is_employee = ?, 
			 municipality_id = ?, municipality = ?, enterprise_id = ?
			 WHERE id = ?`,
			tp.UserID, tp.FirstName, tp.LastName, tp.DocumentNumber, tp.DocumentType,
			tp.PersonalEmail, tp.CommercialName, tp.Address, tp.Phone, tp.AdditionalEmail,
			tp.TaxResponsibility, tp.IsClient, tp.IsProvider, tp.IsEmployee,
			tp.MunicipalityID, tp.Municipality, tp.EnterpriseID, tp.ID,
		)
	} else {
		tp.CreatedAt = time.Now()
		_, err = r.db.ExecContext(ctx,
			`INSERT INTO third_parties 
			 (id, user_id, first_name, last_name, document_number, document_type, personal_email, commercial_name, 
			  address, phone, additional_email, tax_responsibility, is_client, is_provider, is_employee, 
			  municipality_id, municipality, enterprise_id, created_at, deleted_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			tp.ID, tp.UserID, tp.FirstName, tp.LastName, tp.DocumentNumber, tp.DocumentType,
			tp.PersonalEmail, tp.CommercialName, tp.Address, tp.Phone, tp.AdditionalEmail,
			tp.TaxResponsibility, tp.IsClient, tp.IsProvider, tp.IsEmployee,
			tp.MunicipalityID, tp.Municipality, tp.EnterpriseID, tp.CreatedAt, nil,
		)
	}
	return err
}

// ─── Tenant Operations: Category ─────────────────────────────────────────────────────

func (r *repository) UpsertCategory(ctx context.Context, c *Category) error {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM category WHERE id = ?)", c.ID).Scan(&exists)
	if err != nil {
		return err
	}

	active := 0
	if c.Active {
		active = 1
	}
	parentID := 0
	if c.ParentID != nil {
		parentID = int(*c.ParentID)
	}

	if exists {
		_, err = r.db.ExecContext(ctx,
			`UPDATE category SET name = ?, description = ?, parent_id = ?, default_tax_rate = ?, active = ?, enterprise_id = ?, updated_at = ?
			 WHERE id = ?`,
			c.Name, c.Description, parentID, c.DefaultTaxRate, active, c.EnterpriseID, time.Now(), c.ID,
		)
	} else {
		c.CreatedAt = time.Now()
		c.UpdatedAt = time.Now()
		_, err = r.db.ExecContext(ctx,
			`INSERT INTO category (id, name, description, parent_id, default_tax_rate, active, enterprise_id, created_at, updated_at, deleted_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			c.ID, c.Name, c.Description, parentID, c.DefaultTaxRate, active, c.EnterpriseID, c.CreatedAt, c.UpdatedAt, nil,
		)
	}
	return err
}

// ─── Tenant Operations: Brand ─────────────────────────────────────────────────────

func (r *repository) UpsertBrand(ctx context.Context, b *Brand) error {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM brand WHERE id = ?)", b.ID).Scan(&exists)
	if err != nil {
		return err
	}

	active := 0
	if b.Active {
		active = 1
	}

	if exists {
		_, err = r.db.ExecContext(ctx,
			`UPDATE brand SET name = ?, description = ?, active = ?, enterprise_id = ?, updated_at = ?
			 WHERE id = ?`,
			b.Name, b.Description, active, b.EnterpriseID, time.Now(), b.ID,
		)
	} else {
		b.CreatedAt = time.Now()
		b.UpdatedAt = time.Now()
		_, err = r.db.ExecContext(ctx,
			`INSERT INTO brand (id, name, description, active, enterprise_id, created_at, updated_at, deleted_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			b.ID, b.Name, b.Description, active, b.EnterpriseID, b.CreatedAt, b.UpdatedAt, nil,
		)
	}
	return err
}

// ─── Tenant Operations: Unit ─────────────────────────────────────────────────────

func (r *repository) UpsertUnit(ctx context.Context, u *Unit) error {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM unit WHERE id = ?)", u.ID).Scan(&exists)
	if err != nil {
		return err
	}

	active := 0
	if u.Active {
		active = 1
	}
	allowDecimals := 0
	if u.AllowDecimals {
		allowDecimals = 1
	}

	if exists {
		_, err = r.db.ExecContext(ctx,
			`UPDATE unit SET name = ?, abbreviation = ?, active = ?, allow_decimals = ?, enterprise_id = ?, updated_at = ?
			 WHERE id = ?`,
			u.Name, u.Abbreviation, active, allowDecimals, u.EnterpriseID, time.Now(), u.ID,
		)
	} else {
		u.CreatedAt = time.Now()
		u.UpdatedAt = time.Now()
		_, err = r.db.ExecContext(ctx,
			`INSERT INTO unit (id, name, abbreviation, active, allow_decimals, enterprise_id, created_at, updated_at, deleted_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			u.ID, u.Name, u.Abbreviation, active, allowDecimals, u.EnterpriseID, u.CreatedAt, u.UpdatedAt, nil,
		)
	}
	return err
}

// ─── Tenant Operations: Product ─────────────────────────────────────────────────────

func (r *repository) UpsertProduct(ctx context.Context, p *Product) error {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM product WHERE id = ?)", p.ID).Scan(&exists)
	if err != nil {
		return err
	}

	active := 0
	if p.Active {
		active = 1
	}
	visibleInPOS := 0
	if p.VisibleInPOS {
		visibleInPOS = 1
	}
	managesInventory := 0
	if p.ManagesInventory {
		managesInventory = 1
	}
	managesBatches := 0
	if p.ManagesBatches {
		managesBatches = 1
	}
	managesSerial := 0
	if p.ManagesSerial {
		managesSerial = 1
	}
	allowNegativeStock := 0
	if p.AllowNegativeStock {
		allowNegativeStock = 1
	}

	categoryID := 0
	if p.CategoryID != nil {
		categoryID = int(*p.CategoryID)
	}
	brandID := 0
	if p.BrandID != nil {
		brandID = int(*p.BrandID)
	}

	if exists {
		_, err = r.db.ExecContext(ctx,
			`UPDATE product SET 
			 sku = ?, barcode = ?, name = ?, description = ?, category_id = ?, brand_id = ?, unit_id = ?, 
			 product_type = ?, active = ?, visible_in_pos = ?, cost_price = ?, sale_price = ?, 
			 price_2 = ?, price_3 = ?, iva_percentage = ?, consumption_tax_value = ?, 
			 current_stock = ?, min_stock = ?, max_stock = ?, manages_inventory = ?, manages_batches = ?, 
			 manages_serial = ?, allow_negative_stock = ?, image_url = ?, updated_at = ?
			 WHERE id = ?`,
			p.SKU, p.Barcode, p.Name, p.Description, categoryID, brandID, p.UnitID,
			p.ProductType, active, visibleInPOS, p.CostPrice, p.SalePrice, p.Price2, p.Price3,
			p.IVAPercentage, p.ConsumptionTax,
			p.CurrentStock, p.MinStock, p.MaxStock, managesInventory, managesBatches,
			managesSerial, allowNegativeStock, p.ImageURL, time.Now(), p.ID,
		)
	} else {
		p.CreatedAt = time.Now()
		p.UpdatedAt = time.Now()
		_, err = r.db.ExecContext(ctx,
			`INSERT INTO product 
			 (id, sku, barcode, name, description, category_id, brand_id, unit_id, product_type, active, visible_in_pos, 
			  cost_price, sale_price, price_2, price_3, iva_percentage, consumption_tax_value, 
			  current_stock, min_stock, max_stock, manages_inventory, manages_batches, manages_serial, 
			  allow_negative_stock, image_url, enterprise_id, created_at, updated_at, deleted_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			p.ID, p.SKU, p.Barcode, p.Name, p.Description, categoryID, brandID, p.UnitID,
			p.ProductType, active, visibleInPOS, p.CostPrice, p.SalePrice, p.Price2, p.Price3,
			p.IVAPercentage, p.ConsumptionTax, p.CurrentStock, p.MinStock, p.MaxStock,
			managesInventory, managesBatches, managesSerial, allowNegativeStock, p.ImageURL,
			p.EnterpriseID, p.CreatedAt, p.UpdatedAt, nil,
		)
	}
	return err
}

// ─── Tenant Operations: Presentation ─────────────────────────────────────────────────────

func (r *repository) UpsertPresentation(ctx context.Context, pr *Presentation) error {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM presentation WHERE id = ?)", pr.ID).Scan(&exists)
	if err != nil {
		return err
	}

	defaultPurchase := 0
	if pr.DefaultPurchase {
		defaultPurchase = 1
	}
	defaultSale := 0
	if pr.DefaultSale {
		defaultSale = 1
	}

	if exists {
		_, err = r.db.ExecContext(ctx,
			`UPDATE presentation SET 
			 product_id = ?, name = ?, factor = ?, barcode = ?, cost_price = ?, sale_price = ?, 
			 default_purchase = ?, default_sale = ?, updated_at = ?
			 WHERE id = ?`,
			pr.ProductID, pr.Name, pr.Factor, pr.Barcode, pr.CostPrice, pr.SalePrice,
			defaultPurchase, defaultSale, time.Now(), pr.ID,
		)
	} else {
		pr.CreatedAt = time.Now()
		pr.UpdatedAt = time.Now()
		_, err = r.db.ExecContext(ctx,
			`INSERT INTO presentation 
			 (id, product_id, name, factor, barcode, cost_price, sale_price, default_purchase, default_sale, 
			  enterprise_id, created_at, updated_at, deleted_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			pr.ID, pr.ProductID, pr.Name, pr.Factor, pr.Barcode, pr.CostPrice, pr.SalePrice,
			defaultPurchase, defaultSale, pr.EnterpriseID, pr.CreatedAt, pr.UpdatedAt, nil,
		)
	}
	return err
}