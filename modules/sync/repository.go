package sync

import (
	"context"
	"time"

	"github.com/cloud-tech-develop/aura-back/internal/db"
	"github.com/cloud-tech-develop/aura-back/modules/products"
	"github.com/cloud-tech-develop/aura-back/modules/third-parties"
	"github.com/cloud-tech-develop/aura-back/modules/sales"
)

type repository struct {
	db db.Querier
}

func NewRepository(db db.Querier) *repository {
	return &repository{db: db}
}

// GetUpdates since last login/sync
func (r *repository) GetProductUpdates(ctx context.Context, lastSync time.Time, enterpriseID int64) ([]products.Product, error) {
	var list []products.Product
	query := `SELECT id, sku, name, description, category_id, brand_id, cost_price, sale_price, tax_rate, min_stock, current_stock, image_url, status, enterprise_id, global_id, sync_status, last_synced_at, created_at, updated_at, deleted_at 
	          FROM product WHERE enterprise_id = $1 AND (updated_at > $2 OR created_at > $2)`
	
	rows, err := r.db.QueryContext(ctx, query, enterpriseID, lastSync)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p products.Product
		err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.CategoryID, &p.BrandID, &p.CostPrice, &p.SalePrice, &p.TaxRate, &p.MinStock, &p.CurrentStock, &p.ImageURL, &p.Status, &p.EnterpriseID, &p.GlobalID, &p.SyncStatus, &p.LastSyncedAt, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

// Similar methods for Category, Brand, ThirdParty, Sales, Invoices...
// To keep it short I will implement a generic fetcher or just the most important ones for now.

func (r *repository) UpsertProduct(ctx context.Context, p products.Product) error {
	query := `INSERT INTO product (sku, name, description, category_id, brand_id, cost_price, sale_price, tax_rate, min_stock, current_stock, image_url, status, enterprise_id, global_id, sync_status, last_synced_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, 'SYNCED', NOW())
	          ON CONFLICT (enterprise_id, global_id) DO UPDATE SET
	          sku = EXCLUDED.sku, name = EXCLUDED.name, description = EXCLUDED.description, category_id = EXCLUDED.category_id, brand_id = EXCLUDED.brand_id, 
	          cost_price = EXCLUDED.cost_price, sale_price = EXCLUDED.sale_price, tax_rate = EXCLUDED.tax_rate, min_stock = EXCLUDED.min_stock, 
	          current_stock = EXCLUDED.current_stock, image_url = EXCLUDED.image_url, status = EXCLUDED.status, last_synced_at = NOW(), sync_status = 'SYNCED'`
	
	_, err := r.db.ExecContext(ctx, query, p.SKU, p.Name, p.Description, p.CategoryID, p.BrandID, p.CostPrice, p.SalePrice, p.TaxRate, p.MinStock, p.CurrentStock, p.ImageURL, p.Status, p.EnterpriseID, p.GlobalID)
	return err
}

func (r *repository) UpsertThirdParty(ctx context.Context, tp thirdparties.ThirdParty) error {
	query := `INSERT INTO third_parties (user_id, first_name, last_name, document_number, document_type, personal_email, commercial_name, address, phone, additional_email, tax_responsibility, is_client, is_provider, is_employee, municipality_id, municipality, global_id, sync_status, last_synced_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, 'SYNCED', NOW())
	          ON CONFLICT (global_id) DO UPDATE SET
	          first_name = EXCLUDED.first_name, last_name = EXCLUDED.last_name, document_number = EXCLUDED.document_number, personal_email = EXCLUDED.personal_email, 
	          address = EXCLUDED.address, phone = EXCLUDED.phone, sync_status = 'SYNCED', last_synced_at = NOW()`
	_, err := r.db.ExecContext(ctx, query, tp.UserID, tp.FirstName, tp.LastName, tp.DocumentNumber, tp.DocumentType, tp.PersonalEmail, tp.CommercialName, tp.Address, tp.Phone, tp.AdditionalEmail, tp.TaxResponsibility, tp.IsClient, tp.IsProvider, tp.IsEmployee, tp.MunicipalityID, tp.Municipality, tp.GlobalID)
	return err
}

func (r *repository) GetThirdPartyUpdates(ctx context.Context, lastSync time.Time, enterpriseID int64) ([]thirdparties.ThirdParty, error) {
	var list []thirdparties.ThirdParty
	query := `SELECT id, user_id, first_name, last_name, document_number, document_type, personal_email, commercial_name, address, phone, additional_email, tax_responsibility, is_client, is_provider, is_employee, municipality_id, municipality, global_id, created_at, deleted_at 
	          FROM third_parties WHERE (created_at > $1 OR deleted_at > $1)`
	
	rows, err := r.db.QueryContext(ctx, query, lastSync)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tp thirdparties.ThirdParty
		err := rows.Scan(&tp.ID, &tp.UserID, &tp.FirstName, &tp.LastName, &tp.DocumentNumber, &tp.DocumentType, &tp.PersonalEmail, &tp.CommercialName, &tp.Address, &tp.Phone, &tp.AdditionalEmail, &tp.TaxResponsibility, &tp.IsClient, &tp.IsProvider, &tp.IsEmployee, &tp.MunicipalityID, &tp.Municipality, &tp.GlobalID, &tp.CreatedAt, &tp.DeletedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, tp)
	}
	return list, nil
}

func (r *repository) GetPendingSales(ctx context.Context, enterpriseID int64) ([]sales.SalesOrder, error) {
	var list []sales.SalesOrder
	query := `SELECT id, order_number, customer_id, user_id, branch_id, enterprise_id, subtotal, discount, tax_total, total, status, notes, global_id, created_at, updated_at
	          FROM sales_order WHERE enterprise_id = $1 AND sync_status = 'PENDING'`
	
	rows, err := r.db.QueryContext(ctx, query, enterpriseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s sales.SalesOrder
		err := rows.Scan(&s.ID, &s.OrderNumber, &s.CustomerID, &s.UserID, &s.BranchID, &s.EnterpriseID, &s.Subtotal, &s.Discount, &s.TaxTotal, &s.Total, &s.Status, &s.Notes, &s.GlobalID, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}
