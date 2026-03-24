package cart

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateCart(ctx context.Context, cart *Cart) error {
	if cart.CartCode == "" {
		cart.CartCode = generateCartCode()
	}
	if cart.CartType == "" {
		cart.CartType = CartTypeSale
	}
	query := `
		INSERT INTO cart (cart_code, cart_type, customer_id, user_id, branch_id, enterprise_id, subtotal, discount, tax_total, total, status, notes, valid_until)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		cart.CartCode, cart.CartType, cart.CustomerID, cart.UserID, cart.BranchID, cart.EnterpriseID,
		cart.Subtotal, cart.Discount, cart.TaxTotal, cart.Total, cart.Status, cart.Notes, cart.ValidUntil).
		Scan(&cart.ID, &cart.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create cart: %w", err)
	}
	return nil
}

func (r *repository) GetCartByID(ctx context.Context, id int64) (*Cart, error) {
	cart := &Cart{}
	query := `
		SELECT id, cart_code, cart_type, customer_id, user_id, branch_id, enterprise_id, subtotal, discount, tax_total, total, status, notes, valid_until, converted_at, reference_id, reference_type, created_at, updated_at, deleted_at
		FROM cart WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cart.ID, &cart.CartCode, &cart.CartType, &cart.CustomerID, &cart.UserID, &cart.BranchID, &cart.EnterpriseID,
		&cart.Subtotal, &cart.Discount, &cart.TaxTotal, &cart.Total, &cart.Status, &cart.Notes, &cart.ValidUntil,
		&cart.ConvertedAt, &cart.ReferenceID, &cart.ReferenceType, &cart.CreatedAt, &cart.UpdatedAt, &cart.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	return cart, nil
}

func (r *repository) GetCartByCode(ctx context.Context, code string, enterpriseID int64) (*Cart, error) {
	cart := &Cart{}
	query := `
		SELECT id, cart_code, cart_type, customer_id, user_id, branch_id, enterprise_id, subtotal, discount, tax_total, total, status, notes, valid_until, converted_at, reference_id, reference_type, created_at, updated_at, deleted_at
		FROM cart WHERE cart_code = $1 AND enterprise_id = $2 AND deleted_at IS NULL`

	err := r.db.QueryRowContext(ctx, query, code, enterpriseID).Scan(
		&cart.ID, &cart.CartCode, &cart.CartType, &cart.CustomerID, &cart.UserID, &cart.BranchID, &cart.EnterpriseID,
		&cart.Subtotal, &cart.Discount, &cart.TaxTotal, &cart.Total, &cart.Status, &cart.Notes, &cart.ValidUntil,
		&cart.ConvertedAt, &cart.ReferenceID, &cart.ReferenceType, &cart.CreatedAt, &cart.UpdatedAt, &cart.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get cart by code: %w", err)
	}
	return cart, nil
}

func (r *repository) UpdateCart(ctx context.Context, cart *Cart) error {
	query := `
		UPDATE cart SET cart_type = $1, customer_id = $2, notes = $3, valid_until = $4, updated_at = NOW()
		WHERE id = $5 AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, cart.CartType, cart.CustomerID, cart.Notes, cart.ValidUntil, cart.ID)
	if err != nil {
		return fmt.Errorf("failed to update cart: %w", err)
	}
	return nil
}

func (r *repository) UpdateCartStatus(ctx context.Context, cartID int64, status string) error {
	query := `UPDATE cart SET status = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, status, cartID)
	if err != nil {
		return fmt.Errorf("failed to update cart status: %w", err)
	}
	return nil
}

func (r *repository) DeleteCart(ctx context.Context, id int64) error {
	query := `UPDATE cart SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete cart: %w", err)
	}
	return nil
}

func (r *repository) ListCarts(ctx context.Context, enterpriseID int64, filters CartFilters) ([]Cart, error) {
	query := `
		SELECT id, cart_code, cart_type, customer_id, user_id, branch_id, enterprise_id, subtotal, discount, tax_total, total, status, notes, valid_until, converted_at, reference_id, reference_type, created_at, updated_at, deleted_at
		FROM cart WHERE enterprise_id = $1 AND deleted_at IS NULL`

	args := []interface{}{enterpriseID}
	argPos := 2

	if filters.CartType != "" {
		query += fmt.Sprintf(" AND cart_type = $%d", argPos)
		args = append(args, filters.CartType)
		argPos++
	}

	if filters.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, filters.Status)
		argPos++
	}

	if filters.BranchID != nil {
		query += fmt.Sprintf(" AND branch_id = $%d", argPos)
		args = append(args, *filters.BranchID)
		argPos++
	}

	query += " ORDER BY created_at DESC"

	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, filters.Limit)
		argPos++

		offset := (filters.Page - 1) * filters.Limit
		query += fmt.Sprintf(" OFFSET $%d", argPos)
		args = append(args, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list carts: %w", err)
	}
	defer rows.Close()

	var carts []Cart
	for rows.Next() {
		var cart Cart
		if err := rows.Scan(
			&cart.ID, &cart.CartCode, &cart.CartType, &cart.CustomerID, &cart.UserID, &cart.BranchID, &cart.EnterpriseID,
			&cart.Subtotal, &cart.Discount, &cart.TaxTotal, &cart.Total, &cart.Status, &cart.Notes, &cart.ValidUntil,
			&cart.ConvertedAt, &cart.ReferenceID, &cart.ReferenceType, &cart.CreatedAt, &cart.UpdatedAt, &cart.DeletedAt,
		); err != nil {
			return nil, err
		}
		carts = append(carts, cart)
	}
	return carts, nil
}

func (r *repository) AddItem(ctx context.Context, item *CartItem) error {
	query := `
		INSERT INTO cart_item (cart_id, product_id, product_variant_id, quantity, unit_price, discount_type, discount_value, discount_amount, tax_rate, tax_amount, total, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		item.CartID, item.ProductID, item.ProductVariantID, item.Quantity, item.UnitPrice,
		item.DiscountType, item.DiscountValue, item.DiscountAmount, item.TaxRate, item.TaxAmount, item.LineTotal, item.Notes).
		Scan(&item.ID, &item.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to add item to cart: %w", err)
	}
	return nil
}

func (r *repository) UpdateItem(ctx context.Context, item *CartItem) error {
	query := `
		UPDATE cart_item SET quantity = $1, unit_price = $2, discount_type = $3, discount_value = $4, discount_amount = $5, tax_rate = $6, tax_amount = $7, total = $8, notes = $9, updated_at = NOW()
		WHERE id = $10 AND cart_id = $11`

	_, err := r.db.ExecContext(ctx, query,
		item.Quantity, item.UnitPrice, item.DiscountType, item.DiscountValue, item.DiscountAmount,
		item.TaxRate, item.TaxAmount, item.LineTotal, item.Notes, item.ID, item.CartID)
	if err != nil {
		return fmt.Errorf("failed to update cart item: %w", err)
	}
	return nil
}

func (r *repository) RemoveItem(ctx context.Context, cartID, itemID int64) error {
	query := `DELETE FROM cart_item WHERE id = $1 AND cart_id = $2`
	_, err := r.db.ExecContext(ctx, query, itemID, cartID)
	if err != nil {
		return fmt.Errorf("failed to remove cart item: %w", err)
	}
	return nil
}

func (r *repository) GetItems(ctx context.Context, cartID int64) ([]CartItem, error) {
	query := `
		SELECT id, cart_id, product_id, product_variant_id, quantity, unit_price, discount_type, discount_value, discount_amount, tax_rate, tax_amount, total, notes, created_at, updated_at
		FROM cart_item WHERE cart_id = $1
		ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}
	defer rows.Close()

	var items []CartItem
	for rows.Next() {
		var item CartItem
		if err := rows.Scan(
			&item.ID, &item.CartID, &item.ProductID, &item.ProductVariantID, &item.Quantity, &item.UnitPrice,
			&item.DiscountType, &item.DiscountValue, &item.DiscountAmount, &item.TaxRate, &item.TaxAmount,
			&item.LineTotal, &item.Notes, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *repository) UpdateCartTotals(ctx context.Context, cartID int64) error {
	query := `
		UPDATE cart SET 
			subtotal = (SELECT COALESCE(SUM(unit_price * quantity), 0) FROM cart_item WHERE cart_id = $1),
			discount = (SELECT COALESCE(SUM(discount_amount), 0) FROM cart_item WHERE cart_id = $1),
			tax_total = (SELECT COALESCE(SUM(tax_amount), 0) FROM cart_item WHERE cart_id = $1),
			total = (SELECT COALESCE(SUM(total), 0) FROM cart_item WHERE cart_id = $1),
			updated_at = NOW()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, cartID)
	if err != nil {
		return fmt.Errorf("failed to update cart totals: %w", err)
	}
	return nil
}

// Service implementation
type service struct {
	repo Repository
}

func NewService(db *sql.DB) Service {
	return &service{repo: NewRepository(db)}
}

func (s *service) CreateCart(ctx context.Context, cart *Cart) error {
	if cart.CartCode == "" {
		cart.CartCode = generateCartCode()
	}
	if cart.CartType == "" {
		cart.CartType = CartTypeSale
	}
	cart.Status = CartStatusActive
	return s.repo.CreateCart(ctx, cart)
}

func (s *service) GetCart(ctx context.Context, id int64) (*Cart, error) {
	cart, err := s.repo.GetCartByID(ctx, id)
	if err != nil {
		return nil, err
	}
	items, err := s.repo.GetItems(ctx, id)
	if err != nil {
		return nil, err
	}
	cart.Items = items
	return cart, nil
}

func (s *service) GetCartByCode(ctx context.Context, code string, enterpriseID int64) (*Cart, error) {
	cart, err := s.repo.GetCartByCode(ctx, code, enterpriseID)
	if err != nil {
		return nil, err
	}
	items, err := s.repo.GetItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}
	cart.Items = items
	return cart, nil
}

func (s *service) UpdateCart(ctx context.Context, id int64, cart *Cart) error {
	cart.ID = id
	return s.repo.UpdateCart(ctx, cart)
}

func (s *service) DeleteCart(ctx context.Context, id int64) error {
	return s.repo.DeleteCart(ctx, id)
}

func (s *service) ListCarts(ctx context.Context, enterpriseID int64, filters CartFilters) ([]Cart, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 20
	}
	return s.repo.ListCarts(ctx, enterpriseID, filters)
}

func (s *service) AddItem(ctx context.Context, cartID int64, item *CartItem) error {
	cart, err := s.repo.GetCartByID(ctx, cartID)
	if err != nil {
		return fmt.Errorf("cart not found: %w", err)
	}
	if cart.Status != CartStatusActive {
		return fmt.Errorf("cart is not active")
	}

	item.CartID = cartID
	s.calculateItemTotals(item)

	if err := s.repo.AddItem(ctx, item); err != nil {
		return fmt.Errorf("failed to add item: %w", err)
	}

	if err := s.repo.UpdateCartTotals(ctx, cartID); err != nil {
		return fmt.Errorf("failed to update cart totals: %w", err)
	}

	return nil
}

func (s *service) UpdateItem(ctx context.Context, cartID, itemID int64, item *CartItem) error {
	cart, err := s.repo.GetCartByID(ctx, cartID)
	if err != nil {
		return fmt.Errorf("cart not found: %w", err)
	}
	if cart.Status != CartStatusActive {
		return fmt.Errorf("cart is not active")
	}

	item.ID = itemID
	item.CartID = cartID
	s.calculateItemTotals(item)

	if err := s.repo.UpdateItem(ctx, item); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	if err := s.repo.UpdateCartTotals(ctx, cartID); err != nil {
		return fmt.Errorf("failed to update cart totals: %w", err)
	}

	return nil
}

func (s *service) RemoveItem(ctx context.Context, cartID, itemID int64) error {
	cart, err := s.repo.GetCartByID(ctx, cartID)
	if err != nil {
		return fmt.Errorf("cart not found: %w", err)
	}
	if cart.Status != CartStatusActive {
		return fmt.Errorf("cart is not active")
	}

	if err := s.repo.RemoveItem(ctx, cartID, itemID); err != nil {
		return fmt.Errorf("failed to remove item: %w", err)
	}

	if err := s.repo.UpdateCartTotals(ctx, cartID); err != nil {
		return fmt.Errorf("failed to update cart totals: %w", err)
	}

	return nil
}

func (s *service) GetItems(ctx context.Context, cartID int64) ([]CartItem, error) {
	return s.repo.GetItems(ctx, cartID)
}

func (s *service) ConvertToSale(ctx context.Context, cartID int64) (*Cart, error) {
	cart, err := s.repo.GetCartByID(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("cart not found: %w", err)
	}
	if cart.Status != CartStatusActive {
		return nil, fmt.Errorf("cart is not active")
	}

	items, err := s.repo.GetItems(ctx, cartID)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	now := time.Now()
	referenceID := cartID
	referenceType := "sales_order"

	// Update cart status
	if err := s.repo.UpdateCartStatus(ctx, cartID, CartStatusConverted); err != nil {
		return nil, fmt.Errorf("failed to update cart status: %w", err)
	}

	cart.Status = CartStatusConverted
	cart.ConvertedAt = &now
	cart.ReferenceID = &referenceID
	cart.ReferenceType = referenceType
	cart.Items = items

	return cart, nil
}

func (s *service) ConvertToQuotation(ctx context.Context, cartID int64, validDays int) (*Cart, error) {
	cart, err := s.repo.GetCartByID(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("cart not found: %w", err)
	}
	if cart.Status != CartStatusActive {
		return nil, fmt.Errorf("cart is not active")
	}

	// Update cart type to quotation
	cart.CartType = CartTypeQuotation
	if validDays > 0 {
		validUntil := time.Now().AddDate(0, 0, validDays)
		cart.ValidUntil = &validUntil
	}

	if err := s.repo.UpdateCart(ctx, cart); err != nil {
		return nil, fmt.Errorf("failed to convert to quotation: %w", err)
	}

	return s.GetCart(ctx, cartID)
}

func (s *service) SetCustomer(ctx context.Context, cartID int64, customerID *int64) error {
	cart, err := s.repo.GetCartByID(ctx, cartID)
	if err != nil {
		return err
	}
	if cart.Status != CartStatusActive {
		return fmt.Errorf("cart is not active")
	}

	cart.CustomerID = customerID
	return s.repo.UpdateCart(ctx, cart)
}

func (s *service) ApplyDiscount(ctx context.Context, cartID int64, discountType string, discountValue float64) error {
	cart, err := s.repo.GetCartByID(ctx, cartID)
	if err != nil {
		return err
	}
	if cart.Status != CartStatusActive {
		return fmt.Errorf("cart is not active")
	}

	items, err := s.repo.GetItems(ctx, cartID)
	if err != nil {
		return err
	}

	for i := range items {
		items[i].DiscountType = discountType
		items[i].DiscountValue = discountValue
		s.calculateItemTotals(&items[i])
		if err := s.repo.UpdateItem(ctx, &items[i]); err != nil {
			return err
		}
	}

	return s.repo.UpdateCartTotals(ctx, cartID)
}

func (s *service) ApplyItemDiscount(ctx context.Context, cartID, itemID int64, discountType string, discountValue float64) error {
	cart, err := s.repo.GetCartByID(ctx, cartID)
	if err != nil {
		return err
	}
	if cart.Status != CartStatusActive {
		return fmt.Errorf("cart is not active")
	}

	items, err := s.repo.GetItems(ctx, cartID)
	if err != nil {
		return err
	}

	for i := range items {
		if items[i].ID == itemID {
			items[i].DiscountType = discountType
			items[i].DiscountValue = discountValue
			s.calculateItemTotals(&items[i])
			if err := s.repo.UpdateItem(ctx, &items[i]); err != nil {
				return err
			}
			break
		}
	}

	return s.repo.UpdateCartTotals(ctx, cartID)
}

func (s *service) calculateItemTotals(item *CartItem) {
	subtotal := item.UnitPrice * float64(item.Quantity)
	
	// Calculate discount
	var discountAmount float64
	if item.DiscountType == DiscountTypePercentage {
		discountAmount = subtotal * (item.DiscountValue / 100)
	} else if item.DiscountType == DiscountTypeFixed {
		discountAmount = item.DiscountValue * float64(item.Quantity)
	}
	
	taxableAmount := subtotal - discountAmount
	taxAmount := taxableAmount * (item.TaxRate / 100)
	total := taxableAmount + taxAmount

	item.DiscountAmount = math.Round(discountAmount*100) / 100
	item.TaxAmount = math.Round(taxAmount*100) / 100
	item.LineTotal = math.Round(total*100) / 100
}

func generateCartCode() string {
	return fmt.Sprintf("CART-%s", uuid.New().String()[:8])
}
