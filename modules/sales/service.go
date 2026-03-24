package sales

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateOrder(ctx context.Context, order *SalesOrder) error {
	if order.OrderNumber == "" {
		order.OrderNumber = generateOrderNumber()
	}
	query := `
		INSERT INTO sales_order (order_number, customer_id, user_id, branch_id, enterprise_id, subtotal, discount, tax_total, total, status, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		order.OrderNumber, order.CustomerID, order.UserID, order.BranchID, order.EnterpriseID,
		order.Subtotal, order.Discount, order.TaxTotal, order.Total, order.Status, order.Notes).
		Scan(&order.ID, &order.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

func (r *repository) GetOrderByID(ctx context.Context, id int64) (*SalesOrder, error) {
	order := &SalesOrder{}
	query := `
		SELECT id, order_number, customer_id, user_id, branch_id, enterprise_id, subtotal, discount, tax_total, total, status, notes, created_at, updated_at
		FROM sales_order WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID, &order.OrderNumber, &order.CustomerID, &order.UserID, &order.BranchID, &order.EnterpriseID,
		&order.Subtotal, &order.Discount, &order.TaxTotal, &order.Total, &order.Status, &order.Notes,
		&order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}

func (r *repository) GetOrderByNumber(ctx context.Context, orderNumber string, enterpriseID int64) (*SalesOrder, error) {
	order := &SalesOrder{}
	query := `
		SELECT id, order_number, customer_id, user_id, branch_id, enterprise_id, subtotal, discount, tax_total, total, status, notes, created_at, updated_at
		FROM sales_order WHERE order_number = $1 AND enterprise_id = $2`

	err := r.db.QueryRowContext(ctx, query, orderNumber, enterpriseID).Scan(
		&order.ID, &order.OrderNumber, &order.CustomerID, &order.UserID, &order.BranchID, &order.EnterpriseID,
		&order.Subtotal, &order.Discount, &order.TaxTotal, &order.Total, &order.Status, &order.Notes,
		&order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get order by number: %w", err)
	}
	return order, nil
}

func (r *repository) UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
	query := `UPDATE sales_order SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}

func (r *repository) GetOrders(ctx context.Context, enterpriseID int64, filters OrderFilters) ([]SalesOrder, error) {
	query := `
		SELECT id, order_number, customer_id, user_id, branch_id, enterprise_id, subtotal, discount, tax_total, total, status, notes, created_at, updated_at
		FROM sales_order WHERE enterprise_id = $1`

	args := []interface{}{enterpriseID}
	argPos := 2

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

	if filters.CustomerID != nil {
		query += fmt.Sprintf(" AND customer_id = $%d", argPos)
		args = append(args, *filters.CustomerID)
		argPos++
	}

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argPos)
		args = append(args, *filters.StartDate)
		argPos++
	}

	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argPos)
		args = append(args, *filters.EndDate)
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
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	var orders []SalesOrder
	for rows.Next() {
		var order SalesOrder
		if err := rows.Scan(
			&order.ID, &order.OrderNumber, &order.CustomerID, &order.UserID, &order.BranchID, &order.EnterpriseID,
			&order.Subtotal, &order.Discount, &order.TaxTotal, &order.Total, &order.Status, &order.Notes,
			&order.CreatedAt, &order.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *repository) CreateOrderItem(ctx context.Context, item *SalesOrderItem) error {
	query := `
		INSERT INTO sales_order_item (sales_order_id, product_id, quantity, unit_price, discount_percent, discount_amount, tax_rate, tax_amount, total)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		item.SalesOrderID, item.ProductID, item.Quantity, item.UnitPrice,
		item.DiscountPercent, item.DiscountAmount, item.TaxRate, item.TaxAmount, item.Total).
		Scan(&item.ID, &item.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create order item: %w", err)
	}
	return nil
}

func (r *repository) GetOrderItems(ctx context.Context, orderID int64) ([]SalesOrderItem, error) {
	query := `
		SELECT id, sales_order_id, product_id, quantity, unit_price, discount_percent, discount_amount, tax_rate, tax_amount, total, created_at
		FROM sales_order_item WHERE sales_order_id = $1
		ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	var items []SalesOrderItem
	for rows.Next() {
		var item SalesOrderItem
		if err := rows.Scan(
			&item.ID, &item.SalesOrderID, &item.ProductID, &item.Quantity, &item.UnitPrice,
			&item.DiscountPercent, &item.DiscountAmount, &item.TaxRate, &item.TaxAmount,
			&item.Total, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// Service implementation
type service struct {
	repo Repository
}

func NewService(db *sql.DB) Service {
	return &service{repo: NewRepository(db)}
}

func (s *service) CreateOrderFromCart(ctx context.Context, cartID int64) (*SalesOrder, error) {
	// Get cart from cart module - we'll need to inject it or use database directly
	// For now, we'll create a simple order from provided data
	// This should be called with cart data already loaded
	
	order := &SalesOrder{
		Status: StatusPendingPayment,
	}
	return order, nil
}

func (s *service) GetOrder(ctx context.Context, id int64) (*SalesOrder, error) {
	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	items, err := s.repo.GetOrderItems(ctx, id)
	if err != nil {
		return nil, err
	}
	
	order.Items = items
	return order, nil
}

func (s *service) GetOrders(ctx context.Context, enterpriseID int64, filters OrderFilters) ([]SalesOrder, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 20
	}
	return s.repo.GetOrders(ctx, enterpriseID, filters)
}

func (s *service) UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}
	
	// Validate status transitions
	switch order.Status {
	case StatusPendingPayment:
		if status != StatusPaid && status != StatusCancelled {
			return fmt.Errorf("invalid status transition from %s to %s", order.Status, status)
		}
	case StatusPaid:
		if status != StatusCompleted && status != StatusCancelled {
			return fmt.Errorf("invalid status transition from %s to %s", order.Status, status)
		}
	case StatusCompleted:
		return fmt.Errorf("cannot change status of completed order")
	case StatusCancelled:
		return fmt.Errorf("cannot change status of cancelled order")
	}
	
	return s.repo.UpdateOrderStatus(ctx, orderID, status)
}

func (s *service) CancelOrder(ctx context.Context, orderID int64) error {
	return s.UpdateOrderStatus(ctx, orderID, StatusCancelled)
}

func (s *service) CompleteOrder(ctx context.Context, orderID int64) error {
	return s.UpdateOrderStatus(ctx, orderID, StatusCompleted)
}

func generateOrderNumber() string {
	timestamp := time.Now().Format("200601021504")
	return fmt.Sprintf("ORD-%s-%s", timestamp, uuid.New().String()[:4])
}

// UpdateOrder updates an existing order
func (s *service) UpdateOrder(ctx context.Context, order *SalesOrder) error {
	existing, err := s.repo.GetOrderByID(ctx, order.ID)
	if err != nil {
		return err
	}
	
	if existing.Status == StatusCompleted || existing.Status == StatusCancelled {
		return fmt.Errorf("cannot update a %s order", existing.Status)
	}
	
	query := `
		UPDATE sales_order SET 
			customer_id = $1, notes = $2, subtotal = $3, discount = $4, tax_total = $5, total = $6, updated_at = NOW()
		WHERE id = $7`
	
	_, err = s.repo.(*repository).db.ExecContext(ctx, query,
		order.CustomerID, order.Notes, order.Subtotal, order.Discount, order.TaxTotal, order.Total, order.ID)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}
	
	return nil
}
